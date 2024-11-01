package com.inspiretmstech.supavault.utils.secrets.shims;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.Optional;

public class MacOSSecretsShim implements SecretsShim {

    private static final String ACCOUNT = "supavault@local";

    private static final String CLIENT_SECRET_KEY = "com.inspiretmstech.supavault";

    @Override
    public Optional<String> getSecret() {
        ProcessBuilder builder = new ProcessBuilder(
                "security",
                "find-generic-password",
                "-a", ACCOUNT, // the ID of the account
                "-s", CLIENT_SECRET_KEY, // the unique service name
                "-w" // output only the password
        );

        Process process;
        String password;
        try {
            process = builder.start();
            BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
            password = reader.readLine();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        return Optional.ofNullable(password);
    }

    @Override
    public void setSecret(String secret) {

        ProcessBuilder builder = new ProcessBuilder(
                "security",
                "add-generic-password",
                "-a", ACCOUNT, // the ID of the account
                "-s", CLIENT_SECRET_KEY, // the unique service name
                "-w", secret, // the secret
                "-U" // update if already exists
        );

        Process process;
        try {
            process = builder.start();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }
        if (process.exitValue() != 0) {
            throw new RuntimeException("Failed to store password in Keychain");
        }
    }

    @Override
    public void removeSecret() {
        ProcessBuilder builder = new ProcessBuilder(
                "security",
                "delete-generic-password",
                "-a", ACCOUNT, // the ID of the account
                "-s", CLIENT_SECRET_KEY // the unique service name
        );

        Process process;
        try {
            process = builder.start();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        if (process.exitValue() != 0) {
            throw new RuntimeException("Failed to delete password from Keychain");
        }
    }


}
