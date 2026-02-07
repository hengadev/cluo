package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// RespondWithJSON responds with JSON data.
func RespondWithJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}
