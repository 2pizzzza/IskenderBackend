package plumbing

import (
	"encoding/json"
	"errors"
	"fmt"
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

// GetCollectionsStandart retrieves collections standart
// @Summary Retrieve collections standart
// @Description Returns a list of collections standart
// @Tags collections
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.CollectionResponse "List of collections"
// @Failure 400 {object} models.ErrorMessage "Missing language code"
// @Failure 500 {object} models.ErrorMessage "Failed to get collections"
// @Router /collections [get]
func (s *Server) GetCollectionsStandart(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get all collection", slog.String("lang: ", lang))

	res, err := s.service.GetCollectionByStadart(r.Context(), lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get collections"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}

// GetCollectionsByPainted retrieves collections painted
// @Summary Retrieve collections painted
// @Description Returns a list of collections in the specified language for a category
// @Tags collections
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.CollectionResponse "List of collections"
// @Failure 400 {object} models.ErrorMessage "Missing language code"
// @Failure 500 {object} models.ErrorMessage "Failed to get collections"
// @Router /collections [get]
func (s *Server) GetCollectionsByPainted(w http.ResponseWriter, r *http.Request) {
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

// CreateCollection godoc
// @Summary Create a new collection
// @Description Create a new collection with the specified details and upload photos.
// @Tags collections
// @Accept multipart/form-data
// @Produce json
// @Param collection formData string true "Collection data in JSON format" example="{\"category_id\":1,\"collection_id\":2,\"size\":\"M\",\"price\":100.5,\"isProducer\":false,\"isPainted\":true,\"is_popular\":true,\"is_new\":false,\"items\":[{\"language_code\":\"en\",\"name\":\"Item Name\",\"description\":\"Item Description\"}]}"
// @Param photos formData file false "Photos of the item"
// @Param isMain_{filename} formData bool false "Indicates if the photo is the main one"
// @Param hashColor_{filename} formData string false "Color hash for the photo"
// @Success 201 {object} models.CreateCollectionResponse"Successfully created collection"
// @Failure 400 {object} models.ErrorMessage "Invalid request data or collection already exists"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to create collection"
// @Router /collection [post]
func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid form-data"}, http.StatusBadRequest)
		return
	}

	itemData := r.FormValue("collection")
	if itemData == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing collection data"}, http.StatusBadRequest)
		return
	}

	var req models.CreateCollectionRequest
	if err := json.Unmarshal([]byte(itemData), &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection data"}, http.StatusBadRequest)
		return
	}

	var photos []models.PhotosResponse
	files := r.MultipartForm.File["photos"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to open uploaded file"}, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename, err := saveImage(file, fileHeader.Filename)
		if err != nil {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to save image"}, http.StatusInternalServerError)
			return
		}

		isMain := r.FormValue(fmt.Sprintf("isMain_%s", fileHeader.Filename)) == "false"
		hashColor := r.FormValue(fmt.Sprintf("hashColor_%s", fileHeader.Filename))

		photos = append(photos, models.PhotosResponse{
			URL:       "/media/images/" + filename,
			IsMain:    isMain,
			HashColor: hashColor,
		})
	}

	req.Photos = photos

	res, err := s.service.CreateCollection(r.Context(), req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrCollectionExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Collection with this name already exists"}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 languages"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create collection"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

// UpdateCollection godoc
// @Summary Update an existing collection
// @Description Update a collection with new details and photos.
// @Tags collections
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param collection_id query int true "Collection ID"
// @Param collection formData string true "Collection data in JSON format"
// @Param photos formData file true "Photos to upload"
// @Param isMain_{filename} formData bool false "Is this photo the main one?"
// @Param hashColor_{filename} formData string false "Hash color for the photo"
// @Success 201 {object} models.Message "Successfully updated the collection"
// @Failure 400 {object} models.ErrorMessage "Invalid request (e.g., missing or invalid data)"
// @Failure 401 {object} models.ErrorMessage "Unauthorized or invalid token"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 404 {object} models.ErrorMessage "Collection not found"
// @Failure 500 {object} models.ErrorMessage "Failed to update collection"
// @Router /collections [put]
func (s *Server) UpdateCollection(w http.ResponseWriter, r *http.Request) {

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

	idStr := r.URL.Query().Get("collection_id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid form-data"}, http.StatusBadRequest)
		return
	}

	itemData := r.FormValue("collection")
	if itemData == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing collection data"}, http.StatusBadRequest)
		return
	}

	var req models.CreateCollectionRequest
	if err := json.Unmarshal([]byte(itemData), &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection data"}, http.StatusBadRequest)
		return
	}

	var photos []models.PhotosResponse
	files := r.MultipartForm.File["photos"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to open uploaded file"}, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename, err := saveImage(file, fileHeader.Filename)
		if err != nil {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to save image"}, http.StatusInternalServerError)
			return
		}

		isMain := r.FormValue(fmt.Sprintf("isMain_%s", fileHeader.Filename)) == "false"
		hashColor := r.FormValue(fmt.Sprintf("hashColor_%s", fileHeader.Filename))

		photos = append(photos, models.PhotosResponse{
			URL:       "/media/images/" + filename,
			IsMain:    isMain,
			HashColor: hashColor,
		})
	}

	req.Photos = photos

	err = s.service.UpdateCollection(r.Context(), token, id, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrCollectionNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Collection not fount"}, http.StatusNotFound)
			return
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 languages"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create collection"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful update collection"}, http.StatusCreated)
}

// GetAllCollection retrieves all collections
// @Summary Get all collections
// @Description Retrieves a list of all available collections
// @Tags collections
// @Produce json
// @Success 200 {array} models.CollectionResponses "Successfully retrieved all collections"
// @Failure 500 {object} models.ErrorMessage "Failed to get collections"
// @Router /getAllCollection [get]
func (s *Server) GetAllCollection(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get all review")

	reviews, err := s.service.GetCollection(r.Context())
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get brands"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, reviews, http.StatusOK)
}
