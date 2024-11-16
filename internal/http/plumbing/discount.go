package plumbing

import (
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
)

// GetAllDiscount godoc
// @Summary Get all discounts
// @Description Retrieve a list of all available discounts.
// @Tags Discounts
// @Produce json
// @Success 200 {array} models.Discount "List of discounts"
// @Failure 500 {object} models.ErrorMessage "Failed to get all discounts"
// @Router /discounts [get]
func (s Server) GetAllDiscount(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Get all discount")

	res, err := s.service.GetAllDiscounts(r.Context())
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get all discount"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusOK)
}
