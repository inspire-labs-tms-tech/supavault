package com.inspiretmstech.supavault.commands;

import com.inspiretmstech.supavault.ExceptionHandler;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.db.Database;
import com.inspiretmstech.supavault.db.gen.Tables;
import com.inspiretmstech.supavault.db.gen.tables.records.ProjectsRecord;
import org.jooq.Result;
import org.jooq.exception.IntegrityConstraintViolationException;
import picocli.CommandLine;

import java.sql.SQLException;
import java.util.Objects;
import java.util.UUID;

@CommandLine.Command(
        name = "projects",
        description = "manage projects in a supavault instance"
)
public class Projects extends Loggable {

    public Projects() {
        super(Projects.class);
    }

    @CommandLine.Command(name = "delete")
    public int delete(
            @CommandLine.Parameters(paramLabel = "id", description = "the (UUID) id of the project to delete")
            String id
    ) {
        logger.debug("deleting project {}...", id);

        UUID uuid;
        try {
            uuid = UUID.fromString(id);
        } catch (IllegalArgumentException e) {
            logger.error("\"{}\" is not a valid id", id);
            return 1;
        }

        Database.with().execute(db -> {

            ProjectsRecord r = db.selectFrom(Tables.PROJECTS)
                            .where(Tables.PROJECTS.ID.eq(uuid))
                                    .fetchOne();
            if(Objects.isNull(r)) throw new RuntimeException("project with id \"" + id + "\" does not exist");

            r.delete();
        });

        return 0;


    }

    @CommandLine.Command(name = "create")
    public int create(@CommandLine.Parameters(paramLabel = "display", description = "the name for the project") String display) {
        logger.debug("creating project \"{}\"...", display);
        try {
            Database.with().unsafely(db -> {
                ProjectsRecord record = new ProjectsRecord();
                record.setDisplay(display);
                db.insertInto(Tables.PROJECTS).set(record).execute();
            });
        } catch (Exception e) {
            if (e.getMessage().contains("duplicate key value"))
                logger.error("project with display \"{}\" already exists", display);
            else ExceptionHandler.handle(e);
            return 1;
        }
        return 0;
    }

    @CommandLine.Command(name = "list")
    public int list() {
        Database.with().execute(db -> {
            Result<ProjectsRecord> projects = db.selectFrom(Tables.PROJECTS).fetch();
            if (projects.isEmpty()) logger.warn("no projects created");
            else {
                logger.info("Projects ({}): ", projects.size());
                for (ProjectsRecord project : projects) {
                    logger.info("  - [{}]: {} ", project.getDisplay(), project.getId());
                }
            }
        });
        return 0;
    }

}
