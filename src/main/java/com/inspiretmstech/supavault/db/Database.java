package com.inspiretmstech.supavault.db;

import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.env.RuntimeEnvironment;
import org.jooq.DSLContext;
import org.jooq.impl.DSL;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class Database extends Loggable {

    @FunctionalInterface
    public interface Executor {
        void execute(DSLContext db) throws SQLException;
    }

    private static Database instance;

    private final Connection conn;

    private final DSLContext ctx;

    private String getConnectionString() {
        return "jdbc:postgresql://" + RuntimeEnvironment.getDbHost() + ":" + RuntimeEnvironment.getDbPort() + "/" + RuntimeEnvironment.getDbName();
    }

    private Database() throws SQLException {
        super(Database.class);
        this.conn = DriverManager.getConnection(this.getConnectionString(), RuntimeEnvironment.getDbUser(), RuntimeEnvironment.getDbPass());
        this.conn.isValid(5);
        this.ctx = DSL.using(this.conn);
    }

    public static synchronized Database get() throws SQLException {
        if (instance == null) instance = new Database();
        return instance;
    }

    public static Database with() {
        try {
            return Database.get();
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

    public void unsafely(Executor executor) throws Exception {
        executor.execute(this.ctx);
    }

    public void execute(Executor executor) {
        try {
            executor.execute(this.ctx);
        } catch (SQLException e) {
            logger.error("Database Exception: {}", e.getMessage());
            for (StackTraceElement stackTraceElement : e.getStackTrace())
                logger.debug(stackTraceElement.toString());
        }
    }

}
