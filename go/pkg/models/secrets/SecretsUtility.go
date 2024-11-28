package secrets

import (
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/secrets/shims"
)

func SetSecret(credentials shims.ClientCredentials) error {
	return shims.GetShim().SetSecret(credentials)
}

func GetSecret() (shims.ClientCredentials, error) {
	return shims.GetShim().GetSecret()
}
