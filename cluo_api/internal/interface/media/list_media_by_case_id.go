package mediaHandler

import (
	"net/http"
	"strconv"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (h *handler) ListMediaByCaseID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Parse case ID from path
	caseID := r.PathValue("caseId")

	// Parse query params
	query := r.URL.Query()
	mediaType := query.Get("type")

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	request := &domain.ListMediaByCaseIDRequest{
		CaseID:   caseID,
		Page:     page,
		PageSize: pageSize,
	}

	if mediaType != "" {
		request.Type = &mediaType
	}

	logger.InfoContext(ctx, "Handler: Processing list media by case ID request",
		"operation", "list_media_by_case_id",
		"case_id", caseID,
		"page", page,
		"page_size", pageSize,
		"type", mediaType)

	response, err := h.svc.ListMediaByCaseID(ctx, request)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list media")
		return
	}

	logger.InfoContext(ctx, "Handler: List media by case ID completed successfully",
		"operation", "list_media_by_case_id",
		"case_id", caseID,
		"count", len(response.Media))

	httpx.RespondWithJSON(w, response, http.StatusOK)
}
