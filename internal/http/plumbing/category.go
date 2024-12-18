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

// GetAllCategoriesByCode fetches categories based on the provided language code
// @Summary Retrieve categories by language code
// @Description Returns a list of categories for a specified language code
// @Tags categories
// @Accept  json
// @Produce  json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.Category "List of categories"
// @Failure 400 {object} models.ErrorMessage "Missing language code"
// @Failure 404 {object} models.ErrorMessage "Language code not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/category [get]
func (s *Server) GetAllCategoriesByCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("lang")
	if code == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Fetching categories by code lang", slog.String("code", code))

	categories, err := s.service.GetCategoriesByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, storage.ErrLanguageNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Language code not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get categories by language code"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, categories, http.StatusOK)
}

// CreateCategory creates a new category with multi-language support
// @Summary Create a new category
// @Description Creates a new category. Requires a valid token and at least 3 language entries.
// @Tags categories
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer Token"
// @Param category body models.CreateCategoryRequest true "Category creation details"
// @Success 201 {object} models.CreateCategoryResponse "Successfully created category"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or category exists or less than required languages"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to create category"
// @Router /api/category [post]
func (s *Server) CreateCategory(w http.ResponseWriter, r *http.Request) {
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

	var req models.CreateCategoryRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.CreateCategory(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrCategoryExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Category with this name arledy exist"}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 language"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create category"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

// RemoveCategory removes an existing category
// @Summary Remove a category
// @Description Deletes a category. Requires a valid token.
// @Tags categories
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer Token"
// @Param category body models.RemoveCategoryRequest true "Category removal details"
// @Success 201 {object} models.Message "Successfully removed category"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or category not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to remove category"
// @Router /api/category [delete]
func (s *Server) RemoveCategory(w http.ResponseWriter, r *http.Request) {
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

	var req models.RemoveCategoryRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.RemoveCategory(r.Context(), token, &req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrCategoryNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Category not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove category"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove category"}, http.StatusCreated)
}

// UpdateCategory updates an existing category
// UpdateCategory updates an existing category
// @Summary Update a category
// @Description Updates a category with new details. Requires a valid token.
// @Tags categories
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer Token"
// @Param category_id query integer true "Category ID to retrieve"
// @Param category body []models.UpdateCategoriesResponse true "Category update details"
// @Success 201 {object} models.Message "Successfully updated category"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or category not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to update category"
// @Router /api/category [put]
func (s *Server) UpdateCategory(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.URL.Query().Get("category_id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
		return
	}
	var req []models.UpdateCategoriesResponse
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}
	err = s.service.UpdateCategory(r.Context(), token, id, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrCategoryNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Category not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to update category"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful update category"}, http.StatusCreated)
}

// GetCategoryById retrieves a category by its ID using query parameters
// @Summary Get a category by ID
// @Description Retrieves the details of a category using its ID from the query parameter
// @Tags categories
// @Accept json
// @Produce json
// @Param category_id query integer true "Category ID to retrieve"
// @Success 200 {object} models.GetCategoryRequest "Successfully retrieved category"
// @Failure 400 {object} models.ErrorMessage "Invalid or missing category ID"
// @Failure 404 {object} models.ErrorMessage "Category not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/category/by/id [get]
func (s *Server) GetCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("category_id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.GetCategoryById(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Category Not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get category by id"}, http.StatusInternalServerError)
		return

	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}
