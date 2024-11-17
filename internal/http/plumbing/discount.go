package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
)

// GetAllDiscount godoc
// @Summary Get all discounts
// @Description Retrieve a list of all available discounts.
// @Tags Discounts
// @Produce json
// @Param  lang  query  string  true  "Language code"
// @Success 200 {array} models.Discount "List of discounts"
// @Failure 500 {object} models.ErrorMessage "Failed to get all discounts"
// @Router /discounts [get]
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
