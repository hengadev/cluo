package tokenService

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
	domainToken "github.com/hengadev/cluo_api/internal/domain/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type stubCaseRepo struct {
	existsResult bool
	getCaseErr   error
}

func (r *stubCaseRepo) ExistsCase(_ context.Context, _ uuid.UUID) (bool, error) {
	return r.existsResult, nil
}
func (r *stubCaseRepo) GetCaseByID(_ context.Context, _ uuid.UUID) (*investigation.InvestigationEncx, error) {
	return nil, r.getCaseErr
}
func (r *stubCaseRepo) CreateCase(_ context.Context, _ *investigation.InvestigationEncx) error {
	return nil
}
func (r *stubCaseRepo) UpdateCase(_ context.Context, _ *investigation.InvestigationEncx) error {
	return nil
}
func (r *stubCaseRepo) DeleteCase(_ context.Context, _ uuid.UUID) error { return nil }
func (r *stubCaseRepo) List(_ context.Context, _ investigation.Filter, _ investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	return nil, 0, nil
}
func (r *stubCaseRepo) ListByClient(_ context.Context, _ uuid.UUID, _ investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	return nil, 0, nil
}

type stubTokenRepo struct{}

func (r *stubTokenRepo) CreateToken(_ context.Context, _ *domainToken.Token) error { return nil }
func (r *stubTokenRepo) GetTokenByHash(_ context.Context, _ string) (*domainToken.Token, error) {
	return nil, nil
}
func (r *stubTokenRepo) ListTokensByCaseID(_ context.Context, _ uuid.UUID) ([]*domainToken.Token, error) {
	return nil, nil
}
func (r *stubTokenRepo) RevokeToken(_ context.Context, _ uuid.UUID) error { return nil }

// signalEmailSvc closes a channel on the first Send call.
type signalEmailSvc struct {
	called chan struct{}
	err    error
}

func newSignalEmailSvc(err error) *signalEmailSvc {
	return &signalEmailSvc{called: make(chan struct{}, 1), err: err}
}

func (s *signalEmailSvc) Send(_ context.Context, _, _, _ string) error {
	select {
	case s.called <- struct{}{}:
	default:
	}
	return s.err
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestCreateToken_GoroutineErrorDoesNotPropagate verifies that an error inside
// the fire-and-forget goroutine (here: resolveClientEmail failing) does not
// affect the value returned by CreateToken.
func TestCreateToken_GoroutineErrorDoesNotPropagate(t *testing.T) {
	caseRepo := &stubCaseRepo{
		existsResult: true,
		getCaseErr:   errors.New("simulated db error"),
	}
	emailSvc := newSignalEmailSvc(errors.New("send failed"))

	svc := New(&stubTokenRepo{}, caseRepo, nil, nil, nil, emailSvc, "", slog.Default())

	resp, err := svc.CreateToken(context.Background(), uuid.New())
	require.NoError(t, err, "goroutine error must not propagate to caller")
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.RawToken)

	// Goroutine should not reach emailService.Send when resolveClientEmail fails.
	select {
	case <-emailSvc.called:
		t.Fatal("expected emailService.Send not to be called when email resolution fails")
	case <-time.After(50 * time.Millisecond):
	}
}

// TestCreateToken_SuccessReturnedImmediately verifies that CreateToken returns
// before the background goroutine completes, and returns the expected token shape.
func TestCreateToken_SuccessReturnedImmediately(t *testing.T) {
	caseRepo := &stubCaseRepo{existsResult: true, getCaseErr: errors.New("no case")}

	svc := New(&stubTokenRepo{}, caseRepo, nil, nil, nil, newSignalEmailSvc(nil), "https://portal.example.com", slog.Default())

	start := time.Now()
	resp, err := svc.CreateToken(context.Background(), uuid.New())
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	// The synchronous path should complete well under 1 second (no network calls).
	assert.Less(t, elapsed, time.Second)
}
