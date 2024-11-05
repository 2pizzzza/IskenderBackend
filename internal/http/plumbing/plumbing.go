package plumbing

import (
	"encoding/json"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

func (s *Server) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	res, err := s.service.CreateItem(r.Context(), &req)
	if err != nil {
		s.log.Error("Failed to create item", sl.Err(err))
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	s.log.Info("Item created successfully", slog.String("name", res.Name))
}
