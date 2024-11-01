package com.inspiretmstech.supavault.utils.secrets.shims;

import java.util.Optional;

public interface SecretsShim {

    Optional<String> getSecret();

    void setSecret(String secret);

    void removeSecret();
}
