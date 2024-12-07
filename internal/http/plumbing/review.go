package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
	"strings"
)

// GetAllReviews retrieves all reviews
// @Summary Get all reviews
// @Description Fetches all reviews from the database
// @Tags reviews
// @Produce  json
// @Success 200 {array} models.ReviewResponse "List of reviews"
// @Failure 500 {object} models.ErrorMessage "Failed to get reviews"
// @Router /api/reviews [get]
func (s *Server) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get all review")

	reviews, err := s.service.GetAllReview(r.Context())
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get reviews"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, reviews, http.StatusOK)
}

// CreateReview creates a new review
// @Summary Create a new review
// @Description Creates a review with the provided username, rating, and text
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.CreateReviewRequest true "Review data"
// @Success 201 {object} models.Message "Successfully created review"
// @Failure 400 {object} models.ErrorMessage "Invalid request body"
// @Failure 500 {object} models.ErrorMessage "Failed to create review"
// @Router /api/reviews [post]
func (s *Server) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req models.CreateReviewRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}
	s.log.Info("Create review")

	err := s.service.CreateReview(r.Context(), &req)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create review"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful create review"}, http.StatusCreated)
}

// RemoveReview deletes a review by ID
// @Summary Deletes a review
// @Description Removes a review by ID with authorization token required
// @Tags reviews
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RemoveReview body models.RemoveReview true "review ID to delete"
// @Success 201 {object} models.Message "Successfully removed review"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or review not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/reviews [delete]
func (s *Server) RemoveReview(w http.ResponseWriter, r *http.Request) {
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

	var req models.RemoveReview
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.RemoveReview(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrReviewNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Review not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove review"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove review"}, http.StatusCreated)
}

// SwitchIsShowReview
// @Description Switch isShow a review by ID with authorization token required
// @Tags reviews
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <token>"
// @Param RemoveReview body models.RemoveReview true "review ID to delete"
// @Success 201 {object} models.Message "Successfully switch review"
// @Failure 400 {object} models.ErrorMessage "Invalid request body or review not found"
// @Failure 401 {object} models.ErrorMessage "Token required or invalid token format"
// @Failure 403 {object} models.ErrorMessage "Permissions denied"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /api/switchIsShowReview [post]
func (s *Server) SwitchIsShowReview(w http.ResponseWriter, r *http.Request) {
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

	var req models.RemoveReview
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.SwitchIsShowReview(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrToken) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Permissions denied"}, http.StatusForbidden)
			return
		}
		if errors.Is(err, storage.ErrReviewNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Review not found"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to switch isShow"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful switch isShow review"}, http.StatusCreated)
}
