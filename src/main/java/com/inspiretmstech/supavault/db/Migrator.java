package com.inspiretmstech.supavault.db;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.inspiretmstech.supavault.db.gen.supabase_migrations.Tables;
import com.inspiretmstech.supavault.db.gen.supabase_migrations.tables.records.SchemaMigrationsRecord;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

public class Migrator {

    private final static Logger logger = LoggerFactory.getLogger(Migrator.class);
    // SEE: https://github.com/supabase/cli/blob/8dbf6b2750ad34f3645f1dbc19fdab0afc905da8/pkg/migration/history.go
    private final static List<String> SETUP_MIGRATIONS_IF_NOT_EXISTS = List.of(
            "CREATE SCHEMA IF NOT EXISTS supabase_migrations",
            "CREATE TABLE IF NOT EXISTS supabase_migrations.schema_migrations (version text NOT NULL PRIMARY KEY)",
            "ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS statements text[]",
            "ALTER TABLE supabase_migrations.schema_migrations ADD COLUMN IF NOT EXISTS name text",
            "CREATE TABLE IF NOT EXISTS supabase_migrations.seed_files (path text NOT NULL PRIMARY KEY, hash text NOT NULL)"
    );

    public static void migrate(String version) throws IOException, InterruptedException {
        String downloadUrl = "https://github.com/inspire-labs-tms-tech/supavault/releases/download/" + version + "/migrations.zip";

        byte[] zipData = Migrator.downloadZip(downloadUrl);

        Migrator.processZip(zipData);
    }

    public static void migrate() throws IOException, InterruptedException {
        String latestVersion = Migrator.getLatestRelease();
        Migrator.migrate(latestVersion);
    }

    private static byte[] downloadZip(String url) throws IOException, InterruptedException {
        // Create an HttpClient that follows redirects
        HttpClient client = HttpClient.newBuilder()
                .followRedirects(HttpClient.Redirect.ALWAYS)
                .build();

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .GET()
                .build();

        HttpResponse<byte[]> response = client.send(request, HttpResponse.BodyHandlers.ofByteArray());

        if (response.statusCode() == 200) {
            return response.body();
        } else {
            throw new IOException("Failed to download the ZIP file: HTTP " + response.statusCode());
        }
    }

    private static void processZip(byte[] zipData) throws IOException {
        List<String> sortedFileNames = new ArrayList<>();
        Map<String, String> fileContents = new HashMap<>();

        // Step 1: Read all entries into a map
        try (ZipInputStream zipStream = new ZipInputStream(new ByteArrayInputStream(zipData))) {
            ZipEntry entry;
            while ((entry = zipStream.getNextEntry()) != null) {
                if (!entry.isDirectory()) {
                    String fileName = entry.getName();
                    sortedFileNames.add(fileName);

                    // Read file contents into memory
                    String content = new String(zipStream.readAllBytes(), StandardCharsets.UTF_8);
                    fileContents.put(fileName, content);
                }
                zipStream.closeEntry();
            }
        }

        // Step 2: Sort filenames
        sortedFileNames.sort(Comparator.naturalOrder());

        // setup migrations if not exists
        Database.with().execute(db -> {
            for (String command : SETUP_MIGRATIONS_IF_NOT_EXISTS)
                logger.debug("setup migrations history table command: {}", db.execute(command));
        });

        // load existing migrations
        Map<String, SchemaMigrationsRecord> migrations = new HashMap<>();
        Database.with().execute(db -> {
            for (SchemaMigrationsRecord record : db.selectFrom(Tables.SCHEMA_MIGRATIONS)
                    .fetch().stream().toList())
                migrations.put(record.getVersion(), record);
        });

        // Step 3: Process files in sorted order
        for (String fileName : sortedFileNames) {

            String[] pathParts = fileName.split("/");
            String[] parts = Arrays.stream(pathParts).toList().getLast().split("_", 2);
            String version = parts[0];
            String description = parts[1].replace(".sql", "");
            String sql = fileContents.get(fileName);

            if (migrations.containsKey(version)) {
                logger.info("migration {}: already applied", version);
            } else {
                SchemaMigrationsRecord migration = new SchemaMigrationsRecord();
                migration.setVersion(version);
                migration.setName(description);
                migration.setStatements(parseMigrationFile(sql));

                Database.with().execute(db -> {
                    db.transaction(tx -> {
                        tx.dsl().execute(sql);
                        tx.dsl().insertInto(Tables.SCHEMA_MIGRATIONS).set(migration).execute();
                    });
                });
                logger.info("migration {}: applied", version);
            }
        }
    }

    private static String[] parseMigrationFile(String migrationContents) {
        List<String> statements = new ArrayList<>();

        // Split the content into individual SQL statements by semicolon
        String[] sqlStatements = migrationContents.split(";");

        for (String statement : sqlStatements) {
            // Trim whitespace from each statement
            statement = statement.trim();

            // Skip empty or blank lines
            if (!statement.isEmpty()) {
                statements.add(statement);
            }
        }

        return statements.toArray(new String[0]);
    }

    private static String getLatestRelease() throws IOException, InterruptedException {
        String apiUrl = "https://api.github.com/repos/inspire-labs-tms-tech/supavault/releases/latest";

        HttpClient client = HttpClient.newHttpClient();
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(apiUrl))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());

        if (response.statusCode() == 200) {
            // Parse the JSON response to get the tag name using GSON
            JsonObject json = JsonParser.parseString(response.body()).getAsJsonObject();
            return json.get("tag_name").getAsString();
        } else {
            throw new IOException("Failed to fetch the latest release: HTTP " + response.statusCode());
        }
    }

}
