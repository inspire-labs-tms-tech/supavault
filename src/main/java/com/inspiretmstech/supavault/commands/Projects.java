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

@CommandLine.Command(
        name = "projects",
        description = "manage projects in a supavault instance"
)
public class Projects extends Loggable {

    public Projects() {
        super(Projects.class);
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
