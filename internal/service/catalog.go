package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"log/slog"
)

func (pr *Plumping) CreateCatalog(ctx context.Context, req *schemas.CreateCatalogRequest) (*schemas.CreateCatalogResponse, error) {
	const op = "service.CreateCatalog"

	log := pr.log.With(slog.String("op: ", op))

	catalog, err := pr.plumpingRepository.CreateCatalog(ctx, req.Name, req.Description, req.LanguageCode, req.Price, req.Color)

	if err != nil {
		if errors.Is(err, storage.ErrCatalogExists) {
			log.Error("Catalog with this name exist", sl.Err(err))
			return nil, storage.ErrCatalogExists
		}

		log.Error("Failed to create catalog", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return catalog, nil
}

func (pr *Plumping) AddNewCatalogLocalization(ctx context.Context, req *schemas.CatalogLocalizationRequest) (*schemas.CatalogLocalization, error) {
	const op = "service.AddNewCatalogLocalization"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	catalog, err := pr.plumpingRepository.InsertCatalogLocalization(ctx, req.CatalogID, req.LanguageCode, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogExists) {
			log.Error("Catalog with this name already exist", sl.Err(err))
			return nil, storage.ErrCatalogExists
		}

		log.Error("Failed to create localization for catalog", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return catalog, nil
}

func (pr *Plumping) GetCatalogsByLangCode(ctx context.Context, req *schemas.CatalogsByLanguageRequest) ([]*schemas.CatalogResponse, error) {
	const op = "service.GetCatalogByLangCode"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	catalogs, err := pr.plumpingRepository.GetCatalogsByLanguage(ctx, req.LanguageCode)

	if err != nil {
		log.Error("Failed to get all catalogs", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return catalogs, nil
}

func (pr *Plumping) RemoveCatalog(ctx context.Context, req *schemas.CatalogRemoveRequest) error {
	const op = "service.RemoveCatalog"

	log := pr.log.With(
		slog.String("op: ", op),
	)
	err := pr.plumpingRepository.DeleteCatalog(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogNotFound) {
			log.Error("Failed to found catalog by Id", sl.Err(err))
			return storage.ErrCatalogNotFound
		}

		log.Error("Failed to remove catalog", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) UpdateCatalog(ctx context.Context, req *schemas.UpdateCatalogRequest) error {
	const op = "service.UpdateCatalog"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	err := pr.plumpingRepository.UpdateCatalog(ctx, req.ID, req.LanguageCode, req.NewName, req.NewDescription, req.NewPrice)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogNotFound) {
			log.Error("Failed to found catalog", sl.Err(err))
			return storage.ErrCatalogNotFound
		}
		log.Error("Failed to update catalog", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) GetCatalogById(ctx context.Context, id int) (*schemas.CatalogDetailResponse, error) {
	const op = "service.GetCatalogByID"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	catalog, err := pr.plumpingRepository.GetCatalogByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrCatalogNotFound) {
			log.Error("Failed to found catalog", sl.Err(err))
			return nil, storage.ErrCatalogNotFound
		}
		log.Error("Failed to get catalog", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return catalog, nil
}
