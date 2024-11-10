package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
)

// GetPopular retrieves popular collections and items based on language
// @Summary Retrieve popular collections and items by language
// @Description Returns a list of popular collections and items for the specified language
// @Tags popular
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {object} models.PopularResponse "List of popular collections and items"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid language parameter"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /popular [get]
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

// GetNew retrieves new collections and items based on language
// @Summary Retrieve new collections and items by language
// @Description Returns a list of new collections and items for the specified language
// @Tags new
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {object} models.PopularResponse "List of new collections and items"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid language parameter"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /new [get]
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
