package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
)

func (s *Server) GetAllCategoriesByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("lang")
	if code == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Fetching categories by code lang", slog.String("code", code))

	categories, err := s.service.GetCategoriesByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, storage.ErrLanguageNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Language code not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get categories by language code"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, categories, http.StatusOK)
}
