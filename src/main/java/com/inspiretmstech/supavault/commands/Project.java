package com.inspiretmstech.supavault.commands;

import picocli.CommandLine;

import java.util.UUID;

@CommandLine.Command(
        name = "project",
        description = "manage an individual project",
        subcommands = {
                Environments.class,
                Variables.class
        }
)
public class Project {

    public static UUID projectID;

    @CommandLine.Parameters(paramLabel = "project", description = "the id of the project to create this environment in", scope = CommandLine.ScopeType.INHERIT)
    private void setProjectID(UUID projectID) {
        Project.projectID = projectID;
    }

}
