package document

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/hengadev/cluo_api/internal/domain/document"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type sendTestDocRepo struct {
	doc domain.Documentable
}

func (r *sendTestDocRepo) GetByID(_ context.Context, _ string, _ domain.DocumentType) (domain.Documentable, error) {
	return r.doc, nil
}
func (r *sendTestDocRepo) Create(_ context.Context, _ domain.Documentable) error  { return nil }
func (r *sendTestDocRepo) Update(_ context.Context, _ domain.Documentable) error  { return nil }
func (r *sendTestDocRepo) Delete(_ context.Context, _ string, _ domain.DocumentType) error {
	return nil
}
func (r *sendTestDocRepo) List(_ context.Context, _ domain.DocumentFilter, _ domain.Pagination) ([]domain.DocumentSummary, int, error) {
	return nil, 0, nil
}
func (r *sendTestDocRepo) GetLinkedDocuments(_ context.Context, _ string, _ domain.DocumentType) ([]domain.Documentable, error) {
	return nil, nil
}
func (r *sendTestDocRepo) GetFirstByCaseAndType(_ context.Context, _ string, _ domain.DocumentType) (domain.Documentable, error) {
	return nil, nil
}
func (r *sendTestDocRepo) CreateEstimate(_ context.Context, _ *domain.EstimateEncx) error {
	return nil
}
func (r *sendTestDocRepo) GetEstimateByID(_ context.Context, _ string) (*domain.EstimateEncx, error) {
	return nil, nil
}
func (r *sendTestDocRepo) UpdateEstimate(_ context.Context, _ *domain.EstimateEncx) error {
	return nil
}
func (r *sendTestDocRepo) DeleteEstimate(_ context.Context, _ string) error { return nil }
func (r *sendTestDocRepo) ListEstimatesByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.EstimateEncx, int, error) {
	return nil, 0, nil
}
func (r *sendTestDocRepo) CreateMandate(_ context.Context, _ *domain.MandateEncx) error {
	return nil
}
func (r *sendTestDocRepo) GetMandateByID(_ context.Context, _ string) (*domain.MandateEncx, error) {
	return nil, nil
}
func (r *sendTestDocRepo) UpdateMandate(_ context.Context, _ *domain.MandateEncx) error {
	return nil
}
func (r *sendTestDocRepo) DeleteMandate(_ context.Context, _ string) error { return nil }
func (r *sendTestDocRepo) ListMandatesByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.MandateEncx, int, error) {
	return nil, 0, nil
}
func (r *sendTestDocRepo) CreateContract(_ context.Context, _ *domain.ContractEncx) error {
	return nil
}
func (r *sendTestDocRepo) GetContractByID(_ context.Context, _ string) (*domain.ContractEncx, error) {
	return nil, nil
}
func (r *sendTestDocRepo) UpdateContract(_ context.Context, _ *domain.ContractEncx) error {
	return nil
}
func (r *sendTestDocRepo) DeleteContract(_ context.Context, _ string) error { return nil }
func (r *sendTestDocRepo) ListContractsByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.ContractEncx, int, error) {
	return nil, 0, nil
}
func (r *sendTestDocRepo) CreateInvoice(_ context.Context, _ *domain.InvoiceEncx) error {
	return nil
}
func (r *sendTestDocRepo) GetInvoiceByID(_ context.Context, _ string) (*domain.InvoiceEncx, error) {
	return nil, nil
}
func (r *sendTestDocRepo) UpdateInvoice(_ context.Context, _ *domain.InvoiceEncx) error {
	return nil
}
func (r *sendTestDocRepo) DeleteInvoice(_ context.Context, _ string) error { return nil }
func (r *sendTestDocRepo) ListInvoicesByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.InvoiceEncx, int, error) {
	return nil, 0, nil
}
func (r *sendTestDocRepo) ListOverdueInvoices(_ context.Context, _ domain.Pagination) ([]*domain.InvoiceEncx, int, error) {
	return nil, 0, nil
}

type sendTestVersionRepo struct{}

func (r *sendTestVersionRepo) CreateVersion(_ context.Context, _ *domain.DocumentVersion) error {
	return nil
}
func (r *sendTestVersionRepo) GetDocumentHistory(_ context.Context, _ string, _ domain.DocumentType, _ domain.Pagination) ([]*domain.DocumentVersion, int, error) {
	return nil, 0, nil
}
func (r *sendTestVersionRepo) GetVersion(_ context.Context, _ string, _ domain.DocumentType, _ int) (*domain.DocumentVersion, error) {
	return nil, nil
}
func (r *sendTestVersionRepo) DeleteVersions(_ context.Context, _ string, _ domain.DocumentType) error {
	return nil
}

// signalEmailSvc notifies via a buffered channel when Send is called.
type sendTestEmailSvc struct {
	called chan struct{}
	err    error
}

func newSendTestEmailSvc(err error) *sendTestEmailSvc {
	return &sendTestEmailSvc{called: make(chan struct{}, 1), err: err}
}

func (s *sendTestEmailSvc) Send(_ context.Context, _, _, _ string) error {
	select {
	case s.called <- struct{}{}:
	default:
	}
	return s.err
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestSendDocument_EmailFailureDoesNotPropagate verifies that a Send failure
// inside the fire-and-forget goroutine does not affect the return value of
// SendDocument.
func TestSendDocument_EmailFailureDoesNotPropagate(t *testing.T) {
	caseID := uuid.New()
	clientID := uuid.New()
	estimate := domain.NewEstimate(caseID, clientID, "EST-TEST-001", nil)
	// Draft status is the default; Draft → Sent is a valid transition.

	repo := &sendTestDocRepo{doc: estimate}
	emailSvc := newSendTestEmailSvc(errors.New("SMTP connection refused"))

	svc := New(repo, &sendTestVersionRepo{}, nil, nil, nil, emailSvc, slog.Default())

	req := &domain.SendDocumentRequest{Recipients: []string{"client@example.com"}}
	err := svc.SendDocument(context.Background(), estimate.ID.String(), domain.DocumentTypeEstimate, req)
	require.NoError(t, err, "email send failure must not propagate to caller")

	// The goroutine must have been fired and attempted the send.
	select {
	case <-emailSvc.called:
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected emailService.Send to be called within timeout")
	}
}

// TestSendDocument_MultipleRecipients_EmailFailureDoesNotPropagate verifies
// that failures across multiple goroutines are all isolated.
func TestSendDocument_MultipleRecipients_EmailFailureDoesNotPropagate(t *testing.T) {
	caseID := uuid.New()
	clientID := uuid.New()
	estimate := domain.NewEstimate(caseID, clientID, "EST-TEST-002", nil)

	repo := &sendTestDocRepo{doc: estimate}
	called := make(chan struct{}, 3)
	emailSvc := &sendTestEmailSvc{called: called, err: errors.New("SMTP down")}

	svc := New(repo, &sendTestVersionRepo{}, nil, nil, nil, emailSvc, slog.Default())

	req := &domain.SendDocumentRequest{
		Recipients: []string{"a@example.com", "b@example.com", "c@example.com"},
	}
	err := svc.SendDocument(context.Background(), estimate.ID.String(), domain.DocumentTypeEstimate, req)
	require.NoError(t, err, "email failures must not propagate even with multiple recipients")

	deadline := time.After(200 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case <-called:
		case <-deadline:
			assert.FailNow(t, "not all goroutines completed within timeout")
		}
	}
}
