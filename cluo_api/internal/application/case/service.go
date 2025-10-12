package application

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

type CaseService struct {
	repo ports.CaseRepository
}

func New(repo ports.CaseRepository) *CaseService {
	return &CaseService{
		repo: repo,
	}
}

func (s *CaseService) CreateCase(ctx context.Context, c *domain.Case) error {
	return s.repo.Create(ctx, c)
}

func (s *CaseService) GetCaseByID(ctx context.Context, id string) (*domain.Case, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CaseService) UpdateCase(ctx context.Context, c *domain.Case) error {
	return s.repo.Update(ctx, c)
}

func (s *CaseService) DeleteCase(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CaseService) MarkCaseAsReleased(ctx context.Context, id string) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	c.MarkAsReleased()

	return s.repo.Update(ctx, c)
}
