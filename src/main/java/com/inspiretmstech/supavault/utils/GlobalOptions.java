package com.inspiretmstech.supavault.utils;


import ch.qos.logback.classic.Logger;
import com.inspiretmstech.supavault.constants.LogLevel;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

public class GlobalOptions {

    @CommandLine.Option(names = {"-V", "--verbose"}, description = "show verbose output", scope = CommandLine.ScopeType.INHERIT)
    public void configureLogger(boolean verbose) {
        Logger rootLogger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
        rootLogger.setLevel(verbose ? LogLevel.VERBOSE : LogLevel.DEFAULT);
    }

}
