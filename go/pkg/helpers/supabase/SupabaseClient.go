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

func GetClient() (*SupabaseClient, error) {

	secret, err := secrets.GetSecret()
	if err != nil {
		return nil, err
	}

	return &SupabaseClient{
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

func (c *SupabaseClient) IsAuthenticated() bool {
	return c.token != nil
}

func (c *SupabaseClient) GetUser() (*supabase.User, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("supabase client not authenticated")
	}

	return &c.token.User, nil
}

// Get retrieves rows from the specified table endpoint and maps them to the provided type
func (c *SupabaseClient) Get(tableName string, output interface{}) error {

	if !c.IsAuthenticated() {
		return errors.New("supabase client not authenticated")
	}

	// Build the endpoint URL
	endpoint := fmt.Sprintf("%s/rest/v1/%s", c.getBaseURL(), tableName)

	// Create the HTTP request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("apikey", c.credentials.AnonKey)
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to retrieve rows from table %s (HTTP response code: %d): %s", tableName, resp.StatusCode, string(body))
	}

	// Read the entire response body into a buffer
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert the body to a string and print/log it
	//bodyString := string(bodyBytes)
	//fmt.Println("Response Body:", bodyString)

	// Decode the body into the output object
	if err := json.Unmarshal(bodyBytes, &output); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// getBaseURL returns the base URL of the Supabase client.
func (c *SupabaseClient) getBaseURL() string {
	return strings.TrimSuffix(c.credentials.Url, "/")
}
