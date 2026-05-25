package archive_test

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/archive"
	"github.com/hengadev/cluo_api/internal/domain/document"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mock dependencies
// ---------------------------------------------------------------------------

type mockDeps struct {
	docSummaries []document.DocumentSummary
	docByGetID   map[string]document.Documentable
	rapportJSON  []byte
	mediaList    []*domainMedia.MediaFile
	mediaBytes   map[string][]byte // URL → content
	downloadErr  map[string]error  // URL → error
}

func (m *mockDeps) ListDocumentsByCase(_ context.Context, _ uuid.UUID) ([]document.DocumentSummary, error) {
	return m.docSummaries, nil
}

func (m *mockDeps) GetDocument(_ context.Context, id string, _ document.DocumentType) (document.Documentable, error) {
	d, ok := m.docByGetID[id]
	if !ok {
		return nil, fmt.Errorf("not found: %s", id)
	}
	return d, nil
}

func (m *mockDeps) GetRapportContent(_ context.Context, _ uuid.UUID) ([]byte, error) {
	return m.rapportJSON, nil
}

func (m *mockDeps) ListPublishedMediaByCase(_ context.Context, _ uuid.UUID) ([]*domainMedia.MediaFile, error) {
	return m.mediaList, nil
}

func (m *mockDeps) DownloadMedia(_ context.Context, url string) (io.ReadCloser, error) {
	if err, ok := m.downloadErr[url]; ok {
		return nil, err
	}
	data, ok := m.mediaBytes[url]
	if !ok {
		return nil, fmt.Errorf("not found: %s", url)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (m *mockDeps) DecryptDocument(_ context.Context, encDoc document.Documentable) (interface{}, error) {
	return encDoc, nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestBuildFullArchive_ContainsRapportAndMedia(t *testing.T) {
	ctx := context.Background()
	caseID := uuid.New()

	rapportJSON := []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hello"}]}]}`)

	mediaURL1 := "https://storage.example.com/photo1.jpg"
	mediaURL2 := "https://storage.example.com/video1.mp4"

	deps := &mockDeps{
		docSummaries: []document.DocumentSummary{
			// A draft document — should be excluded
			{
				ID:          uuid.New(),
				CaseID:      caseID,
				Type:        document.DocumentTypeEstimate,
				Status:      document.DocumentStatusDraft,
				DocumentRef: "DEV-001",
			},
		},
		docByGetID:  map[string]document.Documentable{},
		rapportJSON: rapportJSON,
		mediaList: []*domainMedia.MediaFile{
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL1,
				Type: domainMedia.MediaTypeImage, MimeType: "image/jpeg",
				FileName: "photo1.jpg", FileSize: 100, IsPublished: true, CreatedAt: time.Now(),
			},
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL2,
				Type: domainMedia.MediaTypeVideo, MimeType: "video/mp4",
				FileName: "video1.mp4", FileSize: 200, IsPublished: true, CreatedAt: time.Now(),
			},
		},
		mediaBytes: map[string][]byte{
			mediaURL1: []byte("fake-jpeg-data"),
			mediaURL2: []byte("fake-mp4-data"),
		},
	}

	var buf bytes.Buffer
	err := archive.BuildFullArchive(ctx, deps, caseID, &buf)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)

	entryNames := make(map[string]bool)
	for _, f := range reader.File {
		entryNames[f.Name] = true
	}

	// Rapport PDF should be present
	assert.True(t, entryNames["rapport.pdf"], "archive should contain rapport.pdf")

	// Media files should be present
	assert.True(t, entryNames["media/photo1.jpg"], "archive should contain media/photo1.jpg")
	assert.True(t, entryNames["media/video1.mp4"], "archive should contain media/video1.mp4")

	// Draft documents should NOT be present
	for name := range entryNames {
		assert.NotContains(t, name, "draft", "no draft documents should be in the archive")
	}

	// rapport + 2 media = 3 entries (draft doc skipped)
	assert.Equal(t, 3, len(reader.File))
}

