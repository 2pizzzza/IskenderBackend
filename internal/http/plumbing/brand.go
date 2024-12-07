package plumbing

import (
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	imageDir = "media/images"
)

// CreateBrand creates a new brand
// @Summary Create a new brand
// @Description Creates a new brand with a name and an optional image
// @Tags brand
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param name formData string true "Brand name"
// @Param photo formData file true "Brand photo"
// @Success 201 {object} models.BrandResponse "Successfully created brand"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or brand already exists"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/brand [post]
func (s *Server) CreateBrand(w http.ResponseWriter, r *http.Request) {

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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := saveImage(file, handler.Filename)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	req := &models.BrandRequest{
		Name: name,
		Url:  fmt.Sprintf("%s/%s", imageDir, filename),
	}

	res, err := s.service.CreateBrand(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrBrandExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Brand already exist"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Something to wri=ong"}, http.StatusInternalServerError)

		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

// GetAllBrands retrieves all brands
// @Summary Get all brands
// @Description Retrieves a list of all available brands
// @Tags brand
// @Produce json
// @Success 200 {array} models.BrandResponse "Successfully retrieved all brands"
// @Failure 500 {object} models.ErrorMessage "Failed to get brands"
// @Router /api/brands [get]
func (s *Server) GetAllBrands(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get all review")

	reviews, err := s.service.GetAllBrand(r.Context())
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get brands"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, reviews, http.StatusOK)
}

// RemoveBrand deletes a brand by ID
// @Summary Deletes a brand
// @Description Removes a brand by ID with authorization token required
// @Tags brand
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RemoveBrandRequest body models.RemoveBrandRequest true "Brand ID to delete"
// @Success 201 {object} models.Message "Successfully removed brand"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or brand not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/brand [delete]
func (s *Server) RemoveBrand(w http.ResponseWriter, r *http.Request) {
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

	var req models.RemoveBrandRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.RemoveBrand(r.Context(), token, &req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrBrandNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Brand not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to brand category"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove brand"}, http.StatusCreated)
}

// UpdateBrand updates a brand by ID
// @Summary Updates a brand's information
// @Description Updates the details of an existing brand by ID, including the option to upload a new image
// @Tags brand
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param id query integer true "Brand ID to update"
// @Param name formData string true "Updated brand name"
// @Param photo formData file false "Updated brand image"
// @Success 200 {object} models.BrandResponse "Successfully updated brand"
// @Failure 400 {object} models.ErrorMessage "Invalid brand ID or form data"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 404 {object} models.ErrorMessage "Brand not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/brand [put]
func (s *Server) UpdateBrand(w http.ResponseWriter, r *http.Request) {
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

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid brand ID"}, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Unable to parse form data"}, http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	file, handler, err := r.FormFile("photo")
	var url string

	if err == nil {
		defer file.Close()
		filename, err := saveImage(file, handler.Filename)
		if err != nil {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to save image"}, http.StatusInternalServerError)
			return
		}
		url = fmt.Sprintf("/%s/%s", imageDir, filename)
	}

	res, err := s.service.UpdateBrand(r.Context(), token, id, name, url)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if err == storage.ErrBrandNotFound {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Brand not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to update brand"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

// GetBrandById retrieves a brand by its ID using query parameters
// @Summary Get a brand by ID
// @Description Retrieves the details of a brand using its ID from the query parameter
// @Tags brand
// @Accept json
// @Produce json
// @Param brand_id query integer true "Brand ID to retrieve"
// @Success 200 {object} models.BrandResponse "Successfully retrieved brand"
// @Failure 400 {object} models.ErrorMessage "Invalid or missing brand ID"
// @Failure 404 {object} models.ErrorMessage "Brand not found"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/brand [get]
func (s *Server) GetBrandById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("brand_id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.GetBrandById(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrBrandNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Brand Not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get brand by id"}, http.StatusInternalServerError)
		return

	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

func saveImage(file io.Reader, filename string) (string, error) {
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = re.ReplaceAllString(filename, "_")

	filename = fmt.Sprintf("%s_%d%s", strings.TrimSuffix(filename, filepath.Ext(filename)), time.Now().Unix(), filepath.Ext(filename))
	path := filepath.Join(imageDir, filename)

	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			fmt.Printf("failed to close file: %v\n", cerr)
		}
	}()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return filename, nil
}
