package investigationService

import (
	// "context"

	// "github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/hengadev/encx"
)

type CaseService struct {
	repo            ports.CaseRepository
	clientRepo      ports.ClientRepository
	caseSubjectRepo ports.CaseSubjectRepository
	crypto          encx.CryptoService
}

func New(repo ports.CaseRepository, clientRepo ports.ClientRepository, caseSubjectRepo ports.CaseSubjectRepository, crypto encx.CryptoService) *CaseService {
	return &CaseService{
		repo:            repo,
		clientRepo:      clientRepo,
		caseSubjectRepo: caseSubjectRepo,
		crypto:          crypto,
	}
}

// func (s *CaseService) GetCaseByID(ctx context.Context, id string) (*domain.Case, error) {
// 	return s.repo.GetByID(ctx, id)
// }
//
// func (s *CaseService) UpdateCase(ctx context.Context, c *domain.Case) error {
// 	return s.repo.Update(ctx, c)
// }
//
// func (s *CaseService) DeleteCase(ctx context.Context, id string) error {
// 	return s.repo.Delete(ctx, id)
// }
//
// func (s *CaseService) MarkCaseAsReleased(ctx context.Context, id string) error {
// 	c, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}
//
// 	c.MarkAsReleased()
//
// 	return s.repo.Update(ctx, c)
// }
