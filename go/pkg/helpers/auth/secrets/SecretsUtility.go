package secrets

import (
	"github.com/inspire-labs-tms-tech/supavault/pkg/helpers/auth/secrets/shims"
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/auth"
)

func SetSecret(credentials auth.ClientCredentials) error {
	return shims.GetShim().SetSecret(credentials)
}

func GetSecret() (auth.ClientCredentials, error) {
	return shims.GetShim().GetSecret()
}

func RemoveSecret() error {
	return shims.GetShim().RemoveSecret()
}
