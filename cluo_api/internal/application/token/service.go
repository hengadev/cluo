package tokenService

import (
	"log/slog"

	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

type Service struct {
	repo            ports.TokenRepository
	caseRepo        ports.CaseRepository
	mediaRepo       ports.MediaRepository
	clientRepo      ports.ClientRepository
	crypto          encx.CryptoService
	emailService    ports.EmailService
	portalPublicURL string
	logger          *slog.Logger
}

func New(
	repo ports.TokenRepository,
	caseRepo ports.CaseRepository,
	mediaRepo ports.MediaRepository,
	clientRepo ports.ClientRepository,
	crypto encx.CryptoService,
	emailService ports.EmailService,
	portalPublicURL string,
	logger *slog.Logger,
) ports.TokenService {
	return &Service{
		repo:            repo,
		caseRepo:        caseRepo,
		mediaRepo:       mediaRepo,
		clientRepo:      clientRepo,
		crypto:          crypto,
		emailService:    emailService,
		portalPublicURL: portalPublicURL,
		logger:          logger,
	}
}
