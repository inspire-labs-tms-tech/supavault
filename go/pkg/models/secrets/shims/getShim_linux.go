//go:build linux

package shims

func GetShim() SecretShim {
	return &UbuntuSecretsShim{}
}
