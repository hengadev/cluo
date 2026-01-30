package updater

import (
	"fmt"
	"os"
	"runtime"
)

// InstallBinary replaces the current executable with the new binary.
// Strategy varies by platform:
// - Windows: Rename running exe to .old, copy new to original path
// - macOS/Linux: Backup current, copy new, set permissions, remove backup
func InstallBinary(newBinaryPath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		return installWindows(execPath, newBinaryPath)
	default:
		return installUnix(execPath, newBinaryPath)
	}
}

// installWindows handles binary replacement on Windows.
// Windows allows renaming a running executable but not overwriting it.
func installWindows(execPath, newBinaryPath string) error {
	backupPath := execPath + ".old"

	// Remove old backup if exists
	os.Remove(backupPath)

	// Rename current executable to .old
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to rename current executable: %w", err)
	}

	// Rename new binary to original path
	if err := os.Rename(newBinaryPath, execPath); err != nil {
		// Try to restore the old executable
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Note: The .old file will be cleaned up on next update or can be
	// manually deleted after restart
	return nil
}

// installUnix handles binary replacement on macOS and Linux.
func installUnix(execPath, newBinaryPath string) error {
	backupPath := execPath + ".backup"

	// Get permissions from current executable
	info, err := os.Stat(execPath)
	if err != nil {
		return fmt.Errorf("failed to stat current executable: %w", err)
	}
	mode := info.Mode()

	// Remove old backup if exists
	os.Remove(backupPath)

	// Backup current executable
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current executable: %w", err)
	}

	// Move new binary to executable path
	if err := os.Rename(newBinaryPath, execPath); err != nil {
		// Try to restore backup
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Set executable permissions
	if err := os.Chmod(execPath, mode); err != nil {
		// Not fatal, but log it
		fmt.Printf("warning: failed to set permissions: %v\n", err)
	}

	// Remove backup
	os.Remove(backupPath)

	return nil
}

// GetExecutablePath returns the path to the current executable.
func GetExecutablePath() (string, error) {
	return os.Executable()
}
