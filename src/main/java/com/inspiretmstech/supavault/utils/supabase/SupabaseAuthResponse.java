package com.inspiretmstech.supavault.utils.supabase;

import com.inspiretmstech.supavault.utils.gson.GSON;

import java.util.List;
import java.util.Map;

public record SupabaseAuthResponse(
        String access_token,
        String token_type,
        int expires_in,
        long expires_at,
        String refresh_token,
        User user
) {

    public static SupabaseAuthResponse fromJson(String json) {
        return GSON.GLOBAL.fromJson(json, SupabaseAuthResponse.class);
    }

    public record User(
            String id,
            String aud,
            String role,
            String email,
            String email_confirmed_at,
            String phone,
            String confirmed_at,
            String last_sign_in_at,
            AppMetadata app_metadata,
            Map<String, Object> user_metadata,
            List<Identity> identities,
            String created_at,
            String updated_at,
            boolean is_anonymous
    ) {
        public record AppMetadata(
                String provider,
                List<String> providers
        ) {}

        public record Identity(
                String identity_id,
                String id,
                String user_id,
                IdentityData identity_data,
                String provider,
                String last_sign_in_at,
                String created_at,
                String updated_at,
                String email
        ) {
            public record IdentityData(
                    String email,
                    boolean email_verified,
                    boolean phone_verified,
                    String sub
            ) {}
        }
    }
}
