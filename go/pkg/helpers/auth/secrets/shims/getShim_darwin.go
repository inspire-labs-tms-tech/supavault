//go:build darwin

package shims

func GetShim() SecretShim {
	return &MacSecretsShim{}
}
