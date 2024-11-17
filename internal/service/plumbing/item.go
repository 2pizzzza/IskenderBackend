package plumbing

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"log/slog"
)

func (pr *Plumping) GetItemsByCategoryId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error) {
	const op = "service.GetItemByCategoryId"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetItemsByCategoryID(ctx, id, code)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			log.Error("Failed to found category", sl.Err(err))
			return nil, storage.ErrCategoryNotFound
		}
		log.Error("Failed to found collection", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return items, nil
}

func (pr *Plumping) GetItemById(ctx context.Context, id int, code string) (*models.ItemResponse, error) {
	const op = "service.GetItemById"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetItemByID(ctx, id, code)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			log.Error("Failed to found item", sl.Err(err))
			return nil, storage.ErrItemNotFound
		}
		log.Error("Failed to get item", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return items, nil
}

func (pr *Plumping) GetItemsByCollectionId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error) {
	const op = "service.GetItemById"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetItemsByCollectionID(ctx, id, code)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			log.Error("Failed to found item", sl.Err(err))
			return nil, storage.ErrCollectionNotFound
		}
		log.Error("Failed to get item", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return items, nil
}

func (pr *Plumping) GetItemsRec(ctx context.Context, id int, code string) ([]*models.ItemResponse, error) {
	const op = "service.GetItemsRec"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetRandomItemsWithPopularity(ctx, code, id)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			log.Error("Failed to found item", sl.Err(err))
			return nil, storage.ErrCollectionNotFound
		}
		log.Error("Failed to get item", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return items, nil
}

func (pr *Plumping) CreateItem(ctx context.Context, req models.CreateItem) (*models.CreateItemResponse, error) {
	const op = "service.CreateItem"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	item, err := pr.plumpingRepository.CreateItem(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrRequiredLanguage) {
			log.Error("Required 3 languages", sl.Err(err))
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		if errors.Is(err, storage.ErrInvalidLanguageCode) {
			log.Error("Required 3 languages kgz, ru, en", sl.Err(err))
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		if errors.Is(err, storage.ErrItemExists) {
			log.Error("Item already exist", sl.Err(err))
			return nil, fmt.Errorf("%s, %w", op, err)
		}
		log.Error("Failed to create item", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return item, nil
}
