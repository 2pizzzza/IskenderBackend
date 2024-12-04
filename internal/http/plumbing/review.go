package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
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
