package whisper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MaxAudioSize is the maximum audio file size (100 MB).
	MaxAudioSize = 100 * 1024 * 1024
)

// SupportedAudioFormats is a list of supported audio formats.
var SupportedAudioFormats = map[string]bool{
	"audio/wav":        true,
	"audio/wave":       true,
	"audio/x-wav":      true,
	"audio/mpeg":       true, // MP3
	"audio/mp3":        true,
	"audio/mp4":        true, // M4A
	"audio/x-m4a":      true,
	"audio/ogg":        true,
	"audio/aac":        true,
	"audio/flac":       true,
	"audio/x-flac":     true,
	"audio/webm":       true,
	"audio/webm;codecs=opus": true,
}

// AudioProcessor handles audio file validation and processing.
type AudioProcessor struct {
	tempDir string
}

// NewAudioProcessor creates a new audio processor.
func NewAudioProcessor(tempDir string) *AudioProcessor {
	return &AudioProcessor{
		tempDir: tempDir,
	}
}

// ValidateAudioHeader validates an audio file from multipart form data.
func (p *AudioProcessor) ValidateAudioHeader(fh *multipart.FileHeader) error {
	// Check file size
	if fh.Size > MaxAudioSize {
		return fmt.Errorf("audio file too large: %d bytes (max %d bytes)", fh.Size, MaxAudioSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !p.isValidAudioExtension(ext) {
		return fmt.Errorf("unsupported audio format: %s", ext)
	}

	return nil
}

// ValidateAudioData validates audio data from a reader.
func (p *AudioProcessor) ValidateAudioData(reader multipart.File) error {
	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read audio header: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(buffer[:n])
	if !SupportedAudioFormats[contentType] {
		return fmt.Errorf("unsupported audio content type: %s", contentType)
	}

	return nil
}

// SaveTempAudio saves audio data to a temporary file.
func (p *AudioProcessor) SaveTempAudio(data []byte, filename string) (string, error) {
	// Create temp file with original extension
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".wav" // Default to wav
	}

	tempFile, err := os.CreateTemp(p.tempDir, "audio_*"+ext)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tempFile.Close()

	// Write audio data
	if _, err := tempFile.Write(data); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("write audio data: %w", err)
	}

	return tempFile.Name(), nil
}

// CleanupAudio deletes an audio file.
func (p *AudioProcessor) CleanupAudio(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete audio file: %w", err)
	}
	return nil
}

// isValidAudioExtension checks if the file extension is a valid audio format.
func (p *AudioProcessor) isValidAudioExtension(ext string) bool {
	validExtensions := map[string]bool{
		".wav":  true,
		".wave": true,
		".mp3":  true,
		".mpeg": true,
		".m4a":  true,
		".ogg":  true,
		".aac":  true,
		".flac": true,
		".webm": true,
		".opus": true,
	}

	return validExtensions[strings.ToLower(ext)]
}

// GetAudioFormat returns the audio format from a file extension.
func (p *AudioProcessor) GetAudioFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".wav", ".wave":
		return "wav"
	case ".mp3", ".mpeg":
		return "mp3"
	case ".m4a":
		return "m4a"
	case ".ogg":
		return "ogg"
	case ".aac":
		return "aac"
	case ".flac":
		return "flac"
	case ".webm", ".opus":
		return "webm"
	default:
		return "unknown"
	}
}
