package document

import (
	"log/slog"

	"github.com/hengadev/encx"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Service handles business logic for document operations.
type Service struct {
	repo         ports.DocumentRepository
	versionRepo  ports.DocumentVersionRepository
	caseRepo     ports.CaseRepository
	clientRepo   ports.ClientRepository
	crypto       encx.CryptoService
	emailService ports.EmailService
	logger       *slog.Logger
}

// New creates a new DocumentService instance.
func New(
	repo ports.DocumentRepository,
	versionRepo ports.DocumentVersionRepository,
	caseRepo ports.CaseRepository,
	clientRepo ports.ClientRepository,
	crypto encx.CryptoService,
	emailService ports.EmailService,
	logger *slog.Logger,
) ports.DocumentService {
	return &Service{
		repo:         repo,
		versionRepo:  versionRepo,
		caseRepo:     caseRepo,
		clientRepo:   clientRepo,
		crypto:       crypto,
		emailService: emailService,
		logger:       logger,
	}
}