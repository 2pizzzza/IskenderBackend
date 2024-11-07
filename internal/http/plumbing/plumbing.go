package plumbing

import (
	"encoding/json"
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strconv"
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

func (s *Server) GetItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	s.log.Info("Fetching item by ID", slog.Int("item_id", itemID))

	req := &schemas.GetItemByIdRequest{ItemID: itemID}

	item, err := s.service.GetItemById(r.Context(), req)
	if err != nil {

		if errors.Is(err, schemas.ErrItemNotFound) {
			s.log.Error("Item not found", sl.Err(err))
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		s.log.Error("Failed to get item by ID", sl.Err(err))
		http.Error(w, "Failed to get item by ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		s.log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	s.log.Info("Item fetched successfully", slog.Int("item_id", item.ItemID))
}
