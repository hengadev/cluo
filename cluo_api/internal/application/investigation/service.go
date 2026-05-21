package investigationService

import (
	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/hengadev/encx"
)

type CaseService struct {
	repo            ports.CaseRepository
	clientRepo      ports.ClientRepository
	caseSubjectRepo ports.CaseSubjectRepository
	rapportRepo     ports.RapportRepository
	tokenService    ports.TokenService
	crypto          encx.CryptoService
}

func New(repo ports.CaseRepository, clientRepo ports.ClientRepository, caseSubjectRepo ports.CaseSubjectRepository, rapportRepo ports.RapportRepository, tokenService ports.TokenService, crypto encx.CryptoService) *CaseService {
	return &CaseService{
		repo:            repo,
		clientRepo:      clientRepo,
		caseSubjectRepo: caseSubjectRepo,
		rapportRepo:     rapportRepo,
		tokenService:    tokenService,
		crypto:          crypto,
	}
}
