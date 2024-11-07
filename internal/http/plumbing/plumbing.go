package plumbing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"io"
	"log/slog"
	"net/http"
	"os"
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

func (s *Server) CreateItemWithDetails(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	isProduced := r.FormValue("is_produced") == "true"
	colors := r.Form["colors"]

	files := r.MultipartForm.File["photos"]
	var photoPaths []string

	for _, file := range files {
		filePath := fmt.Sprintf("media/images/%s", file.Filename)

		dst, err := os.Create(filePath)
		if err != nil {
			s.log.Error("error", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		src, err := file.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusInternalServerError)
			return
		}
		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}

		photoPaths = append(photoPaths, filePath)
	}

	req := &schemas.CreateItemWithDetailsRequest{
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
		Price:       price,
		IsProduced:  isProduced,
		Colors:      colors,
		Photos:      photoPaths,
	}

	item, err := s.service.SaveItemWithDetails(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
