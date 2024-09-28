package com.inspiretmstech.supavault;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

@CommandLine.Command(name = "app", subcommands = {})
class Main implements Runnable {

    static {
        System.setProperty("slf4j.internal.verbosity", "WARN");
        ((Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME)).setLevel(Level.ERROR);
    }

    private static Logger logger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);

    @CommandLine.Option(names = {"-v"}, scope = CommandLine.ScopeType.INHERIT)
    public void configureLogger(boolean[] verbose) {
        // Get the root logger
        Logger rootLogger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
        rootLogger.setLevel(verbose.length > 0 ? Level.INFO : Level.ERROR);
    }

    public static void main(String[] args) {
        int exitCode = new CommandLine(new Main())
                .execute(args);
        System.exit(exitCode);
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
