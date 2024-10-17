package com.inspiretmstech.supavault;

import com.inspiretmstech.supavault.bases.Loggable;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

import java.sql.SQLException;

public class ExceptionHandler extends Loggable implements CommandLine.IExecutionExceptionHandler {

    public ExceptionHandler() {
        super(ExceptionHandler.class);
    }

    public static void handle(Exception ex) {
        Logger logger = LoggerFactory.getLogger(ExceptionHandler.class);
        try {
            throw ex;
        } catch (RuntimeException e) {
            logger.error(e.getMessage());
        } catch (Exception e) {
            logger.error("A critical error occurred: {}", e.getMessage());
        }
    }

    @Override
    public int handleExecutionException(Exception ex, CommandLine commandLine, CommandLine.ParseResult fullParseResult) {
        ExceptionHandler.handle(ex);
        return 1;
    }

}
