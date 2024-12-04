package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
	"strings"
)

// GetAllDiscount godoc
// @Summary Get all discounts
// @Description Retrieve a list of all available discounts.
// @Tags Discounts
// @Produce json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.Discount "List of discounts"
// @Failure 500 {object} models.ErrorMessage "Failed to get all discounts"
// @Router /api/discounts [get]
func (s Server) GetAllDiscount(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get all discount")
	code := r.URL.Query().Get("lang")
	if code == "" {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Missing Language"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.GetAllDiscounts(r.Context(), code)
	if err != nil {
		if errors.Is(err, storage.ErrLanguageNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Language with this code not fount"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get all discount"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}

// CreateDiscount godoc
// @Summary Create a new discount
// @Description Create a new discount with the specified details.
// @Tags Discounts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param discount body models.DiscountCreate true "Discount creation details"
// @Success 201 {object} models.DiscountCreate "Successfully created discount"
// @Failure 400 {object} models.ErrorMessage "Invalid request or discount already exists"
// @Failure 401 {object} models.ErrorMessage "Unauthorized or invalid token"
// @Failure 500 {object} models.ErrorMessage "Failed to create discount"
// @Router /api/discount [post]
func (s *Server) CreateDiscount(w http.ResponseWriter, r *http.Request) {
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
	var req models.DiscountCreate
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	res, err := s.service.CreateDiscount(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrDiscountExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Discount already exist"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to creaye discount"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

// RemoveDiscount godoc
// @Summary Remove a  discount
// @Description Remove a discount.
// @Tags Discounts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param discount body models.DiscountRequest true "Discount remove"
// @Success 201 {object} models.Message "Successfully remove discount"
// @Failure 400 {object} models.ErrorMessage "Invalid request or discount already exists"
// @Failure 401 {object} models.ErrorMessage "Unauthorized or invalid token"
// @Failure 500 {object} models.ErrorMessage "Failed to create discount"
// @Router /api/discount [delete]
func (s *Server) RemoveDiscount(w http.ResponseWriter, r *http.Request) {
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
	var req models.DiscountRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	err := s.service.DeleteDiscount(r.Context(), token, req)
	if err != nil {
		if errors.Is(err, storage.ErrDiscountNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Discount not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove discount"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.Message{Message: "Successful remove discount"}, http.StatusOK)
}
