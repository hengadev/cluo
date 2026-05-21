package pieceService

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

// Service implements ports.PieceService.
type Service struct {
	repo     ports.PieceRepository
	caseRepo ports.CaseRepository
	storage  ports.StorageService
	crypto   encx.CryptoService
}

// New creates a new PieceService.
func New(repo ports.PieceRepository, caseRepo ports.CaseRepository, storage ports.StorageService, crypto encx.CryptoService) ports.PieceService {
	return &Service{
		repo:     repo,
		caseRepo: caseRepo,
		storage:  storage,
		crypto:   crypto,
	}
}
