package com.inspiretmstech.supavault.commands.auth;

import com.inspiretmstech.supavault.bases.Loggable;
import com.inspiretmstech.supavault.models.ClientAuth;
import com.inspiretmstech.supavault.utils.gson.GSON;
import com.inspiretmstech.supavault.utils.secrets.SecretsManager;
import picocli.CommandLine;

import java.util.Optional;
import java.util.UUID;

@CommandLine.Command(
        name = "auth",
        description = "store default login information"
)
public class Auth extends Loggable {

    public Auth() {
        super(Auth.class);
    }

    @CommandLine.Command(
            name = "show",
            description = "show the default login credentials"
    )
    public void show() {
        logger.debug("showing login credentials");
        SecretsManager sm = new SecretsManager();
        try {
            Optional<String> rawSecret = sm.getSecret();

            if (rawSecret.isEmpty()) {
                logger.error("not logged in!");
                return;
            }

            logger.debug(rawSecret.get());
            ClientAuth auth = GSON.GLOBAL.fromJson(rawSecret.get(), ClientAuth.class);

            logger.info("         Client ID: {}", auth.id());
            logger.info("     Client Secret: {}", auth.secret().replaceAll(".", "*"));
            logger.info("      Supabase URL: {}", auth.url());
            logger.info(" Supabase Anon Key: {}", auth.anonKey());
        } catch (Exception e) {
            logger.error("unable to retrieve secret (are you logged in?)");
        }
    }

    @CommandLine.Command(
            name = "logout",
            description = "remove the default login credentials"
    )
    public void logout() {
        logger.debug("clearing login credentials");
        SecretsManager sm = new SecretsManager();
        try {
            sm.removeSecret();
            logger.info("successfully logged out");
        } catch (Exception e) {
            logger.error("unable to log out (were you logged in?)");
        }
    }

    @CommandLine.Command(
            name = "login",
            description = "store the default login credentials"
    )
    public void login(
            @CommandLine.Option(names = {"--client-id"}, required = true, description = "id (email) of the client to authenticate with") String clientID,
            @CommandLine.Option(names = {"--client-secret"}, required = true, description = "secret of the client to authenticate with") String clientSecret,
            @CommandLine.Option(names = {"--url"}, required = true, description = "URL to supabase api (ex.: https://<project-id>.supabase.com)", defaultValue = "http://127.0.0.1:54321") String url,
            @CommandLine.Option(names = {"--anon"}, required = true, description = "Anon key to supabase instance", defaultValue = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0") String anonKey
    ) {
        logger.debug("setting login credentials");

        ClientAuth auth = new ClientAuth(clientID.toString(), clientSecret, url, anonKey);
        SecretsManager sm = new SecretsManager();
        try {
            String serialized = GSON.GLOBAL.toJson(auth);
            logger.debug("serializing: {}", serialized);
            sm.setSecret(serialized);
            logger.info("successfully logged in");
        } catch (Exception e) {
            logger.error("unable to log in");
        }
    }

}
