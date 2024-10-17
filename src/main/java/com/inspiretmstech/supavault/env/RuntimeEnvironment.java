package com.inspiretmstech.supavault.env;

import java.util.Optional;

public class RuntimeEnvironment {

    private static String DB_HOST = Optional.ofNullable(System.getenv("DB_HOST")).orElse("127.0.0.1");
    private static String DB_PORT = Optional.ofNullable(System.getenv("DB_PORT")).orElse("54322");
    private static String DB_NAME = Optional.ofNullable(System.getenv("DB_NAME")).orElse("postgres");
    private static String DB_USER = Optional.ofNullable(System.getenv("DB_USER")).orElse("postgres");
    private static String DB_PASS = Optional.ofNullable(System.getenv("DB_PASS")).orElse("postgres");

    public static String _toString() {
        return "RuntimeEnvironment{" +
               "DB_HOST='" + DB_HOST + '\'' +
               ", DB_PORT='" + DB_PORT + '\'' +
               ", DB_NAME='" + DB_NAME + '\'' +
               ", DB_USER='" + DB_USER + '\'' +
               ", DB_PASS='" + DB_PASS + '\'' +
               "}";
    }

    public static String getDbHost() {
        return DB_HOST;
    }

    public static String getDbPort() {
        return DB_PORT;
    }

    public static String getDbName() {
        return DB_NAME;
    }

    public static String getDbUser() {
        return DB_USER;
    }

    public static String getDbPass() {
        return DB_PASS;
    }

    public static synchronized void setDbHost(String dbHost) {
        DB_HOST = dbHost;
    }

    public static synchronized void setDbPort(String dbPort) {
        DB_PORT = dbPort;
    }

    public static synchronized void setDbName(String dbName) {
        DB_NAME = dbName;
    }

    public static synchronized void setDbUser(String dbUser) {
        DB_USER = dbUser;
    }

    public static synchronized void setDbPass(String dbPass) {
        DB_PASS = dbPass;
    }
}
