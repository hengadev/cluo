package updater

import "runtime"

// GetPlatform returns the current platform identifier in the format "os_arch".
// Examples: "linux_amd64", "darwin_arm64", "windows_amd64"
func GetPlatform() string {
	return runtime.GOOS + "_" + runtime.GOARCH
}
