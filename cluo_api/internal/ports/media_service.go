package ports

import (
	"context"

	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

// MediaService defines the business logic for media operations
type MediaService interface {
	// UploadMedia handles file upload to S3 and creates the media record in the database
	UploadMedia(ctx context.Context, r *domain.UploadMediaRequest) (*domain.MediaResponse, error)
	GetMediaByID(ctx context.Context, r *domain.GetMediaByIDRequest) (*domain.MediaResponse, error)
	UpdateMedia(ctx context.Context, r *domain.UpdateMediaRequest) (*domain.MediaResponse, error)
	DeleteMedia(ctx context.Context, r *domain.DeleteMediaRequest) error
	ListMediaByCaseID(ctx context.Context, r *domain.ListMediaByCaseIDRequest) (*domain.ListMediaResponse, error)
}
