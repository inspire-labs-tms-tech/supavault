package com.inspiretmstech.supavault.utils.secrets.shims;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.Optional;

public class UbuntuSecretsShim implements SecretsShim {

    private static final String ACCOUNT = "supavault@local";
    private static final String CLIENT_SECRET_KEY = "com.inspiretmstech.supavault";

    @Override
    public Optional<String> getSecret() {
        ProcessBuilder builder = new ProcessBuilder(
                "secret-tool",
                "lookup",
                "account", ACCOUNT, // account identifier
                "service", CLIENT_SECRET_KEY // unique service name
        );

        Process process;
        String secret;
        try {
            process = builder.start();
            BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
            secret = reader.readLine();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        return Optional.ofNullable(secret);
    }

    @Override
    public void setSecret(String secret) {
        ProcessBuilder builder = new ProcessBuilder(
                "secret-tool",
                "store",
                "--label=Supavault Secret",
                "account", ACCOUNT, // account identifier
                "service", CLIENT_SECRET_KEY // unique service name
        );

        Process process;
        try {
            process = builder.start();
            process.getOutputStream().write(secret.getBytes()); // Write the secret to the process
            process.getOutputStream().close();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        if (process.exitValue() != 0) {
            throw new RuntimeException("Failed to store secret in Gnome Keyring");
        }
    }

    @Override
    public void removeSecret() {
        ProcessBuilder builder = new ProcessBuilder(
                "secret-tool",
                "clear",
                "account", ACCOUNT, // account identifier
                "service", CLIENT_SECRET_KEY // unique service name
        );

        Process process;
        try {
            process = builder.start();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        if (process.exitValue() != 0) {
            throw new RuntimeException("Failed to delete secret from Gnome Keyring");
        }
    }
}
