package com.inspiretmstech.supavault.db;

import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.env.RuntimeEnvironment;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.SQLException;

public class Database extends Loggable {

    private static Database instance;

    private final Connection conn;

    private String getConnectionString() {
        return "jdbc:postgresql://" + RuntimeEnvironment.getDbHost() + ":" + RuntimeEnvironment.getDbPort() + "/" + RuntimeEnvironment.getDbName();
    }

    private Database() throws SQLException {
        super(Database.class);
        this.conn = DriverManager.getConnection(this.getConnectionString(), RuntimeEnvironment.getDbUser(), RuntimeEnvironment.getDbPass());
        this.conn.isValid(5);
    }

    public static synchronized Database get() throws SQLException {
        if (instance == null) instance = new Database();
        return instance;
    }

}
