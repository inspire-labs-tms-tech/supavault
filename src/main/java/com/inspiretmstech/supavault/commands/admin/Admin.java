package com.inspiretmstech.supavault.commands.admin;

import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.db.Migrator;
import picocli.CommandLine;


@CommandLine.Command(
        name = "admin",
        description = "manage administrative supavault functions"
)
public class Admin extends Loggable {

    public Admin() {
        super(Admin.class);
    }

    @CommandLine.Command(
            name = "update",
            description = "update supavault migrations to the latest version"
    )
    public int update() {
        try {
            Migrator.migrate();
            logger.info("migrations applied!");
            return 0;
        } catch (Exception e) {
            logger.error("migration failed: {}", e.getMessage());
            return 1;
        }
    }

}
