package supabase

import "encoding/json"

// AuthResponse represents the main authentication response.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// FromJSON parses a JSON string into a AuthResponse struct.
func (r *AuthResponse) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), r)
}

// User represents a user within the authentication response.
type User struct {
	ID               string                 `json:"id"`
	Aud              string                 `json:"aud"`
	Role             string                 `json:"role"`
	Email            string                 `json:"email"`
	EmailConfirmedAt string                 `json:"email_confirmed_at"`
	Phone            string                 `json:"phone"`
	ConfirmedAt      string                 `json:"confirmed_at"`
	LastSignInAt     string                 `json:"last_sign_in_at"`
	AppMetadata      AppMetadata            `json:"app_metadata"`
	UserMetadata     map[string]interface{} `json:"user_metadata"`
	Identities       []Identity             `json:"identities"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
	IsAnonymous      bool                   `json:"is_anonymous"`
}

// AppMetadata represents application-specific metadata for a user.
type AppMetadata struct {
	Provider  string   `json:"provider"`
	Providers []string `json:"providers"`
}

// Identity represents an identity provider for a user.
type Identity struct {
	IdentityID   string       `json:"identity_id"`
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	IdentityData IdentityData `json:"identity_data"`
	Provider     string       `json:"provider"`
	LastSignInAt string       `json:"last_sign_in_at"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	Email        string       `json:"email"`
}

// IdentityData represents data associated with a specific identity provider.
type IdentityData struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`
	Sub           string `json:"sub"`
}
