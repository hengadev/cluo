package document

import (
	"github.com/hengadev/encx"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Service handles business logic for document operations.
type Service struct {
	repo        ports.DocumentRepository
	versionRepo ports.DocumentVersionRepository
	caseRepo    ports.CaseRepository
	clientRepo  ports.ClientRepository
	crypto      encx.CryptoService
	// TODO: Add additional dependencies like:
	// - emailService ports.EmailService
	// - pdfService ports.PDFService
	// - notificationService ports.NotificationService
}

// New creates a new DocumentService instance.
func New(
	repo ports.DocumentRepository,
	versionRepo ports.DocumentVersionRepository,
	caseRepo ports.CaseRepository,
	clientRepo ports.ClientRepository,
	crypto encx.CryptoService,
) ports.DocumentService {
	return &Service{
		repo:        repo,
		versionRepo: versionRepo,
		caseRepo:    caseRepo,
		clientRepo:  clientRepo,
		crypto:      crypto,
	}
}