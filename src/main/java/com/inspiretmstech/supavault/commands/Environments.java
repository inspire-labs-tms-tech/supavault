package com.inspiretmstech.supavault.commands;

import com.google.gson.JsonArray;
import com.inspiretmstech.supavault.ExceptionHandler;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.db.Database;
import com.inspiretmstech.supavault.db.gen.Tables;
import com.inspiretmstech.supavault.db.gen.tables.records.EnvironmentsRecord;
import com.inspiretmstech.supavault.utils.gson.GSON;
import org.jooq.Field;
import org.jooq.Result;
import picocli.CommandLine;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.UUID;

@CommandLine.Command(
        name = "environments",
        description = "manage environments in a project"
)
public class Environments extends Loggable {

    public Environments() {
        super(Environments.class);
    }

    @CommandLine.Command(
            name = "create",
            description = "create a new environment"
    )
    public int create(
            @CommandLine.Parameters(paramLabel = "display", description = "the name for the environment") String display
    ) {
        logger.debug("creating environment \"{}\" in project \"{}\"...", display, Project.projectID);

        // validate display
        display = display.trim();
        if (display.isEmpty()) {
            logger.error("display is empty");
            return 1;
        }

        try {
            String finalDisplay = display;
            Database.with().unsafely(db -> {
                EnvironmentsRecord record = new EnvironmentsRecord();
                record.setDisplay(finalDisplay);
                record.setProjectId(Project.projectID);
                record = db.insertInto(Tables.ENVIRONMENTS).set(record).returning().fetchOne();
                if (Objects.isNull(record)) throw new RuntimeException("unable to create environment");
                logger.info(record.getId().toString());
            });
        } catch (Exception e) {
            if (e.getMessage().contains("duplicate key value"))
                logger.error("environment with display \"{}\" already exists", display);
            else ExceptionHandler.handle(e);
            return 1;
        }
        return 0;
    }

    @CommandLine.Command(
            name = "list",
            description = "list all existing environments"
    )
    public int list(
            @CommandLine.Option(names = {"--json"}, description = "output as a JSON object") boolean json
    ) {
        Database.with().execute(db -> {
            Result<EnvironmentsRecord> envs = db
                    .selectFrom(Tables.ENVIRONMENTS)
                    .where(Tables.ENVIRONMENTS.PROJECT_ID.eq(Project.projectID))
                    .fetch();
            JsonArray output = new JsonArray();
            if (json) {
                for (EnvironmentsRecord env : envs) {
                    Map<String, Object> rowMap = new HashMap<>();
                    for (Field<?> field : env.fields()) rowMap.put(field.getName(), env.getValue(field));
                    output.add(GSON.GLOBAL.toJsonTree(rowMap).getAsJsonObject());
                }
                logger.info(output.toString());
            } else {
                if (envs.isEmpty())
                    logger.warn("no environments created for project \"{}\"", Project.projectID);
                else {
                    logger.info("Environments ({}): ", envs.size());
                    for (EnvironmentsRecord env : envs)
                        logger.info("  - [{}]: {} ", env.getDisplay(), env.getId());
                }
            }
        });
        return 0;
    }

    @CommandLine.Command(
            name = "delete",
            description = "delete an existing environment"
    )
    public int delete(
            @CommandLine.Parameters(paramLabel = "id", description = "the (UUID) id of the environment to delete")
            UUID id
    ) {
        logger.debug("deleting environment {} for project {}...", id, Project.projectID);

        // delete the project
        Database.with().execute(db -> {
            EnvironmentsRecord r = db.selectFrom(Tables.ENVIRONMENTS)
                    .where(Tables.ENVIRONMENTS.ID.eq(id))
                    .and(Tables.ENVIRONMENTS.PROJECT_ID.eq(Project.projectID))
                    .fetchOne();
            if (Objects.isNull(r))
                throw new RuntimeException("environment with id \"" + id + "\" in project \"" + Project.projectID + "\" does not exist");
            r.delete();
        });

        return 0;
    }

}
