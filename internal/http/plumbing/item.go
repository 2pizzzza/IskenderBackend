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

// GetItemsByCategoryId retrieves items by category ID and language code
// @Summary Retrieve items by category ID and language code
// @Description Returns a list of items in the specified language for a category
// @Tags items
// @Accept  json
// @Produce  json
// @Param  category_id  query  int  true  "Category ID"
// @Param  lang         query  string  true  "Language code"
// @Success 200 {array} models.ItemResponse "List of items"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid parameters"
// @Failure 404 {object} models.ErrorMessage "Category not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/items [get]
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

// GetItemsById retrieves an item by its ID and language code
// @Summary Retrieve an item by ID and language code
// @Description Returns details of a specific item in the specified language
// @Tags items
// @Accept  json
// @Produce  json
// @Param  item_id  query  int  true  "Item ID"
// @Param  lang     query  string  true  "Language code"
// @Success 200 {object} models.ItemResponse "Item details"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid parameters"
// @Failure 404 {object} models.ErrorMessage "Item not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/item [get]
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

// GetItemsByCollectionId retrieves items by collection ID and language code
// @Summary Retrieve items by collection ID and language code
// @Description Returns a list of items in the specified language for a collection
// @Tags items
// @Accept  json
// @Produce  json
// @Param  collection_id  query  int  true  "Collection ID"
// @Param  lang           query  string  true  "Language code"
// @Success 200 {array} models.ItemResponse "List of items"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid parameters"
// @Failure 404 {object} models.ErrorMessage "Collection not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/items/collection [get]
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

// GetItemsRec retrieves items recommendation
// @Summary Retrieve items recommendation language code
// @Description Returns a list of items in the specified language for recommendation
// @Tags items
// @Accept  json
// @Produce  json
// @Param  item_id  query  int  true  "Item id"
// @Param  lang           query  string  true  "Language code"
// @Success 200 {array} models.ItemResponse "List of items"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid parameters"
// @Failure 404 {object} models.ErrorMessage "Collection not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/items/rec [get]
func (s *Server) GetItemsRec(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("item_id")
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

	s.log.Info("Get recommendation items", slog.String("lang: ", lang), slog.Int("id: ", id))

	res, err := s.service.GetItemsRec(r.Context(), id, lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get items"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}

// CreateItem godoc
// @Summary Create a new item
// @Description Create a new item with the specified details and upload photos.
// @Tags items
// @Accept multipart/form-data
// @Produce json
// @Param item formData string true "Item data in JSON format" example="{\"category_id\":1,\"collection_id\":2,\"size\":\"M\",\"price\":100.5,\"isProducer\":false,\"isPainted\":true,\"is_popular\":true,\"is_new\":false,\"items\":[{\"language_code\":\"en\",\"name\":\"Item Name\",\"description\":\"Item Description\"}]}"
// @Param photos formData file false "Photos of the item"
// @Param isMain_{filename} formData bool false "Indicates if the photo is the main one"
// @Param hashColor_{filename} formData string false "Color hash for the photo"
// @Success 201 {object} models.CreateItemResponse "Successfully created item"
// @Failure 400 {object} models.ErrorMessage "Invalid request data or item already exists"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to create item"
// @Router /api/items [post]
func (s *Server) CreateItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid form-data"}, http.StatusBadRequest)
		return
	}

	itemData := r.FormValue("item")
	if itemData == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing item data"}, http.StatusBadRequest)
		return
	}

	var req models.CreateItem
	if err := json.Unmarshal([]byte(itemData), &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid item data"}, http.StatusBadRequest)
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

		isMain := r.FormValue(fmt.Sprintf("isMain_%s", fileHeader.Filename)) == "true"
		hashColor := r.FormValue(fmt.Sprintf("hashColor_%s", fileHeader.Filename))

		photos = append(photos, models.PhotosResponse{
			URL:       "/media/images/" + filename,
			IsMain:    isMain,
			HashColor: hashColor,
		})
	}

	req.Photos = photos

	res, err := s.service.CreateItem(r.Context(), req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrItemExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Item with this name already exists"}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 languages"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create item"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

// UpdateItem godoc
// @Summary Update an existing item
// @Description Update a item with new details and photos.
// @Tags items
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param item_id query int true "Item ID"
// @Param item formData string true "Collection data in JSON format"
// @Param photos formData file true "Photos to upload"
// @Param isMain_{filename} formData bool false "Is this photo the main one?"
// @Param hashColor_{filename} formData string false "Hash color for the photo"
// @Success 201 {object} models.Message "Successfully updated the item"
// @Failure 400 {object} models.ErrorMessage "Invalid request (e.g., missing or invalid data)"
// @Failure 401 {object} models.ErrorMessage "Unauthorized or invalid token"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 404 {object} models.ErrorMessage "Collection not found"
// @Failure 500 {object} models.ErrorMessage "Failed to update collection"
// @Router /api/items [put]
func (s *Server) UpdateItem(w http.ResponseWriter, r *http.Request) {

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

	idStr := r.URL.Query().Get("item_id")
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

// RemoveItem godoc
// @Summary Remove an item from the collection
// @Description Remove a specific item from the collection based on the provided ID.
// @Tags items
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param item body models.ItemRequest true "Item ID to remove"
// @Success 201 {object} models.Message "Successfully removed the item"
// @Failure 400 {object} models.ErrorMessage "Invalid request or item not found"
// @Failure 401 {object} models.ErrorMessage "Unauthorized or invalid token"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to remove item"
// @Router /api/items [delete]
func (s *Server) RemoveItem(w http.ResponseWriter, r *http.Request) {
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

	var req models.ItemRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.RemoveItem(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrItemNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Item not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove collection"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove iten"}, http.StatusCreated)
}

// GetItemId retrieves a item by its ID and language code
// @Summary Retrieve a item by ID and language code
// @Description Returns details of a specific collection in the specified language
// @Tags items
// @Accept  json
// @Produce  json
// @Param  item_id  query  int  true  "Collection ID"
// @Success 200 {object} models.ItemResponseForAdmin "Collection details"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid parameters"
// @Failure 404 {object} models.ErrorMessage "Collection not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/getItemById [get]
func (s *Server) GetItemId(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("item_id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.GetItemID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found collection"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get collection"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

// GetAllItems retrieves all items
// @Summary Get all Items
// @Description Retrieves a list of all available items
// @Tags items
// @Produce json
// @Success 200 {array} models.ItemResponses "Successfully retrieved all items"
// @Failure 500 {object} models.ErrorMessage "Failed to get items"
// @Router /api/getAllItems [get]
func (s *Server) GetAllItems(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get all review")

	reviews, err := s.service.GetItems(r.Context())
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get brands"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, reviews, http.StatusOK)
}
