package com.inspiretmstech.supavault.commands;

import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.db.Database;
import com.inspiretmstech.supavault.db.gen.Tables;
import com.inspiretmstech.supavault.db.gen.tables.records.ProjectsRecord;
import org.jooq.Result;
import org.jooq.SelectWhereStep;
import picocli.CommandLine;

@CommandLine.Command(
        name = "projects",
        description = "manage projects in a supavault instance"
)
public class Projects extends Loggable {

    public Projects() {
        super(Projects.class);
    }

//    @CommandLine.Command(name = "create")
//    public int create(@CommandLine.ArgGroup(names = {"--display"}, required = true, paramLabel = "\"New Project Name\"") String display) {
//        logger.info("{}", display);
////        Database.with().execute(db -> {
////
////            ProjectsRecord record = new ProjectsRecord();
////            record.setDisplay();
////        });
//        return 0;
//    }

    @CommandLine.Command(name = "list")
    public int list() {
        Database.with().execute(db -> {
            Result<ProjectsRecord> projects = db.selectFrom(Tables.PROJECTS).fetch();
            if(projects.isEmpty()) logger.warn("no projects created");
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
