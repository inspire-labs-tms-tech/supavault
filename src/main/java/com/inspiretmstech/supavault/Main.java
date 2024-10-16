package com.inspiretmstech.supavault;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import com.inspiretmstech.supavault.constants.Version;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

@CommandLine.Command(
        name = "supavault",
        subcommands = {},
        versionProvider = Version.class
)
class Main implements Runnable {

    private static Logger logger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);

    static {
        System.setProperty("slf4j.internal.verbosity", "OFF");
        ((Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME)).setLevel(Level.ERROR);
    }

    @CommandLine.Option(names = {"-v", "--version"}, versionHelp = true,
            description = "print version information and exit")
    boolean versionRequested;
    @CommandLine.Option(names = {"-h", "--help"}, usageHelp = true, description = "display this usage menu")
    private boolean help;

    public static void main(String[] args) {
        int exitCode = new CommandLine(new Main()).execute(args);
        System.exit(exitCode);
    }

    @CommandLine.Option(names = {"-V", "--verbose"}, description = "show verbose output", scope = CommandLine.ScopeType.INHERIT)
    public void configureLogger(boolean verbose) {
        Logger rootLogger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
        rootLogger.setLevel(verbose ? Level.INFO : Level.ERROR);
    }

    @Override
    public void run() {
        logger.trace("Starting... (trace) from app");
        logger.debug("Starting... (debug) from app");
        logger.info("Starting... (info)  from app");
        logger.warn("Starting... (warn)  from app");
        logger.error("Starting... (error) from app");
    }
}
