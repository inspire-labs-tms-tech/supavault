package com.inspiretmstech.supavault.commands;

import com.google.gson.JsonArray;
import com.inspiretmstech.supavault.ExceptionHandler;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.db.Database;
import com.inspiretmstech.supavault.db.gen.Tables;
import com.inspiretmstech.supavault.db.gen.tables.records.VariablesRecord;
import com.inspiretmstech.supavault.utils.gson.GSON;
import org.jooq.Field;
import org.jooq.Result;
import picocli.CommandLine;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

@CommandLine.Command(
        name = "variables",
        description = "manage variables in a project"
)
public class Variables extends Loggable {

    public Variables() {
        super(Variables.class);
    }

    @CommandLine.Command(
            name = "create",
            description = "create a new variable"
    )
    public int create(
            @CommandLine.Parameters(paramLabel = "name", description = "the name for the variable (e.g., ENV_VAR)") String var,
            @CommandLine.Option(names = {"--description"}, description = "a description of the variable", defaultValue = "") String description,
            @CommandLine.Option(names = {"--default"}, description = "a default value for the variable at the project-level (can be over-ridden at the client-level)", defaultValue = "") String defaultValue
    ) {

        logger.debug("creating vairble {} in project {}...", var, Project.projectID);

        // validate variable name
        if (!var.matches("^[A-Z_]{1,}[A-Z0-9_]*$"))
            throw new RuntimeException("display must start with a capital letter or an underscore and contain only capital letters, numbers, and underscores");

        try {
            Database.with().unsafely(db -> {
                VariablesRecord record = new VariablesRecord();
                record.setId(var);
                record.setDescription(description);
                record.setDefault(defaultValue);
                record.setProjectId(Project.projectID);
                record = db.insertInto(Tables.VARIABLES).set(record).returning().fetchOne();
                if (Objects.isNull(record)) throw new RuntimeException("unable to create variable");
                logger.info("created!");
            });
        } catch (Exception e) {
            if (e.getMessage().contains("duplicate key value"))
                logger.error("variable \"{}\" already exists", var);
            else ExceptionHandler.handle(e);
            return 1;
        }
        return 0;
    }

    @CommandLine.Command(
            name = "list",
            description = "list all existing variables"
    )
    public int list(
            @CommandLine.Option(names = {"--json"}, description = "output as a JSON object") boolean json
    ) {
        Database.with().execute(db -> {
            Result<VariablesRecord> vars = db
                    .selectFrom(Tables.VARIABLES)
                    .where(Tables.VARIABLES.PROJECT_ID.eq(Project.projectID))
                    .fetch();
            JsonArray output = new JsonArray();
            if (json) {
                for (VariablesRecord var : vars) {
                    Map<String, Object> rowMap = new HashMap<>();
                    for (Field<?> field : var.fields()) rowMap.put(field.getName(), var.getValue(field));
                    output.add(GSON.GLOBAL.toJsonTree(rowMap).getAsJsonObject());
                }
                logger.info(output.toString());
            } else {
                if (vars.isEmpty())
                    logger.warn("no variables created for project \"{}\"", Project.projectID);
                else {
                    logger.info("Variables ({}): ", vars.size());
                    for (VariablesRecord var : vars)
                        logger.info("  - [{}]: {}(default: \"{}\")", var.getId(), var.getDescription().isBlank() ? "" : var.getDescription().trim() + " ", var.getDefault());
                }
            }
        });
        return 0;
    }

    @CommandLine.Command(
            name = "delete",
            description = "delete an existing variable"
    )
    public int delete(
            @CommandLine.Parameters(paramLabel = "id", description = "the name of the variable to delete")
            String var
    ) {
        logger.debug("deleting variable {} for project {}...", var, Project.projectID);

        // delete the project
        Database.with().execute(db -> {
            VariablesRecord r = db.selectFrom(Tables.VARIABLES)
                    .where(Tables.VARIABLES.ID.eq(var))
                    .and(Tables.VARIABLES.PROJECT_ID.eq(Project.projectID))
                    .fetchOne();
            if (Objects.isNull(r))
                throw new RuntimeException("variable \"" + var + "\" in project \"" + Project.projectID + "\" does not exist");
            r.delete();
        });

        return 0;
    }
}
