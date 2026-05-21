package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

func (h *handler) GetDocumentsByToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	rawToken := r.PathValue("token")

	logger.InfoContext(ctx, "Handler: Processing get documents by token request",
		"operation", "get_documents_by_token",
		"method", r.Method,
		"path", r.URL.Path)

	caseID, err := h.svc.ValidateToken(ctx, rawToken)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get documents by token")
		return
	}

	// Document repo is nil while document packages are disabled.
	if h.documentRepo == nil {
		httpx.RespondWithJSON(w, []document.DocumentSummary{}, http.StatusOK)
		return
	}

	filter := document.DocumentFilter{CaseID: &caseID}
	pagination := document.Pagination{Page: 1, PageSize: 100}

	docs, _, err := h.documentRepo.List(ctx, filter, pagination)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get documents by token")
		return
	}

	if docs == nil {
		docs = []document.DocumentSummary{}
	}

	httpx.RespondWithJSON(w, docs, http.StatusOK)
}
