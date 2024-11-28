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

func (pr *Plumping) GetCategoriesByCode(ctx context.Context, languageCode string) ([]*models.Category, error) {
	const op = "service.GetCategoryByCode"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	category, err := pr.plumpingRepository.GetCategoriesByLanguageCode(ctx, languageCode)
	if err != nil {
		if errors.Is(err, storage.ErrLanguageNotFound) {
			log.Error("Failed to found language code", sl.Err(err))
			return nil, storage.ErrLanguageNotFound
		}
		log.Error("Failed to get categories", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return category, nil
}

func (pr *Plumping) CreateCategory(ctx context.Context, token string, req models.CreateCategoryRequest) (*models.CreateCategoryResponse, error) {
	const op = "service.CreateCategory"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return nil, storage.ErrToken
	}

	category, err := pr.plumpingRepository.CreateCategory(ctx, req)

	if err != nil {
		if errors.Is(err, storage.ErrCategoryExists) {
			log.Error("Category Exist", sl.Err(err))
			return nil, storage.ErrCategoryExists
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			log.Error("Language 3 need", sl.Err(err))
			return nil, storage.ErrRequiredLanguage
		}

		log.Error("Failed to create category", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", sl.Err(err))
	}

	return category, nil
}

func (pr *Plumping) UpdateCategory(ctx context.Context, token string, categoryID int, req []models.UpdateCategoriesResponse) error {
	const op = "service.UpdateCategory"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.UpdateCategory(ctx, categoryID, req)
	log.Info("id", categoryID)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			log.Error("Category already exist", sl.Err(err))
			return storage.ErrCategoryNotFound
		}
		log.Error("Failed to update category", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (pr *Plumping) RemoveCategory(ctx context.Context, token string, req *models.RemoveCategoryRequest) error {
	const op = "service.RemoveCategory"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.DeleteCategory(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			log.Error("Category already exist", sl.Err(err))
			return storage.ErrCategoryNotFound
		}
		log.Error("Failed to update category", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (pr Plumping) GetCategoryById(ctx context.Context, id int) (*models.GetCategoryRequest, error) {
	const op = "service.GetBrandById"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	brand, err := pr.plumpingRepository.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			log.Error("Brand not found", sl.Err(err))
			return nil, storage.ErrCategoryNotFound

		}
		log.Error("Failed to get category", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return brand, nil
}
