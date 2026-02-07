package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
)

// RespondWithError responds with a structured JSON error.
// For 5xx errors, a generic message is returned to avoid leaking internal details.
func RespondWithError(w http.ResponseWriter, err error, status int) {
	message := err.Error()
	if status >= 500 {
		message = "Internal server error"
	}

	errorResponse := map[string]string{"error": message}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if encodeErr := encoder.Encode(errorResponse); encodeErr != nil {
		slog.Error("Failed to encode error response", "error", encodeErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, writeErr := w.Write(buf.Bytes()); writeErr != nil {
		slog.Error("Failed to write error response", "error", writeErr)
	}
}
