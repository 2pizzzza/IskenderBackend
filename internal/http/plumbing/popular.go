package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
)

func (s *Server) GetPopular(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("lang")
	if code == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing lang"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get popular by lang", slog.String("lang", code))

	res, err := s.service.GetPopular(r.Context(), code)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get Popular"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

func (s *Server) GetNew(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("lang")
	if code == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing lang"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get new by lang", slog.String("lang", code))

	res, err := s.service.GetNew(r.Context(), code)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get new"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}
