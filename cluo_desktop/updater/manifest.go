package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Manifest represents the update manifest JSON structure.
type Manifest struct {
	Version      string            `json:"version"`
	ReleaseNotes string            `json:"release_notes"`
	Downloads    map[string]string `json:"downloads"`
	Checksums    map[string]string `json:"checksums"`
}

// FetchManifest downloads and parses the manifest from the given URL.
func FetchManifest(url string) (*Manifest, error) {
	if url == "" {
		return nil, fmt.Errorf("manifest URL is not configured")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manifest fetch failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest body: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// GetDownloadURL returns the download URL for the current platform.
func (m *Manifest) GetDownloadURL(platform string) (string, error) {
	url, ok := m.Downloads[platform]
	if !ok {
		return "", fmt.Errorf("no download available for platform: %s", platform)
	}
	return url, nil
}

// GetChecksum returns the checksum for the current platform.
func (m *Manifest) GetChecksum(platform string) (string, error) {
	checksum, ok := m.Checksums[platform]
	if !ok {
		return "", fmt.Errorf("no checksum available for platform: %s", platform)
	}
	return checksum, nil
}
