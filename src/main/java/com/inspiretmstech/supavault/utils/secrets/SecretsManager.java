package com.inspiretmstech.supavault.utils.secrets;

import com.inspiretmstech.supavault.models.ClientAuth;
import com.inspiretmstech.supavault.utils.gson.GSON;
import com.inspiretmstech.supavault.utils.secrets.shims.MacOSSecretsShim;
import com.inspiretmstech.supavault.utils.secrets.shims.SecretsShim;
import com.inspiretmstech.supavault.utils.secrets.shims.UbuntuSecretsShim;

import java.util.Optional;

public class SecretsManager implements SecretsShim {

    private final SecretsShim shim;

    public SecretsManager() {
        this.shim = switch (System.getProperty("os.name").toLowerCase()) {
            case String os when os.contains("win") -> throw new RuntimeException("Windows Shim Not Implemented!");
            case String os when os.contains("mac") -> new MacOSSecretsShim();
            case String os when os.contains("nix") || os.contains("nux") || os.contains("aix") -> new UbuntuSecretsShim();
            default -> throw new RuntimeException("Unsupported OS System!");
        };
    }

    public static ClientAuth parseClientAuth(String clientAuth) {
        return GSON.GLOBAL.fromJson(clientAuth, ClientAuth.class);
    }

    @Override
    public Optional<String> getSecret() {
        return this.shim.getSecret();
    }

    @Override
    public void setSecret(String secret) {
        this.shim.setSecret(secret);
    }

    @Override
    public void removeSecret() {
        this.shim.removeSecret();
    }
}
