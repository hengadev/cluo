package updater

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// UpdateInfo contains information about an available update.
type UpdateInfo struct {
	Available    bool   `json:"available"`
	Version      string `json:"version"`
	ReleaseNotes string `json:"release_notes"`
	DownloadURL  string `json:"download_url"`
}

// ProgressEvent is emitted during download.
type ProgressEvent struct {
	Downloaded int64   `json:"downloaded"`
	Total      int64   `json:"total"`
	Percent    float64 `json:"percent"`
}

// Updater manages application updates.
type Updater struct {
	ctx           context.Context
	mu            sync.Mutex
	manifest      *Manifest
	downloadPath  string
	cancelDownload context.CancelFunc
}

// NewUpdater creates a new Updater instance.
func NewUpdater() *Updater {
	return &Updater{}
}

// Startup is called when the app starts. Store the context for event emission.
func (u *Updater) Startup(ctx context.Context) {
	u.ctx = ctx
}

// GetCurrentVersion returns the current application version.
func (u *Updater) GetCurrentVersion() string {
	return Version
}

// GetManifestURL returns the configured manifest URL.
func (u *Updater) GetManifestURL() string {
	return ManifestURL
}

// CheckForUpdate checks if an update is available.
func (u *Updater) CheckForUpdate() (UpdateInfo, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	info := UpdateInfo{Available: false}

	if ManifestURL == "" {
		return info, fmt.Errorf("update checking is not configured")
	}

	manifest, err := FetchManifest(ManifestURL)
	if err != nil {
		return info, err
	}

	u.manifest = manifest

	currentVersion, err := ParseSemVer(Version)
	if err != nil {
		return info, fmt.Errorf("failed to parse current version: %w", err)
	}

	remoteVersion, err := ParseSemVer(manifest.Version)
	if err != nil {
		return info, fmt.Errorf("failed to parse remote version: %w", err)
	}

	platform := GetPlatform()
	downloadURL, err := manifest.GetDownloadURL(platform)
	if err != nil {
		return info, err
	}

	info.Version = manifest.Version
	info.ReleaseNotes = manifest.ReleaseNotes
	info.DownloadURL = downloadURL
	info.Available = remoteVersion.IsNewerThan(currentVersion)

	return info, nil
}

// DownloadAndInstall downloads and installs the update.
func (u *Updater) DownloadAndInstall() error {
	u.mu.Lock()
	if u.manifest == nil {
		u.mu.Unlock()
		return fmt.Errorf("no update available, call CheckForUpdate first")
	}
	manifest := u.manifest
	u.mu.Unlock()

	platform := GetPlatform()

	downloadURL, err := manifest.GetDownloadURL(platform)
	if err != nil {
		u.emitError(err.Error())
		return err
	}

	checksum, _ := manifest.GetChecksum(platform) // Checksum is optional

	// Create download context that can be cancelled
	ctx, cancel := context.WithCancel(u.ctx)
	u.mu.Lock()
	u.cancelDownload = cancel
	u.mu.Unlock()

	defer func() {
		u.mu.Lock()
		u.cancelDownload = nil
		u.mu.Unlock()
	}()

	// Emit downloading status
	u.emitStatus("downloading")

	// Download to temp directory
	tempDir := os.TempDir()
	execPath, err := GetExecutablePath()
	if err != nil {
		u.emitError(err.Error())
		return err
	}
	filename := filepath.Base(execPath)
	downloadPath := filepath.Join(tempDir, filename+".new")

	// Clean up any previous download
	os.Remove(downloadPath)

	err = DownloadFile(ctx, downloadURL, downloadPath, checksum, func(downloaded, total int64, percent float64) {
		u.emitProgress(downloaded, total, percent)
	})

	if err != nil {
		os.Remove(downloadPath)
		u.emitError(err.Error())
		return err
	}

	u.mu.Lock()
	u.downloadPath = downloadPath
	u.mu.Unlock()

	// Emit installing status
	u.emitStatus("installing")

	// Install the update
	err = InstallBinary(downloadPath)
	if err != nil {
		u.emitError(err.Error())
		return err
	}

	// Emit ready status
	u.emitStatus("ready")

	return nil
}

// CancelDownload cancels an ongoing download.
func (u *Updater) CancelDownload() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.cancelDownload != nil {
		u.cancelDownload()
	}
}

// RestartApp restarts the application.
func (u *Updater) RestartApp() error {
	execPath, err := GetExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Start a new instance of the application
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start new instance: %w", err)
	}

	// Quit the current application
	runtime.Quit(u.ctx)

	return nil
}

// emitProgress emits a progress event to the frontend.
func (u *Updater) emitProgress(downloaded, total int64, percent float64) {
	if u.ctx != nil {
		runtime.EventsEmit(u.ctx, "updater:progress", ProgressEvent{
			Downloaded: downloaded,
			Total:      total,
			Percent:    percent,
		})
	}
}

// emitStatus emits a status change event to the frontend.
func (u *Updater) emitStatus(status string) {
	if u.ctx != nil {
		runtime.EventsEmit(u.ctx, "updater:status", status)
	}
}

// emitError emits an error event to the frontend.
func (u *Updater) emitError(message string) {
	if u.ctx != nil {
		runtime.EventsEmit(u.ctx, "updater:error", message)
	}
}
