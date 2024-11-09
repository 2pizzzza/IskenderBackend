package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) GetItemsByCategoryId(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("category_id")
	lang := r.URL.Query().Get("lang")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Category Id"}, http.StatusBadRequest)
		return
	}
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid category id"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get all items", slog.String("lang: ", lang), slog.Int("id: ", id))

	res, err := s.service.GetItemsByCategoryId(r.Context(), id, lang)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found category"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get items"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

func (s *Server) GetItemsById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("item_id")
	lang := r.URL.Query().Get("lang")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Item Id"}, http.StatusBadRequest)
		return
	}
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid item id"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get all items", slog.String("lang: ", lang), slog.Int("id: ", id))

	res, err := s.service.GetItemById(r.Context(), id, lang)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found item"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get item"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}

func (s *Server) GetItemsByCollectionId(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("collection_id")
	lang := r.URL.Query().Get("lang")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing collection Id"}, http.StatusBadRequest)
		return
	}
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get all items", slog.String("lang: ", lang), slog.Int("id: ", id))

	res, err := s.service.GetItemsByCollectionId(r.Context(), id, lang)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found collection"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get items"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}
