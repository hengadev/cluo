package pieceService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (s *Service) GetPieceByID(ctx context.Context, id uuid.UUID) (*piece.PieceResponse, error) {
	pieceEncx, err := s.repo.GetPieceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get piece: %w", err)
	}

	p, err := piece.DecryptPieceEncx(ctx, s.crypto, pieceEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("piece", err)
	}

	return p.ToResponse(), nil
}
