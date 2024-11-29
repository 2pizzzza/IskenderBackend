package plumbing

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	token2 "github.com/2pizzzza/plumbing/internal/lib/jwt"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"log/slog"
)

func (pr *Plumping) CreateBrand(ctx context.Context, token string, req *models.BrandRequest) (*models.BrandResponse, error) {
	const op = "service.CreateBrand"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return nil, storage.ErrToken
	}

	brand, err := pr.plumpingRepository.CreateBrand(ctx, req.Name, req.Url)
	if err != nil {
		if errors.Is(err, storage.ErrBrandExists) {
			log.Error("Brand already exist", sl.Err(err))
			return nil, storage.ErrBrandExists

		}
		log.Error("Failed to create brand", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return brand, nil
}

func (pr *Plumping) GetAllBrand(ctx context.Context) ([]*models.BrandResponse, error) {
	const op = "service.GetAllBrand"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	brands, err := pr.plumpingRepository.GetAllBrand(ctx)
	if err != nil {
		log.Error("Failed to get all brands", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	if brands == nil {
		brands = []*models.BrandResponse{}
	}
	return brands, nil
}

func (pr *Plumping) UpdateBrand(ctx context.Context, token string, id int, name, url string) (*models.BrandResponse, error) {
	const op = "service.UpdateBrand"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return nil, storage.ErrToken
	}

	brand, err := pr.plumpingRepository.UpdateBrand(ctx, id, name, url)
	if err != nil {
		if errors.Is(err, storage.ErrBrandNotFound) {
			log.Error("Brand not found", sl.Err(err))
			return nil, storage.ErrBrandNotFound

		}
		log.Error("Failed to update brand", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return brand, nil
}

func (pr *Plumping) RemoveBrand(ctx context.Context, token string, req *models.RemoveBrandRequest) error {
	const op = "service.RemoveBrand"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.RemoveBrand(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrBrandNotFound) {
			log.Error("Brand not found", sl.Err(err))
			return storage.ErrBrandNotFound

		}
		log.Error("Failed to remove brand", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr Plumping) GetBrandById(ctx context.Context, id int) (*models.BrandResponse, error) {
	const op = "service.GetBrandById"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	brand, err := pr.plumpingRepository.GetBrandByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrBrandNotFound) {
			log.Error("Brand not found", sl.Err(err))
			return nil, storage.ErrBrandNotFound

		}
		log.Error("Failed to get brand", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return brand, nil
}
