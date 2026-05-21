package tokenService

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

type Service struct {
	repo      ports.TokenRepository
	caseRepo  ports.CaseRepository
	mediaRepo ports.MediaRepository
	crypto    encx.CryptoService
}

func New(repo ports.TokenRepository, caseRepo ports.CaseRepository, mediaRepo ports.MediaRepository, crypto encx.CryptoService) ports.TokenService {
	return &Service{
		repo:      repo,
		caseRepo:  caseRepo,
		mediaRepo: mediaRepo,
		crypto:    crypto,
	}
}
