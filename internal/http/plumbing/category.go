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

func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	s.log.Info("Creating category", slog.String("name", req.Name))

	res, err := s.service.CreateCategory(r.Context(), &req)
	if err != nil {
		s.log.Error("Failed to create category", sl.Err(err))
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	s.log.Info("Category created successfully", slog.String("name", res.Name))
}

func (s *Server) GetAllCategories(w http.ResponseWriter, r *http.Request) {

	res, err := s.service.GetAllCategory(r.Context())
	if err != nil {
		s.log.Error("Failed get all categories", sl.Err(err))
		http.Error(w, "Failed get all categories", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		s.log.Error("Failed to encode response", sl.Err(err))
		http.Error(w, "Failed encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing category Id", http.StatusBadRequest)
		return
	}

	categoryId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	s.log.Info("Fetching category by Item", slog.Int("category_id", categoryId))

	req := &schemas.CategoryByIdRequest{Id: categoryId}

	category, err := s.service.GetCategoryById(r.Context(), req)
	if err != nil {
		if errors.Is(err, schemas.ErrItemNotFound) {
			s.log.Error("Category not found", sl.Err(err))
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		s.log.Error("Failed to get category by ID", sl.Err(err))
		http.Error(w, "Failed to get category by ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		s.log.Error("Failed encode response", sl.Err(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
