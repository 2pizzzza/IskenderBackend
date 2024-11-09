package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
)

func (s *Server) GetAllLanguages(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get All Languages")

	languages, err := s.service.GetAllLanguages(r.Context())
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get all languages"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, languages, http.StatusOK)

}
