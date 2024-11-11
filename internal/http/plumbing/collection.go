package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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

// RemoveCollection removes a collection by its ID
// @Summary Remove a collection
// @Description Removes a collection from the database by ID. Requires a valid token.
// @Tags collections
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer Token"
// @Param collection body models.RemoveCollectionRequest true "Collection removal details"
// @Success 201 {object} models.Message "Successfully removed collection"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or collection not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to remove collection"
// @Router /collection [delete]
func (s *Server) RemoveCollection(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Token required"}, http.StatusUnauthorized)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid token format"}, http.StatusUnauthorized)
		return
	}
	token := parts[1]

	var req models.RemoveCollectionRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.RemoveCollection(r.Context(), token, &req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrCollectionNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Collection not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove collection"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove collection"}, http.StatusCreated)
}

// GetCollectionsRec retrieves collections recommendation
// @Summary Retrieve collections recommendation by language code
// @Description Returns a list of collections recommendation
// @Tags collections
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.CollectionResponse "List of collections"
// @Failure 400 {object} models.ErrorMessage "Missing language code"
// @Failure 500 {object} models.ErrorMessage "Failed to get collections"
// @Router /collections/rec [get]
func (s *Server) GetCollectionsRec(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get rec collection", slog.String("lang: ", lang))

	res, err := s.service.GetCollectionRec(r.Context(), lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get collections"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}

//func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request) {
//	authHeader := r.Header.Get("Authorization")
//	if authHeader == "" {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Token required"}, http.StatusUnauthorized)
//		return
//	}
//
//	parts := strings.Split(authHeader, " ")
//	if len(parts) != 2 || parts[0] != "Bearer" {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid token format"}, http.StatusUnauthorized)
//		return
//	}
//	token := parts[1]
//
//	err := r.ParseMultipartForm(10 << 20)
//	if err != nil {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Unable to parse form data"}, http.StatusBadRequest)
//		return
//	}
//
//	files := r.MultipartForm.File["photos"]
//	if len(files) == 0 {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "No photos uploaded"}, http.StatusBadRequest)
//		return
//	}
//
//	uploadedFiles, err := s.service.UploadPhotos(r.Context(), files)
//	if err != nil {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to upload photos"}, http.StatusInternalServerError)
//		return
//	}
//
//	// Подготовка запроса
//	var req models.UpdateCollectionRequest
//	err = json.NewDecoder(r.Body).Decode(&req)
//	if err != nil {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid JSON in request body"}, http.StatusBadRequest)
//		return
//	}
//
//	err = s.service.UpdateCollection(r.Context(), token, &req)
//	if err != nil {
//		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to update collection"}, http.StatusInternalServerError)
//		return
//	}
//
//	utils.WriteResponseBody(w, models.UploadedPhotosResponse{Files: uploadedFiles}, http.StatusOK)
//}
