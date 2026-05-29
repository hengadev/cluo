package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ProgressCallback is called during download with progress information.
type ProgressCallback func(downloaded, total int64, percent float64)

// DownloadFile downloads a file from the given URL to the destination path.
// It reports progress via the callback and verifies the checksum if provided.
func DownloadFile(ctx context.Context, url, destPath, expectedChecksum string, progressCallback ProgressCallback) error {
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create temp file for download
	tempPath := destPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	var succeeded bool
	defer func() {
		file.Close()
		if !succeeded {
			os.Remove(tempPath)
		}
	}()

	total := resp.ContentLength
	var downloaded int64

	// Create a hasher for checksum verification
	hasher := sha256.New()

	// Create a multi-writer to write to both file and hasher
	writer := io.MultiWriter(file, hasher)

	// Read in chunks and report progress
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := writer.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("failed to write: %w", writeErr)
			}
			downloaded += int64(n)

			if progressCallback != nil {
				var percent float64
				if total > 0 {
					percent = float64(downloaded) / float64(total) * 100
				}
				progressCallback(downloaded, total, percent)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}
	}

	// Verify checksum if provided
	if expectedChecksum != "" {
		actualChecksum := "sha256:" + hex.EncodeToString(hasher.Sum(nil))
		// Support both "sha256:..." and plain hex formats
		expected := expectedChecksum
		if !strings.HasPrefix(expected, "sha256:") {
			expected = "sha256:" + expected
		}
		if actualChecksum != expected {
			return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actualChecksum)
		}
	}

	// Close before rename — required on Windows where open files cannot be renamed.
	// The deferred close is a harmless no-op after this.
	file.Close()

	// Move temp file to final destination
	if err := os.Rename(tempPath, destPath); err != nil {
		return fmt.Errorf("failed to move download to destination: %w", err)
	}

	succeeded = true
	return nil
}
