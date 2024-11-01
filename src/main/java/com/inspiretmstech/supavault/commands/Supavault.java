package com.inspiretmstech.supavault.commands;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import com.inspiretmstech.supavault.commands.auth.Auth;
import com.inspiretmstech.supavault.utils.ConnectionSettings;
import com.inspiretmstech.supavault.utils.GlobalOptions;
import com.inspiretmstech.supavault.commands.projects.Projects;
import com.inspiretmstech.supavault.constants.LogLevel;
import com.inspiretmstech.supavault.constants.Version;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;


@CommandLine.Command(
        name = "supavault",
        subcommands = {
                Projects.class,
                Auth.class
        },
        versionProvider = Version.class,
        mixinStandardHelpOptions = true
)
public class Supavault {

    // setup logging
    static {
        System.setProperty("org.jooq.no-logo", "true");
        System.setProperty("org.jooq.no-tips", "true");
        System.setProperty("slf4j.internal.verbosity", "OFF");
        ((Logger) LoggerFactory.getLogger("org.jooq")).setLevel(Level.OFF);
        ((Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME)).setLevel(LogLevel.DEFAULT);
    }

    // inject the DB connection settings
    @CommandLine.ArgGroup(exclusive = false)
    ConnectionSettings settings;

    // inject global options
    @CommandLine.ArgGroup(exclusive = false)
    GlobalOptions globalOptions;

    @CommandLine.Option(names = {"-v", "--version"}, versionHelp = true,
            description = "print version information and exit")
    boolean versionRequested;

}
