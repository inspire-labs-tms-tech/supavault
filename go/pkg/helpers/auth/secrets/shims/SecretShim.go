package shims

import "github.com/inspire-labs-tms-tech/supavault/pkg/models/auth"

type SecretShim interface {
	GetSecret() (auth.ClientCredentials, error)

	SetSecret(credentials auth.ClientCredentials) error

	RemoveSecret() error
}
