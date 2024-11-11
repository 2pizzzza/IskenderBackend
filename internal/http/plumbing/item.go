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

	s.log.Info("Get rec items", slog.String("lang: ", lang), slog.Int("id: ", id))

	res, err := s.service.GetItemsRec(r.Context(), id, lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get items"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)

}
