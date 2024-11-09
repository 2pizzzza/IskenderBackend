package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
)

func (s *Server) GetCollectionsByCategoryId(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get all collection", slog.String("lang: ", lang))

	res, err := s.service.GetCollectionByCategoryId(r.Context(), lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get collections"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}
