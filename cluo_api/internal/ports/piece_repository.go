package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

// PieceRepository defines persistence operations for exhibit pieces.
type PieceRepository interface {
	CreatePiece(ctx context.Context, p *piece.PieceEncx) error
	GetPieceByID(ctx context.Context, id uuid.UUID) (*piece.PieceEncx, error)
	ListPiecesByCaseID(ctx context.Context, caseID uuid.UUID, pagination piece.Pagination) ([]*piece.PieceEncx, int, error)
	DeletePiece(ctx context.Context, id uuid.UUID) error
}