func TestBuildMediaArchive_ContainsOnlyMedia(t *testing.T) {
	ctx := context.Background()
	caseID := uuid.New()

	mediaURL1 := "https://storage.example.com/photo1.jpg"
	mediaURL2 := "https://storage.example.com/audio1.mp3"

	deps := &mockDeps{
		mediaList: []*domainMedia.MediaFile{
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL1,
				Type: domainMedia.MediaTypeImage, MimeType: "image/jpeg",
				FileName: "photo1.jpg", FileSize: 100, IsPublished: true, CreatedAt: time.Now(),
			},
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL2,
				Type: domainMedia.MediaTypeAudio, MimeType: "audio/mpeg",
				FileName: "audio1.mp3", FileSize: 200, IsPublished: true, CreatedAt: time.Now(),
			},
		},
		mediaBytes: map[string][]byte{
			mediaURL1: []byte("fake-jpeg-data"),
			mediaURL2: []byte("fake-mp3-data"),
		},
	}

	var buf bytes.Buffer
	err := archive.BuildMediaArchive(ctx, deps, caseID, &buf)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)

	entryNames := make(map[string]bool)
	for _, f := range reader.File {
		entryNames[f.Name] = true
	}

	assert.True(t, entryNames["media/photo1.jpg"], "archive should contain media/photo1.jpg")
	assert.True(t, entryNames["media/audio1.mp3"], "archive should contain media/audio1.mp3")
	assert.False(t, entryNames["rapport.pdf"], "media archive should NOT contain rapport.pdf")
	assert.Equal(t, 2, len(reader.File))

	// Verify file content is non-empty
	for _, f := range reader.File {
		rc, err := f.Open()
		require.NoError(t, err)
		content, err := io.ReadAll(rc)
		rc.Close()
		require.NoError(t, err)
		assert.NotEmpty(t, content, "file %s should have content", f.Name)
	}
}

func TestBuildMediaArchive_EmptyMediaList(t *testing.T) {
	ctx := context.Background()
	caseID := uuid.New()

	deps := &mockDeps{
		mediaList: []*domainMedia.MediaFile{},
	}

	var buf bytes.Buffer
	err := archive.BuildMediaArchive(ctx, deps, caseID, &buf)
	require.NoError(t, err)

	// Should produce a valid (empty) zip
	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)
	assert.Empty(t, reader.File, "empty media list should produce empty zip")
}

func TestBuildMediaArchive_S3FetchFailure_SkipsFile(t *testing.T) {
	// Decision: S3 fetch failure for one file causes that file to be skipped
	// and the archive continues with remaining files (best-effort approach).
	ctx := context.Background()
	caseID := uuid.New()

	mediaURL1 := "https://storage.example.com/ok-file.jpg"
	mediaURL2 := "https://storage.example.com/failing-file.jpg"

	deps := &mockDeps{
		mediaList: []*domainMedia.MediaFile{
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL1,
				Type: domainMedia.MediaTypeImage, MimeType: "image/jpeg",
				FileName: "ok-file.jpg", FileSize: 100, IsPublished: true, CreatedAt: time.Now(),
			},
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL2,
				Type: domainMedia.MediaTypeImage, MimeType: "image/jpeg",
				FileName: "failing-file.jpg", FileSize: 200, IsPublished: true, CreatedAt: time.Now(),
			},
		},
		mediaBytes: map[string][]byte{
			mediaURL1: []byte("good-data"),
		},
		downloadErr: map[string]error{
			mediaURL2: fmt.Errorf("S3 unavailable"),
		},
	}

	var buf bytes.Buffer
	err := archive.BuildMediaArchive(ctx, deps, caseID, &buf)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)

	// Only the working file should be present
	assert.Equal(t, 1, len(reader.File), "archive should contain only the successfully downloaded file")
	assert.Equal(t, "media/ok-file.jpg", reader.File[0].Name)
}

func TestBuildFullArchive_NoRapport_NoError(t *testing.T) {
	ctx := context.Background()
	caseID := uuid.New()

	mediaURL := "https://storage.example.com/photo.jpg"

	deps := &mockDeps{
		rapportJSON: nil, // no rapport
		mediaList: []*domainMedia.MediaFile{
			{
				ID: uuid.New(), CaseID: caseID, URL: mediaURL,
				Type: domainMedia.MediaTypeImage, MimeType: "image/jpeg",
				FileName: "photo.jpg", FileSize: 50, IsPublished: true, CreatedAt: time.Now(),
			},
		},
		mediaBytes: map[string][]byte{
			mediaURL: []byte("photo-data"),
		},
	}

	var buf bytes.Buffer
	err := archive.BuildFullArchive(ctx, deps, caseID, &buf)
	require.NoError(t, err)

	reader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	require.NoError(t, err)

	entryNames := make(map[string]bool)
	for _, f := range reader.File {
		entryNames[f.Name] = true
	}

	assert.False(t, entryNames["rapport.pdf"], "archive should NOT contain rapport.pdf when none exists")
	assert.True(t, entryNames["media/photo.jpg"], "archive should contain media files")
}
