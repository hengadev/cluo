package rapportService

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

type Service struct {
	repo     ports.RapportRepository
	caseRepo ports.CaseRepository
	crypto   encx.CryptoService
}

func New(repo ports.RapportRepository, caseRepo ports.CaseRepository, crypto encx.CryptoService) ports.RapportService {
	return &Service{repo: repo, caseRepo: caseRepo, crypto: crypto}
}
