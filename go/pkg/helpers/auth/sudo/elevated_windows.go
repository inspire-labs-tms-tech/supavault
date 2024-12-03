//go:build windows
// +build windows

package sudo

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

// IsElevated checks if the program is running with Administrator privileges on Windows.
func IsElevated() bool {
	var token windows.Token
	if err := windows.OpenCurrentProcessToken(&token); err != nil {
		return false
	}
	defer token.Close()

	// Check if the token has the elevated privilege
	var elevation windows.TokenElevation
	var size uint32
	if err := windows.GetTokenInformation(token, windows.TokenElevation, &elevation, uint32(unsafe.Sizeof(elevation)), &size); err != nil {
		return false
	}

	return elevation.TokenIsElevated != 0
}
