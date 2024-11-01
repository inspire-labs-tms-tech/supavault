package com.inspiretmstech.supavault.utils;

import com.inspiretmstech.supavault.env.RuntimeEnvironment;
import picocli.CommandLine;

import java.util.Objects;

public class ConnectionSettings {

    @CommandLine.Option(names = {"-H", "--db-host"}, description = "host of the Postgres database (e.g., 127.0.0.1 or my.db.example.com)", scope = CommandLine.ScopeType.INHERIT)
    private void dbHost(String dbHost) {
        if (Objects.nonNull(dbHost) && !dbHost.isEmpty())
            RuntimeEnvironment.setDbHost(dbHost);
    }

    @CommandLine.Option(names = {"-P", "--db-port"}, description = "port of the Postgres database (e.g., 54321)", scope = CommandLine.ScopeType.INHERIT)
    private void dbPort(String dbPort) {
        if (Objects.nonNull(dbPort) && !dbPort.isEmpty())
            RuntimeEnvironment.setDbPort(dbPort);
    }

    @CommandLine.Option(names = {"-N", "--db-name"}, description = "name of the Postgres database (e.g., postgres)", scope = CommandLine.ScopeType.INHERIT)
    private void dbName(String dbName) {
        if (Objects.nonNull(dbName) && !dbName.isEmpty())
            RuntimeEnvironment.setDbName(dbName);
    }

    @CommandLine.Option(names = {"-u", "--db-user"}, description = "username for the user to authenticate into the Postgres database with (e.g., postgres)", scope = CommandLine.ScopeType.INHERIT)
    private void dbUser(String dbUser) {
        if (Objects.nonNull(dbUser) && !dbUser.isEmpty())
            RuntimeEnvironment.setDbUser(dbUser);
    }

    @CommandLine.Option(names = {"-p", "--db-pass"}, description = "password for the user to authenticate into the Postgres database with (e.g., my$uper$ecretPa$$word!)", scope = CommandLine.ScopeType.INHERIT)
    private void dbPass(String dbPass) {
        if (Objects.nonNull(dbPass) && !dbPass.isEmpty())
            RuntimeEnvironment.setDbPass(dbPass);
    }

}
