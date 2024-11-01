package com.inspiretmstech.supavault;

import com.inspiretmstech.supavault.commands.Supavault;
import com.inspiretmstech.supavault.utils.ExceptionHandler;
import picocli.CommandLine;

class Main {

    public static void main(String[] args) {
        int exitCode;
        try {
            exitCode = new CommandLine(new Supavault())
                    .setExecutionExceptionHandler(new ExceptionHandler())
                    .execute(args);
        } catch (Exception e) {
            ExceptionHandler.handle(e);
            exitCode = 1;
        }
        System.exit(exitCode);
    }

}
