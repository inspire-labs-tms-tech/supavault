package com.inspiretmstech.supavault.utils.supabase;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.models.ClientAuth;
import com.inspiretmstech.supavault.utils.gson.GSON;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.List;

public class SupabaseClient extends Loggable implements AutoCloseable {

    private final ClientAuth clientAuth;
    private SupabaseAuthResponse auth;

    public SupabaseClient(ClientAuth auth) throws IOException {
        super(SupabaseClient.class);
        this.clientAuth = auth;
        if (!this.authenticate()) logger.error("authentication failure");
    }

    public boolean authenticate() throws IOException {

        this.logger.debug("authenticating supabase client");

        // Create the JSON payload
        JsonObject payload = new JsonObject();
        payload.addProperty("email", this.clientAuth.id());
        payload.addProperty("password", this.clientAuth.secret());

        // Open HTTP connection
        URL url = new URL(this.getBaseURL() + "/auth/v1/token?grant_type=password");
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Content-Type", "application/json");
        conn.setRequestProperty("Authorization", "Bearer " + this.clientAuth.anonKey());
        conn.setDoOutput(true);

        // Send request
        try (OutputStream os = conn.getOutputStream()) {
            os.write(GSON.GLOBAL.toJson(payload).getBytes(StandardCharsets.UTF_8));
            os.flush();
        }

        // Read the response
        int responseCode = conn.getResponseCode();
        InputStream inputStream;
        if (responseCode >= 200 && responseCode < 300) inputStream = conn.getInputStream();
        else inputStream = conn.getErrorStream();

        // Read the response body
        if (inputStream != null) {
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(inputStream, StandardCharsets.UTF_8))) {
                StringBuilder response = new StringBuilder();
                String line;
                while ((line = reader.readLine()) != null) {
                    response.append(line);
                }

                if (responseCode >= 200 && responseCode < 300) {
                    this.auth = SupabaseAuthResponse.fromJson(response.toString());
                    this.logger.debug("authenticated successfully as user {}", this.auth.user().id());
                } else {
                    this.logger.error("error response during supabase authentication: {}", response);
                }
            }
        }

        return responseCode == HttpURLConnection.HTTP_OK;
    }

    /**
     * Generic function to retrieve all rows from a table and map them to instances of a specified class.
     *
     * @param tableName   the name of the table
     * @param recordClass the class to map rows into
     * @param <T>         the type of the record class
     * @return a list of objects of type T
     * @throws Exception if an error occurs during the request
     */
    public <T> List<T> get(String tableName, Class<T> recordClass) throws Exception {

        String endpoint = this.getBaseURL() + "/rest/v1/" + tableName;
        HttpURLConnection conn = (HttpURLConnection) (new URL(endpoint)).openConnection();
        conn.setRequestMethod("GET");
        conn.setRequestProperty("apikey", this.clientAuth.anonKey());
        conn.setRequestProperty("Authorization", "Bearer " + this.auth.access_token());
        conn.setRequestProperty("Content-Type", "application/json");

        int responseCode = conn.getResponseCode();
        if (responseCode == HttpURLConnection.HTTP_OK)
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(conn.getInputStream()))) {
                StringBuilder response = new StringBuilder();
                String line;
                while ((line = reader.readLine()) != null) response.append(line);
                return JsonParser
                        .parseString(response.toString())
                        .getAsJsonArray()
                        .asList()
                        .stream()
                        .map(row -> GSON.GLOBAL.fromJson(row, recordClass))
                        .toList();
            }
        else {
            this.logger.error("failed to retrieve rows from table {} (HTTP response code: {})", tableName, responseCode);
            return List.of();
        }
    }

    private String getBaseURL() {
        String base = this.clientAuth.url();
        if (base.endsWith("/")) base = base.substring(0, base.length() - 1);
        return base;
    }

    /**
     * Close the client by invalidating the local session
     */
    @Override
    public void close() throws Exception {
        this.logger.debug("closing supabase client");

        URL url = new URL(this.getBaseURL() + "/auth/v1/logout");
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Authorization", "Bearer " + this.auth.access_token());
        conn.setRequestProperty("apikey", this.clientAuth.anonKey());

        int responseCode = conn.getResponseCode();
        if (responseCode == HttpURLConnection.HTTP_NO_CONTENT) {
            logger.debug("supabase client closed successfully");
        } else {
            logger.debug("failed invalidate local session (response code: {})", responseCode);
        }
    }
}
