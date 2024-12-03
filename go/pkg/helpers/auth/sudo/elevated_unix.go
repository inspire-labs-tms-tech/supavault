//go:build !windows

package sudo

import "os"

// IsElevated checks if the program is running with root privileges on UNIX-like systems.
func IsElevated() bool {
	return os.Geteuid() == 0
}
