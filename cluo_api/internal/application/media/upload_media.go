package mediaService

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) UploadMedia(ctx context.Context, r *domain.UploadMediaRequest) (*domain.MediaResponse, error) {
	// Validate request
	if err := r.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Validate case exists
	caseID, _ := uuid.Parse(r.CaseID)
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(
			fmt.Errorf("case with ID %s not found", r.CaseID),
			"case",
		)
	}

	// Determine media type from MIME type
	mediaType := determineMediaTypeFromMimeType(r.MimeType)
	if mediaType == "" {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("unsupported mime type: %s", r.MimeType))
	}

	// Upload file to S3 storage
	fileURL, err := s.storage.UploadFile(ctx, r.File, r.FileName, r.MimeType, r.FileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// Create media record in database
	createRequest := &domain.CreateMediaRequest{
		CaseID:      r.CaseID,
		URL:         fileURL,
		Type:        mediaType,
		MimeType:    r.MimeType,
		FileName:    r.FileName,
		FileSize:    r.FileSize,
		Caption:     r.Caption,
		IsPublished: r.IsPublished,
	}

	// Use the private createMedia method to save to database
	response, err := s.createMedia(ctx, createRequest)
	if err != nil {
		// If database save fails, attempt to delete the uploaded file from S3
		deleteErr := s.storage.DeleteFile(ctx, fileURL)
		if deleteErr != nil {
			// Log the error but return the original error
			return nil, fmt.Errorf("failed to create media record (and failed to cleanup uploaded file): %w", err)
		}
		return nil, fmt.Errorf("failed to create media record: %w", err)
	}

	return response, nil
}

// determineMediaTypeFromMimeType determines the media type from MIME type
func determineMediaTypeFromMimeType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	} else if strings.HasPrefix(mimeType, "video/") {
		return "video"
	} else if strings.HasPrefix(mimeType, "audio/") {
		return "audio"
	}
	return ""
}
