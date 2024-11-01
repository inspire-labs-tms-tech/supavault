package com.inspiretmstech.supavault.commands;

import com.inspiretmstech.supavault.ExceptionHandler;
import com.inspiretmstech.supavault.bases.CRUDCommand;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.db.Database;
import com.inspiretmstech.supavault.db.gen.Tables;
import com.inspiretmstech.supavault.db.gen.tables.records.VariablesRecord;
import picocli.CommandLine;

import java.util.Objects;
import java.util.UUID;

@CommandLine.Command(
        name = "variables",
        description = "manage variables in a project"
)
public class Variables extends Loggable implements CRUDCommand {

    public Variables() {
        super(Variables.class);
    }

    @CommandLine.Command(
            name = "create",
            description = "create a new variable"
    )
    public int create(
            @CommandLine.Parameters(paramLabel = "name", description = "the name for the variable (e.g., ENV_VAR)") String var
    ) {

        logger.debug("creating vairble {} in project {}...", var, Project.projectID);

        // validate variable name
        if (!var.matches("^[A-Z_]{1,}[A-Z0-9_]*$"))
            throw new RuntimeException("display must start with a capital letter or an underscore and contain only capital letters, numbers, and underscores");

        try {
            Database.with().unsafely(db -> {
                VariablesRecord record = new VariablesRecord();
                record.setId(var);
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

    @Override
    public int list(boolean json) {
        return 0;
    }

    @Override
    public int delete(UUID id) {
        return 0;
    }
}
