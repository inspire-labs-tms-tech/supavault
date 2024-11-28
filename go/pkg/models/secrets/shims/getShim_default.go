//go:build !darwin && !linux

package shims

func GetShim() shims.SecretShim {
	panic(runtime.GOOS + " secrets shim not implemented")
}
