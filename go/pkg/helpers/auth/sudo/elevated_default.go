//go:build !windows && !unix

package sudo

// IsElevated stub for IDE analysis
func IsElevated() bool {
	panic("stub implementation of IsElevated being used in build")
}
