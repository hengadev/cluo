package investigationService

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/hengadev/cluo_api/internal/domain/token"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- mocks ---

type mockCaseRepo struct{ mock.Mock }

func (m *mockCaseRepo) CreateCase(ctx context.Context, c *investigation.InvestigationEncx) error {
	return m.Called(ctx, c).Error(0)
}
func (m *mockCaseRepo) GetCaseByID(ctx context.Context, id uuid.UUID) (*investigation.InvestigationEncx, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*investigation.InvestigationEncx), args.Error(1)
}
func (m *mockCaseRepo) UpdateCase(ctx context.Context, c *investigation.InvestigationEncx) error {
	return m.Called(ctx, c).Error(0)
}
func (m *mockCaseRepo) DeleteCase(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockCaseRepo) List(ctx context.Context, f investigation.Filter, p investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	args := m.Called(ctx, f, p)
	return args.Get(0).([]*investigation.InvestigationEncx), args.Int(1), args.Error(2)
}
func (m *mockCaseRepo) ListByClient(ctx context.Context, clientID uuid.UUID, p investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	args := m.Called(ctx, clientID, p)
	return args.Get(0).([]*investigation.InvestigationEncx), args.Int(1), args.Error(2)
}
func (m *mockCaseRepo) ExistsCase(ctx context.Context, caseID uuid.UUID) (bool, error) {
	args := m.Called(ctx, caseID)
	return args.Bool(0), args.Error(1)
}

type mockTokenService struct{ mock.Mock }

func (m *mockTokenService) CreateToken(ctx context.Context, caseID uuid.UUID) (*token.CreateTokenResponse, error) {
	args := m.Called(ctx, caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*token.CreateTokenResponse), args.Error(1)
}
func (m *mockTokenService) ValidateToken(ctx context.Context, rawToken string) (uuid.UUID, error) {
	args := m.Called(ctx, rawToken)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
func (m *mockTokenService) ListTokensByCaseID(ctx context.Context, caseID uuid.UUID) ([]*token.TokenResponse, error) {
	args := m.Called(ctx, caseID)
	return args.Get(0).([]*token.TokenResponse), args.Error(1)
}
func (m *mockTokenService) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	return m.Called(ctx, tokenID).Error(0)
}
func (m *mockTokenService) GetCaseSummaryByToken(ctx context.Context, rawToken string) (*investigation.PortalCaseResponse, error) {
	args := m.Called(ctx, rawToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*investigation.PortalCaseResponse), args.Error(1)
}
func (m *mockTokenService) GetPublishedMediaByToken(ctx context.Context, rawToken string) ([]*domainMedia.MediaResponse, error) {
	args := m.Called(ctx, rawToken)
	return args.Get(0).([]*domainMedia.MediaResponse), args.Error(1)
}
func (m *mockTokenService) GetPublishedMediaByIDAndToken(ctx context.Context, rawToken string, mediaID uuid.UUID) (*domainMedia.MediaResponse, error) {
	args := m.Called(ctx, rawToken, mediaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainMedia.MediaResponse), args.Error(1)
}

// verify that our mocks satisfy the interfaces
var _ ports.CaseRepository = (*mockCaseRepo)(nil)
var _ ports.TokenService = (*mockTokenService)(nil)

// --- helpers ---

func testCtx(t *testing.T) context.Context {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	return context.WithValue(context.Background(), ctxutil.LoggerKey, logger)
}

func encryptedCase(t *testing.T, crypto encx.CryptoService, status investigation.Status) (*investigation.InvestigationEncx, uuid.UUID) {
	t.Helper()
	now := time.Now()
	caseID := uuid.New()
	c := &investigation.Investigation{
		ID:        caseID,
		Title:     "Test Case",
		ClientID:  uuid.New(),
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	encx, err := investigation.ProcessInvestigationEncx(testCtx(t), crypto, c)
	require.NoError(t, err)
	return encx, caseID
}

func fakeTokenResponse(caseID uuid.UUID) *token.CreateTokenResponse {
	return &token.CreateTokenResponse{
		ID:        uuid.New().String(),
		CaseID:    caseID.String(),
		RawToken:  "rawtoken123",
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

// --- tests ---

func TestRelease_InProgress_Returns409(t *testing.T) {
	crypto, err := encx.NewTestCrypto(t)
	require.NoError(t, err)

	caseEncx, caseID := encryptedCase(t, crypto, investigation.StatusInProgress)

	repo := &mockCaseRepo{}
	repo.On("GetCaseByID", mock.Anything, caseID).Return(caseEncx, nil)

	svc := &CaseService{repo: repo, tokenService: &mockTokenService{}, crypto: crypto}

	_, err = svc.Release(testCtx(t), caseID)

	assert.ErrorIs(t, err, errs.ErrConflict)
	repo.AssertNotCalled(t, "UpdateCase")
}

func TestRelease_Ready_TransitionsToReleased(t *testing.T) {
	crypto, err := encx.NewTestCrypto(t)
	require.NoError(t, err)

	caseEncx, caseID := encryptedCase(t, crypto, investigation.StatusReady)
	tokenResp := fakeTokenResponse(caseID)

	repo := &mockCaseRepo{}
	repo.On("GetCaseByID", mock.Anything, caseID).Return(caseEncx, nil)
	repo.On("UpdateCase", mock.Anything, mock.AnythingOfType("*investigation.InvestigationEncx")).Return(nil)

	tokenSvc := &mockTokenService{}
	tokenSvc.On("CreateToken", mock.Anything, caseID).Return(tokenResp, nil)

	svc := &CaseService{repo: repo, tokenService: tokenSvc, crypto: crypto}

	resp, err := svc.Release(testCtx(t), caseID)

	require.NoError(t, err)
	assert.Equal(t, caseID.String(), resp.CaseID)
	assert.Equal(t, tokenResp.ID, resp.TokenID)
	assert.Equal(t, tokenResp.RawToken, resp.RawToken)
	assert.Equal(t, "/token/"+tokenResp.RawToken, resp.PortalURL)
	repo.AssertCalled(t, "UpdateCase", mock.Anything, mock.Anything)
}

func TestRelease_AlreadyReleased_NewTokenNoStatusChange(t *testing.T) {
	crypto, err := encx.NewTestCrypto(t)
	require.NoError(t, err)

	caseEncx, caseID := encryptedCase(t, crypto, investigation.StatusReleased)
	tokenResp := fakeTokenResponse(caseID)

	repo := &mockCaseRepo{}
	repo.On("GetCaseByID", mock.Anything, caseID).Return(caseEncx, nil)

	tokenSvc := &mockTokenService{}
	tokenSvc.On("CreateToken", mock.Anything, caseID).Return(tokenResp, nil)

	svc := &CaseService{repo: repo, tokenService: tokenSvc, crypto: crypto}

	resp, err := svc.Release(testCtx(t), caseID)

	require.NoError(t, err)
	assert.Equal(t, tokenResp.RawToken, resp.RawToken)
	// status is already released — no DB update needed
	repo.AssertNotCalled(t, "UpdateCase")
}
