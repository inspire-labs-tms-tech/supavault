package com.inspiretmstech.supavault.commands;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.commands.admin.Admin;
import com.inspiretmstech.supavault.commands.auth.Auth;
import com.inspiretmstech.supavault.commands.projects.Projects;
import com.inspiretmstech.supavault.constants.LogLevel;
import com.inspiretmstech.supavault.constants.Version;
import com.inspiretmstech.supavault.db.gen.public_.Tables;
import com.inspiretmstech.supavault.db.gen.public_.tables.records.EnvironmentVariablesRecord;
import com.inspiretmstech.supavault.db.gen.public_.tables.records.VariablesRecord;
import com.inspiretmstech.supavault.models.ClientAuth;
import com.inspiretmstech.supavault.utils.ConnectionSettings;
import com.inspiretmstech.supavault.utils.GlobalOptions;
import com.inspiretmstech.supavault.utils.secrets.SecretsManager;
import com.inspiretmstech.supavault.utils.supabase.SupabaseClient;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

import java.io.*;
import java.util.*;
import java.util.concurrent.Callable;
import java.util.stream.Collectors;


@CommandLine.Command(
        name = "supavault",
        subcommands = {
                Projects.class,
                Auth.class,
                Admin.class,
        },
        versionProvider = Version.class,
        mixinStandardHelpOptions = true
)
public class Supavault extends Loggable implements Callable<Integer> {

    // setup logging
    static {
        System.setProperty("org.jooq.no-logo", "true");
        System.setProperty("org.jooq.no-tips", "true");
        System.setProperty("slf4j.internal.verbosity", "OFF");
        ((Logger) LoggerFactory.getLogger("org.jooq")).setLevel(Level.OFF);
        ((Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME)).setLevel(LogLevel.DEFAULT);
    }

    // inject the DB connection settings
    @CommandLine.ArgGroup(exclusive = false)
    ConnectionSettings settings;
    // inject global options
    @CommandLine.ArgGroup(exclusive = false)
    GlobalOptions globalOptions;
    @CommandLine.Option(names = {"-v", "--version"}, versionHelp = true,
            description = "print version information and exit")
    boolean versionRequested;
    @CommandLine.Parameters
    List<String> commands;
    @CommandLine.Option(names = {"-e", "--env"}, description = "a .env file to load additional environment variables from")
    String env;

    public Supavault() {
        super(Supavault.class);
    }

    @Override
    public Integer call() {

        if (Objects.isNull(commands) || commands.isEmpty()) {
            logger.error("enter a command (ex.: supavault -- echo 'Environment:' && printenv)");
            return 1;
        }

        logger.debug("command: {}", commands);

        Map<String, String> environment = new HashMap<>();

        // attempt to load .env file vars
        if (Objects.isNull(env)) logger.debug("no env file provided");
        else {
            logger.debug("attempting to load variables from '{}'", env.trim());
            try (BufferedReader reader = new BufferedReader(new FileReader(env.trim()))) {
                Map<String, String> envFileEnvironment = this.getEnvironmentFromFile(reader);
                logger.debug("loaded {} variables from '{}'", envFileEnvironment.size(), env.trim());
                environment.putAll(envFileEnvironment);
            } catch (FileNotFoundException e) {
                logger.error("env file '{}' not found", env.trim());
            } catch (IOException e) {
                logger.error("an error occurred while loading .env file '{}': {}", env.trim(), e.getMessage());
                for (StackTraceElement el : e.getStackTrace()) logger.debug(el.toString());
                return 1;
            }
        }


        // get variables from supabase
        logger.debug("attempting to load variables from supabase");
        Map<String, String> remoteEnvironment = this.getEnvironmentFromVault();
        logger.debug("loaded {} variables from supabase", remoteEnvironment.size());
        environment.putAll(remoteEnvironment);

        ProcessBuilder pb = new ProcessBuilder(commands);
        try {

            // merge runtime environments
            Map<String, String> runtimeEnvironment = pb.environment();
            runtimeEnvironment.putAll(environment);

            // start the process
            Process process = pb.start();

            // Read the output
            try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    System.out.println(line);
                }

                int exitCode = process.waitFor();
                logger.debug("Exited with code: {}", exitCode);
                return exitCode;
            }
        } catch (IOException | InterruptedException e) {
            logger.error("an error occurred while executing command: {}", e.getMessage());
            for (StackTraceElement el : e.getStackTrace()) logger.debug(el.toString());
            return 1;
        }
    }

    private Map<String, String> getEnvironmentFromVault() {
        Map<String, String> environment = new HashMap<>();

        SecretsManager sm = new SecretsManager();
        Optional<String> auth = sm.getSecret();
        if (auth.isEmpty()) logger.warn("no supavault credentials found (are you logged in?)");
        else {
            ClientAuth clientAuth = SecretsManager.parseClientAuth(auth.get());
            try (SupabaseClient client = new SupabaseClient(clientAuth)) {

                Map<String, String> defaults = client.get(Tables.VARIABLES)
                        .stream()
                        .collect(Collectors.toMap(
                                VariablesRecord::getId,    // Replace with the actual method to get the key
                                VariablesRecord::getDefault  // Replace with the actual method to get the value
                        ));
                List<EnvironmentVariablesRecord> envVars = client.get(Tables.ENVIRONMENT_VARIABLES);
                for (EnvironmentVariablesRecord envVar : envVars)
                    environment.put(envVar.getVariableId(), Objects.isNull(envVar.getValue()) || envVar.getValue().isEmpty() ? defaults.get(envVar.getVariableId()) : envVar.getValue());
            } catch (Exception e) {
                logger.error("unable to load environment variables from supabase: {}", e.getMessage());
            }
        }
        return environment;
    }

    private Map<String, String> getEnvironmentFromFile(BufferedReader reader) throws IOException {
        Map<String, String> environment = new HashMap<>();
        String line;
        while ((line = reader.readLine()) != null) {
            line = line.trim();
            if (!line.isEmpty() && !line.startsWith("#")) {
                String[] parts = line.split("=", 2);
                if (parts.length == 2) environment.put(parts[0].trim(), parts[1].trim());
            }
        }
        return environment;
    }

}
