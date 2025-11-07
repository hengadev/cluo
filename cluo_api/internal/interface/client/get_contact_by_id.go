package clientHandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (h *handler) GetContactByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract contact ID from URL path
	contactID := r.PathValue("id")
	if contactID == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("contact ID is required"), http.StatusBadRequest)
		return
	}

	// Build the request
	request := &client.GetContactByIDRequest{
		ContactID: contactID,
	}

	contact, err := h.svc.GetContactByID(ctx, request)
	if err != nil {
		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "contact not found"
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid request validation"
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			logLevel = "error"
			errorContext = "database connection failure"
		case errors.Is(err, errs.ErrResourceExhausted):
			logLevel = "error"
			errorContext = "database resource exhaustion"
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			logLevel = "warn"
			errorContext = "request cancelled"
		case errors.Is(err, context.DeadlineExceeded):
			logLevel = "warn"
			errorContext = "request timeout"
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			logLevel = "error"
			errorContext = "database transaction failure"
		case errors.Is(err, errs.ErrPermissionDenied):
			logLevel = "error"
			errorContext = "database permission denied"
		case errors.Is(err, errs.ErrInvalidInput):
			logLevel = "warn"
			errorContext = "invalid input data"
		case errors.Is(err, errs.ErrDatabase):
			logLevel = "error"
			errorContext = "general database error"
		case errors.Is(err, errs.ErrNotDecrypted), errors.Is(err, errs.ErrNotEncrypted):
			logLevel = "error"
			errorContext = "data encryption/decryption failure"
		default:
			logLevel = "error"
			errorContext = "unexpected error"
		}

		logFields := []any{
			"operation", "get_contact_by_id",
			"error_context", errorContext,
			"method", r.Method,
			"path", r.URL.Path,
			"error", err,
		}
		var statusCode int
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, errs.ErrInvalidValue):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrResourceExhausted):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, context.DeadlineExceeded):
			statusCode = http.StatusRequestTimeout
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			statusCode = http.StatusServiceUnavailable
		case errors.Is(err, errs.ErrPermissionDenied):
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrInvalidInput):
			statusCode = http.StatusBadRequest
		case errors.Is(err, errs.ErrDatabase):
			statusCode = http.StatusInternalServerError
		case errors.Is(err, errs.ErrNotDecrypted), errors.Is(err, errs.ErrNotEncrypted):
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}

		logFields = append(logFields, "status_code", statusCode)

		switch logLevel {
		case "info":
			logger.InfoContext(ctx, "Handler: Get contact by ID request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Get contact by ID request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Get contact by ID request failed", logFields...)
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Get contact by ID request completed successfully",
		"operation", "get_contact_by_id",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with contact data
	httpx.RespondWithJSON(w, struct {
		Message string                  `json:"message"`
		Contact *client.ContactResponse `json:"contact"`
	}{
		Message: "Contact retrieval by ID completed successfully",
		Contact: contact,
	}, http.StatusOK)
}

