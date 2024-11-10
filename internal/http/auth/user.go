package auth

import (
	"encoding/json"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/utils"
	"net/http"
)

// Register registers a new user
// @Summary Register a new user
// @Description Creates a new user account with a username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.UserRequest true "User registration details"
// @Success 201 "User successfully registered"
// @Failure 400 {object} models.ErrorMessage "Invalid input"
// @Failure 500 {object} models.ErrorMessage "Could not register"
// @Router /register [post]
func (h *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid input"}, http.StatusBadRequest)
		return
	}

	if err := h.service.Register(r.Context(), req.Username, req.Password); err != nil {
		h.log.Error("Errors", sl.Err(err))
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Could not register"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login authenticates an existing user and returns a JWT token
// @Summary User login
// @Description Authenticates a user and provides a JWT token if credentials are valid
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.UserRequest true "User login details"
// @Success 200 {object} models.Token "Successful login with JWT token"
// @Failure 400 {object} models.ErrorMessage "Invalid input"
// @Failure 401 {object} models.ErrorMessage "Could not login"
// @Router /login [post]
func (h *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid input"}, http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Could not login"}, http.StatusUnauthorized)
		return
	}

	utils.WriteResponseBody(w, models.Token{Token: token}, http.StatusOK)
}
