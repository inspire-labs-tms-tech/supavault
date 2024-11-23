package com.inspiretmstech.supavault.models;

public record ClientAuth(
        String id,
        String secret,
        String url,
        String anonKey
) {
}
