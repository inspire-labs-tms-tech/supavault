package supabase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/auth/secrets"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/auth"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/supabase"
	"io"
	"net/http"
	"strings"
)

type SupabaseClient struct {
	credentials auth.ClientCredentials
	token       *supabase.AuthResponse
}

func GetClient() (SupabaseClient, error) {

	secret, err := secrets.GetSecret()
	if err != nil {
		return SupabaseClient{}, err
	}

	return SupabaseClient{
		credentials: secret,
		token:       nil,
	}, nil
}

// Authenticate authenticates the Supabase client.
func (c *SupabaseClient) Authenticate() (bool, error) {

	// Prepare the JSON payload
	payload := map[string]string{
		"email":    c.credentials.Email,
		"password": c.credentials.Password,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, supabase.NewClientError("error creating authentication payload", err)
	}

	// Open HTTP connection
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", c.getBaseURL())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, supabase.NewClientError("error creating request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.credentials.AnonKey)
	req.Header.Set("apikey", c.credentials.AnonKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, supabase.NewClientError("error during authentication request", err)
	}
	defer resp.Body.Close()

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, supabase.NewClientError("error reading response body", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := json.Unmarshal(body, &c.token); err != nil {
			return false, supabase.NewClientError("error parsing authentication response", err)
		}
		return true, nil
	} else {
		return false, supabase.NewClientError("error response during supabase authentication", errors.New(string(body)))
	}
}

// Close invalidates the local session.
func (c *SupabaseClient) Close() error {
	url := fmt.Sprintf("%s/auth/v1/logout", c.getBaseURL())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return supabase.NewClientError("error creating logout request", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.credentials.AnonKey)
	req.Header.Set("apikey", c.credentials.AnonKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return supabase.NewClientError("error during logout request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return supabase.NewClientError("logout failed", errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode)))
}

// getBaseURL returns the base URL of the Supabase client.
func (c *SupabaseClient) getBaseURL() string {
	return strings.TrimSuffix(c.credentials.Url, "/")
}
