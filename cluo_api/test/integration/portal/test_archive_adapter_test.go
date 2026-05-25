package portal_test

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/hengadev/cluo_api/internal/common/archive"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

// fakeStorageService satisfies ports.StorageService for integration tests.
// DownloadFile returns placeholder bytes so the archive builder can produce
// real zip entries without a live S3 bucket. The other methods are unused.
type fakeStorageService struct{}

func (f *fakeStorageService) UploadFile(_ context.Context, _ io.Reader, _ string, _ string, _ int64) (string, error) {
	return "", fmt.Errorf("not available in tests")
}

func (f *fakeStorageService) DeleteFile(_ context.Context, _ string) error {
	return fmt.Errorf("not available in tests")
}

func (f *fakeStorageService) GetFileURL(_ context.Context, fileURL string) (string, error) {
	return fileURL, nil
}

func (f *fakeStorageService) DownloadFile(_ context.Context, url string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("fake-media-content-for: %s", url)))), nil
}

var _ ports.StorageService = (*fakeStorageService)(nil)

// newTestArchiveAdapter returns a real archive.Adapter wired to the test
// repositories but with a fake storage service instead of S3.
func newTestArchiveAdapter(
	docRepo ports.DocumentRepository,
	rapportSvc ports.RapportService,
	mediaRepo ports.MediaRepository,
	crypto encx.CryptoService,
) archive.Dependencies {
	return archive.NewAdapter(docRepo, rapportSvc, mediaRepo, &fakeStorageService{}, crypto)
}
