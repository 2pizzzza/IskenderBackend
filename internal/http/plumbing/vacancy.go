package plumbing

import (
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetAllVacancyActive retrieves all active vacancies filtered by language
// @Summary Get all active vacancies
// @Description Retrieves all active vacancies filtered by the specified language code
// @Tags vacancy
// @Accept json
// @Produce json
// @Param lang query string true "Language code for filtering vacancies (e.g., 'en', 'ru')"
// @Success 200 {array} models.VacancyResponse "Successfully retrieved vacancies"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid language parameter"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/vacancies/activ [get]
func (s *Server) GetAllVacancyActive(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing param lang"}, http.StatusBadRequest)
		return
	}

	reviews, err := s.service.GetAllActiveVacancyByLang(r.Context(), lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get brands"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, reviews, http.StatusOK)
}

// RemoveVacancy removes a vacancy by its ID
// @Summary Remove a vacancy
// @Description Deletes a vacancy based on the provided ID in the request body
// @Tags vacancy
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param vacancy body models.RemoveVacancyRequest true "Vacancy ID to remove"
// @Success 201 {object} models.Message "Successfully removed vacancy"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or vacancy not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/vacancy [delete]
func (s *Server) RemoveVacancy(w http.ResponseWriter, r *http.Request) {
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

	var req models.RemoveVacancyRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.RemoveVacancy(r.Context(), token, &req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrVacancyNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Vacancy not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove vacancy"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove vacancy"}, http.StatusCreated)
}

// UpdateVacancy updates a vacancy with new information
// @Summary Update a vacancy
// @Description Updates a vacancy's details based on the provided JSON body and authorization token
// @Tags vacancy
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param vacancy body models.VacancyUpdateRequest true "Updated vacancy details"
// @Success 201 {object} models.Message "Successfully updated vacancy"
// @Failure 400 {object} models.ErrorMessage "Invalid request body, vacancy not found, or vacancy translation not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/vacancy [put]
func (s *Server) UpdateVacancy(w http.ResponseWriter, r *http.Request) {
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

	var req models.VacancyUpdateRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.UpdateVacancy(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrVacancyNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Vacancy not found"}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 languages"}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrInvalidLanguageCode) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 languages kgz, ru, en"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to update vacancy"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful update vacancy"}, http.StatusCreated)
}

// GetAllVacancy retrieves all vacancies by language
// @Summary Get all active vacancies
// @Description Retrieves all  vacancies filtered by the specified language code
// @Tags vacancy
// @Accept json
// @Produce json
// @Param lang query string true "Language code for filtering vacancies (e.g., 'en', 'ru')"
// @Success 200 {array} models.VacancyResponse "Successfully retrieved vacancies"
// @Failure 400 {object} models.ErrorMessage "Missing or invalid language parameter"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/vacancies [get]
func (s *Server) GetAllVacancy(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing param lang"}, http.StatusBadRequest)
		return
	}

	reviews, err := s.service.GetAllVacancyByLang(r.Context(), lang)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get vacancy"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, reviews, http.StatusOK)
}

// GetVacancyById retrieves a vacancy by its ID
// @Summary Get a vacancy by ID
// @Description Retrieves the details of a specific vacancy using the provided vacancy ID parameter
// @Tags vacancy
// @Accept json
// @Produce json
// @Param vacancy_id query int true "The ID of the vacancy to retrieve"
// @Success 200 {object} models.VacancyResponses "Successfully retrieved vacancy"
// @Failure 400 {object} models.ErrorMessage "Invalid or missing vacancy ID"
// @Failure 404 {object} models.ErrorMessage "Vacancy not found"
// @Failure 500 {object} models.ErrorMessage "Failed to get vacancy by ID"
// @Router /api/vacancy [get]
func (s *Server) GetVacancyById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("vacancy_id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Collection Id"}, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid collection id"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.GetVacancyById(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrVacancyNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Vacancy not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get vacancy by id"}, http.StatusInternalServerError)
		return

	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

// CreateVacancy creates a new vacancy
// @Summary Create a new vacancy
// @Description Create a new vacancy with the provided information
// @Tags vacancy
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param body body models.VacancyResponses true "Vacancy details"
// @Success 201 {object} models.VacancyResponses "Successfully created vacancy"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or missing required language"
// @Failure 401 {object} models.ErrorMessage "Unauthorized"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Failed to create vacancy"
// @Router /api/vacancy [post]
func (s *Server) CreateVacancy(w http.ResponseWriter, r *http.Request) {
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

	var req models.VacancyResponses
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}
	res, err := s.service.CreateVacancy(r.Context(), token, &req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Required 3 language"}, http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrLanguageNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "language not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create vacancy"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

// SearchVacancy
// @Description Retrieves
// @Tags search
// @Accept json
// @Produce json
// @Param  q  query  string  false  "Search query"
// @Success 200 {object} models.VacancyResponse "Successfully retrieved vacancy"
// @Failure 400 {object} models.ErrorMessage "Invalid or missing vacancy ID"
// @Failure 404 {object} models.ErrorMessage "Vacancy not found"
// @Failure 500 {object} models.ErrorMessage "Failed to get vacancy by ID"
// @Router /api/searchVacancy [get]
func (s *Server) SearchVacancy(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")

	decodedQuery, err := url.QueryUnescape(searchQuery)
	if err != nil {
		fmt.Println("Error decoding query:", err)
	}

	res, err := s.service.SearchVacancy(r.Context(), decodedQuery)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed search vacancy"}, http.StatusInternalServerError)
		return

	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}
