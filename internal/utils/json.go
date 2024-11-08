package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func ReadRequestBody(r *http.Request, result interface{}) error {
	decoder := json.NewDecoder(r.Body)
	slog.Info("Body", result)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func WriteResponseBody(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
