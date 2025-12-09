package ports

import (
	"context"

	"github.com/google/uuid"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

// MediaRepository stores metadata of uploaded media (file pointers stored in Storage).
type MediaRepository interface {
	CreateMedia(ctx context.Context, m *domain.MediaFileEncx) error
	GetMediaByID(ctx context.Context, id uuid.UUID) (*domain.MediaFileEncx, error)
	UpdateMedia(ctx context.Context, m *domain.MediaFileEncx) error
	DeleteMedia(ctx context.Context, id uuid.UUID) error
	ListMediaByCaseID(ctx context.Context, caseID uuid.UUID, mediaType *domain.MediaType, page, pageSize int) ([]*domain.MediaFileEncx, int, error)
}
