package secrets

import (
	"github.com/inspire-labs-tms-tech/supavault/pkg/models/secrets/shims"
	"runtime"
)

func getShim() shims.SecretShim {
	os := runtime.GOOS
	switch os {
	case "darwin": // macOS
		return &shims.MacSecretsShim{}
	case "linux":
		return &shims.UbuntuSecretsShim{}
	case "windows":
		panic("windows secrets shim not implemented")
	default:
		panic(os + " secrets shim not implemented")
	}
}

func SetSecret(credentials shims.ClientCredentials) error {
	return getShim().SetSecret(credentials)
}

func GetSecret() (shims.ClientCredentials, error) {
	return getShim().GetSecret()
}
