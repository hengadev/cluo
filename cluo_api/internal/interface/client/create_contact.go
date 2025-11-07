package clientHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (h *handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Extract client ID from URL path
	clientID := r.PathValue("id")
	if clientID == "" {
		httpx.RespondWithError(w, errs.NewInvalidValueErr("client ID is required"), http.StatusBadRequest)
		return
	}

	var payload client.CreateContactRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "create_contact",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, errs.NewInvalidValueErr(fmt.Sprintf("invalid request body: %v", err)), http.StatusBadRequest)
		return
	}

	payload.ClientID = clientID

	if err = h.svc.CreateContact(ctx, &payload); err != nil {

		// Log with specific error context based on error type
		var logLevel string
		var errorContext string
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			logLevel = "warn"
			errorContext = "client not found"
		case errors.Is(err, errs.ErrInvalidValue):
			logLevel = "warn"
			errorContext = "invalid request validation"
		case errors.Is(err, errs.ErrAlreadyExists):
			logLevel = "warn"
			errorContext = "conflict with existing contact"
		case errors.Is(err, errs.ErrDomainNotCreated):
			logLevel = "warn"
			errorContext = "contact not created"
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
			"operation", "create_contact",
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
		case errors.Is(err, errs.ErrAlreadyExists):
			statusCode = http.StatusConflict
		case errors.Is(err, errs.ErrDomainNotCreated):
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
			logger.InfoContext(ctx, "Handler: Create contact request result", logFields...)
		case "warn":
			logger.WarnContext(ctx, "Handler: Create contact request failed", logFields...)
		case "error":
			logger.ErrorContext(ctx, "Handler: Create contact request failed", logFields...)
		}
		httpx.RespondWithError(w, err, statusCode)
		return
	}

	// Log successful operation
	logger.InfoContext(ctx, "Handler: Create contact request completed successfully",
		"operation", "create_contact",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	// Respond with success message (no session cookie changes)
	httpx.RespondWithJSON(w, struct {
		Message string `json:"message"`
	}{
		Message: "Contact creation completed successfully",
	}, http.StatusOK)
}
