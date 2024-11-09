package plumbing

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/2pizzzza/plumbing/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) CreateCatalog(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateCatalogRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}
	s.log.Info("Debug", slog.String("name", req.Name), slog.String("descri[tion", req.Description))
	defer r.Body.Close()

	s.log.Debug("Creating catalog", slog.String("name", req.Name))

	res, err := s.service.CreateCatalog(r.Context(), &req)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Catalog with this name exist"}, http.StatusBadRequest)
			return
		}
		s.log.Error("Failed to create catalog", sl.Err(err))
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create catalog"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

func (s *Server) CreateNewLocalizationForCatalog(w http.ResponseWriter, r *http.Request) {
	var req schemas.CatalogLocalizationRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	s.log.Info("Creating Localization for id: ", slog.Int("id: ", req.CatalogID))

	res, err := s.service.AddNewCatalogLocalization(r.Context(), &req)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogExists) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Catalog with this name arledy exist"}, http.StatusBadRequest)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to create localization"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, res, http.StatusCreated)
}

func (s *Server) GetAllCatalogsByLangCode(w http.ResponseWriter, r *http.Request) {
	var req schemas.CatalogsByLanguageRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Get all catalogs by code", slog.String("code", req.LanguageCode))

	catalogs, err := s.service.GetCatalogsByLangCode(r.Context(), &req)
	if err != nil {
		utils.WriteResponseBody(w, "Failed to get all catalogs", http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, catalogs, http.StatusOK)
}

func (s *Server) RemoveCatalog(w http.ResponseWriter, r *http.Request) {
	var req schemas.CatalogRemoveRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Remove catalog by id", slog.Int("id: ", req.ID))

	err := s.service.RemoveCatalog(r.Context(), &req)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Catalog with this id not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to remove catalog"}, http.StatusInternalServerError)
	}

	utils.WriteResponseBody(w, models.ErrorMessage{Message: "Successfully remove catalog"}, http.StatusOK)
}

func (s *Server) UpdateCatalog(w http.ResponseWriter, r *http.Request) {
	var req schemas.UpdateCatalogRequest
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	s.log.Info("Update catalog by id and code: ", slog.Int("id: ", req.ID), slog.Int("code: ", req.LanguageID))

	err := s.service.UpdateCatalog(r.Context(), &req)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found catalog"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to update catalog"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, models.ErrorMessage{Message: "Successfully update catalog"}, http.StatusOK)
}

func (s *Server) GetCatalogById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing category Id", http.StatusBadRequest)
		return
	}

	catalogId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	s.log.Info("Fetching catalog", slog.Int("catalog_id", catalogId))

	req, err := s.service.GetCatalogById(r.Context(), catalogId)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogNotFound) {
			utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to found catalog"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorMessage{Message: "Failed to get catalog"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, req, http.StatusOK)
}
