package com.inspiretmstech.supavault.commands.auth;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
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

            if(rawSecret.isEmpty()) {
                logger.error("not logged in!");
                return;
            }

            logger.debug(rawSecret.get());
            ClientAuth auth = GSON.GLOBAL.fromJson(rawSecret.get(), ClientAuth.class);

            logger.info("         Client ID: {}", auth.id());
            logger.info("     Client Secret: {}", auth.secret().replaceAll(".", "*"));
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
            @CommandLine.Option(names = {"--client-id"}, required = true, description = "id of the client to authenticate with") UUID clientID,
            @CommandLine.Option(names = {"--client-secret"}, required = true, description = "secret of the client to authenticate with") String clientSecret
    ) {
        logger.debug("setting login credentials");

        ClientAuth auth = new ClientAuth(clientID.toString(), clientSecret);
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
