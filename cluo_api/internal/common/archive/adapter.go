package archive

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/document"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

// Adapter connects the archive builder to the application's real data sources.
type Adapter struct {
	docRepo    ports.DocumentRepository
	rapportSvc ports.RapportService
	mediaRepo  ports.MediaRepository
	storage    ports.StorageService
	crypto     encx.CryptoService
}

// NewAdapter creates an Adapter wired to the application's repositories and services.
func NewAdapter(
	docRepo ports.DocumentRepository,
	rapportSvc ports.RapportService,
	mediaRepo ports.MediaRepository,
	storage ports.StorageService,
	crypto encx.CryptoService,
) *Adapter {
	return &Adapter{
		docRepo:    docRepo,
		rapportSvc: rapportSvc,
		mediaRepo:  mediaRepo,
		storage:    storage,
		crypto:     crypto,
	}
}

func (a *Adapter) ListDocumentsByCase(ctx context.Context, caseID uuid.UUID) ([]document.DocumentSummary, error) {
	if a.docRepo == nil {
		return nil, nil
	}
	filter := document.DocumentFilter{CaseID: &caseID}
	pagination := document.Pagination{Page: 1, PageSize: 100}
	docs, _, err := a.docRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("list documents: %w", err)
	}
	return docs, nil
}

func (a *Adapter) GetDocument(ctx context.Context, id string, docType document.DocumentType) (document.Documentable, error) {
	if a.docRepo == nil {
		return nil, fmt.Errorf("document repository not available")
	}
	return a.docRepo.GetByID(ctx, id, docType)
}

func (a *Adapter) GetRapportContent(ctx context.Context, caseID uuid.UUID) ([]byte, error) {
	resp, err := a.rapportSvc.GetRapportByCaseID(ctx, caseID)
	if err != nil {
		return nil, nil // no rapport is not an error
	}
	return resp.Content, nil
}

func (a *Adapter) ListPublishedMediaByCase(ctx context.Context, caseID uuid.UUID) ([]*domainMedia.MediaFile, error) {
	encxList, _, err := a.mediaRepo.ListMediaByCaseID(ctx, caseID, nil, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("list media: %w", err)
	}

	var result []*domainMedia.MediaFile
	for _, enc := range encxList {
		if !enc.IsPublished {
			continue
		}
		media, err := domainMedia.DecryptMediaFileEncx(ctx, a.crypto, enc)
		if err != nil {
			continue // skip undecryptable media
		}
		result = append(result, media)
	}
	return result, nil
}

func (a *Adapter) DownloadMedia(ctx context.Context, url string) (io.ReadCloser, error) {
	return a.storage.DownloadFile(ctx, url)
}

func (a *Adapter) DecryptDocument(ctx context.Context, encDoc document.Documentable) (interface{}, error) {
	return document.DecryptDocumentable(ctx, a.crypto, encDoc)
}
