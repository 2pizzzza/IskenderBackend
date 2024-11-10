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

// GetCollectionsByCategoryId retrieves collections for a specified category ID and language code
// @Summary Retrieve collections by category ID and language code
// @Description Returns a list of collections in the specified language for a category
// @Tags collections
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.CollectionResponse "List of collections"
// @Failure 400 {object} models.ErrorMessage "Missing language code"
// @Failure 500 {object} models.ErrorMessage "Failed to get collections"
// @Router /collections [get]
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

// GetCollectionById retrieves a collection by its ID and language code
// @Summary Retrieve a collection by ID and language code
// @Description Returns details of a specific collection in the specified language
// @Tags collections
// @Accept  json
// @Produce  json
// @Param  collection_id  query  int  true  "Collection ID"
// @Param  lang  query  string  true  "Language code"
// @Success 200 {object} models.CollectionResponse "Collection details"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid parameters"
// @Failure 404 {object} models.ErrorMessage "Collection not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /collection [get]
func (s *Server) GetCollectionById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("collection_id")
	lang := r.URL.Query().Get("lang")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
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

	res, err := s.service.GetCollectionByID(r.Context(), id, lang)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found collection"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get collection"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}
