package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
)

// Starter initializes the service or performs a specific start-up operation
// @Summary Initialize the service or perform a start-up operation
// @Description Starts the service, potentially creating necessary data or performing required tasks.
//
//	If a category already exists, it returns a BadRequest error, otherwise it returns a success message.
//
// @Tags starter
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ErrorMessage "Successfully created"
// @Failure 400 {object} models.ErrorMessage "already exists"
// @Failure 500 {object} models.ErrorMessage "Internal server error"
// @Router /starter [post]
func (s *Server) Starter(w http.ResponseWriter, r *http.Request) {

	err := s.service.Starter(r.Context())
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Already exist"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Something to wrong"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.ErrorMessage{Message: "Successfully create"}, http.StatusOK)
}
