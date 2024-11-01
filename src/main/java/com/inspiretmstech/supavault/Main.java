package com.inspiretmstech.supavault;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import com.inspiretmstech.supavault.commands.Environments;
import com.inspiretmstech.supavault.commands.Projects;
import com.inspiretmstech.supavault.constants.LogLevel;
import com.inspiretmstech.supavault.constants.Version;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

@CommandLine.Command(
        name = "supavault",
        subcommands = {
                Projects.class,
                Environments.class
        },
        versionProvider = Version.class,
        mixinStandardHelpOptions = true
)
class Main {

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

    public static void main(String[] args) {
        int exitCode;
        try {
            exitCode = new CommandLine(new Main())
                    .setExecutionExceptionHandler(new ExceptionHandler())
                    .execute(args);
        } catch (Exception e) {
            ExceptionHandler.handle(e);
            exitCode = 1;
        }
        System.exit(exitCode);
    }

}
