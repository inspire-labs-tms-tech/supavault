package com.inspiretmstech.supavault.commands.projects;

import com.google.gson.JsonArray;
import com.inspiretmstech.supavault.utils.ExceptionHandler;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.commands.projects.project.Project;
import com.inspiretmstech.supavault.db.Database;
import com.inspiretmstech.supavault.db.gen.Tables;
import com.inspiretmstech.supavault.db.gen.tables.records.ProjectsRecord;
import com.inspiretmstech.supavault.utils.gson.GSON;
import org.jooq.Field;
import org.jooq.Result;
import picocli.CommandLine;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.UUID;

@CommandLine.Command(
        name = "projects",
        description = "manage projects in a supavault instance",
        subcommands = {
                Project.class
        }
)
public class Projects extends Loggable {

    public Projects() {
        super(Projects.class);
    }

    @CommandLine.Command(
            name = "create",
            description = "create a new project"
    )
    public int create(@CommandLine.Parameters(paramLabel = "display", description = "the name for the project") String display) {
        logger.debug("creating project \"{}\"...", display);
        try {
            Database.with().unsafely(db -> {
                ProjectsRecord record = new ProjectsRecord();
                record.setDisplay(display);
                record = db.insertInto(Tables.PROJECTS).set(record).returning().fetchOne();
                if (Objects.isNull(record)) throw new RuntimeException("unable to create project");
                logger.info(record.getId().toString());
            });
        } catch (Exception e) {
            if (e.getMessage().contains("duplicate key value"))
                logger.error("project with display \"{}\" already exists", display);
            else ExceptionHandler.handle(e);
            return 1;
        }
        return 0;
    }

    @CommandLine.Command(
            name = "list",
            description = "list all existing projects"
    )
    public int list(
            @CommandLine.Option(names = {"--json"}, description = "output as a JSON object") boolean json
    ) {
        Database.with().execute(db -> {
            Result<ProjectsRecord> projects = db.selectFrom(Tables.PROJECTS).fetch();
            JsonArray output = new JsonArray();
            if (json) {
                for (ProjectsRecord project : projects) {
                    Map<String, Object> rowMap = new HashMap<>();
                    for (Field<?> field : project.fields()) rowMap.put(field.getName(), project.getValue(field));
                    output.add(GSON.GLOBAL.toJsonTree(rowMap).getAsJsonObject());
                }
                logger.info(output.toString());
            } else {
                if (projects.isEmpty()) logger.warn("no projects created");
                else {
                    logger.info("Projects ({}): ", projects.size());
                    for (ProjectsRecord project : projects)
                        logger.info("  - [{}]: {} ", project.getDisplay(), project.getId());
                }
            }
        });
        return 0;
    }

    @CommandLine.Command(
            name = "delete",
            description = "delete an existing project"
    )
    public int delete(
            @CommandLine.Parameters(paramLabel = "id", description = "the (UUID) id of the project to delete")
            UUID id
    ) {
        logger.debug("deleting project {}...", id);

        // delete the project
        Database.with().execute(db -> {
            ProjectsRecord r = db.selectFrom(Tables.PROJECTS)
                    .where(Tables.PROJECTS.ID.eq(id))
                    .fetchOne();
            if (Objects.isNull(r)) throw new RuntimeException("project with id \"" + id + "\" does not exist");
            r.delete();
        });
        logger.info("deleted!");
        return 0;
    }

}
