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
// @Router /items [get]
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
// @Router /item [get]
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
// @Router /items/collection [get]
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
// @Router /items/rec [get]
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
// @Router /items [post]
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
