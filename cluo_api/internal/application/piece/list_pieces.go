package pieceService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (s *Service) ListPiecesByCaseID(ctx context.Context, caseID uuid.UUID, page, pageSize int) (*piece.ListPiecesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Validate case exists
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(
			fmt.Errorf("case with ID %s not found", caseID),
			"case",
		)
	}

	pagination := piece.Pagination{Page: page, PageSize: pageSize}

	encxList, total, err := s.repo.ListPiecesByCaseID(ctx, caseID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list pieces: %w", err)
	}

	responses := make([]*piece.PieceResponse, 0, len(encxList))
	for _, encx := range encxList {
		p, err := piece.DecryptPieceEncx(ctx, s.crypto, encx)
		if err != nil {
			// Skip pieces that cannot be decrypted
			continue
		}
		responses = append(responses, p.ToResponse())
	}

	return &piece.ListPiecesResponse{
		Pieces:     responses,
		Pagination: piece.NewPaginationInfo(page, pageSize, total),
	}, nil
}
