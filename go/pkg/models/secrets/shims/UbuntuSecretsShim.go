package shims

import (
	"bytes"
	"encoding/json"
	"errors"
	errs "github.com/inspire-labs-tms-tech/supavault/pkg/models/secrets/errors"
	"log"
	"os/exec"
)

type UbuntuSecretsShim struct {
}

func (u *UbuntuSecretsShim) GetSecret() (ClientCredentials, error) {

	cmd := exec.Command(
		"secret-tool", "lookup",
		"account", Account,
		"service", Service,
	)

	output, err := cmd.Output()
	if err != nil {
		return ClientCredentials{}, errors.New("failed to retrieve secret from Gnome Keyring")
	}

	secret := string(bytes.TrimSpace(output))
	if secret == "" {
		return ClientCredentials{}, errs.NewNotFoundError()
	}

	var client ClientCredentials
	if json.Unmarshal([]byte(secret), &client) != nil {
		return ClientCredentials{}, err
	}
	return client, nil
}

func (u *UbuntuSecretsShim) SetSecret(client ClientCredentials) error {

	// remove existing secret, if there is one
	u.RemoveSecret()

	serialized, err := json.Marshal(client)
	if err != nil {
		log.Fatalf("error serializing secret: %v", err)
	}

	cmd := exec.Command(
		"secret-tool", "store",
		"account", Account,
		"service", Service,
	)

	// Provide the secret as input to the command.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.New("failed to create stdin pipe for storing secret")
	}

	// Start the command execution.
	if err := cmd.Start(); err != nil {
		return errors.New("failed to start secret-tool process")
	}

	// Write the secret to the command's input.
	if _, err := stdin.Write(serialized); err != nil {
		return errors.New("failed to write secret to secret-tool")
	}

	// Close the input pipe.
	if err := stdin.Close(); err != nil {
		return errors.New("failed to close stdin pipe for secret-tool")
	}

	// Wait for the command to complete.
	if err := cmd.Wait(); err != nil {
		return errors.New("failed to store secret in Gnome Keyring")
	}

	// Check the exit code of the process.
	if cmd.ProcessState.ExitCode() != 0 {
		return errors.New("secret-tool returned a non-zero exit code")
	}

	return nil
}

func (u *UbuntuSecretsShim) RemoveSecret() error {
	clearCmd := exec.Command(
		"secret-tool", "clear",
		"account", Account,
		"service", Service,
	)
	_ = clearCmd.Run()
	return nil
}
