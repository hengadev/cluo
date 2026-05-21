package ports

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

// PieceService defines business-logic operations for exhibit pieces.
type PieceService interface {
	UploadPiece(ctx context.Context, caseID uuid.UUID, file multipart.File, header *multipart.FileHeader, notes string) (*piece.PieceResponse, error)
	GetPieceByID(ctx context.Context, id uuid.UUID) (*piece.PieceResponse, error)
	ListPiecesByCaseID(ctx context.Context, caseID uuid.UUID, page, pageSize int) (*piece.ListPiecesResponse, error)
	DeletePiece(ctx context.Context, caseID uuid.UUID, pieceID uuid.UUID) error
}
