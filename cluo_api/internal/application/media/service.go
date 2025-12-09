package mediaService

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

type Service struct {
	repo     ports.MediaRepository
	caseRepo ports.CaseRepository
	storage  ports.StorageService
	crypto   encx.CryptoService
}

func New(repo ports.MediaRepository, caseRepo ports.CaseRepository, storage ports.StorageService, crypto encx.CryptoService) ports.MediaService {
	return &Service{
		repo:     repo,
		caseRepo: caseRepo,
		storage:  storage,
		crypto:   crypto,
	}
}
