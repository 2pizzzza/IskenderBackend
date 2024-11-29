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

	if items == nil {
		items = []*models.ItemResponse{}
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
	if items == nil {
		items = []*models.ItemResponse{}
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

	if items == nil {
		items = []*models.ItemResponse{}
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

func (pr *Plumping) UpdateItem(ctx context.Context, token string, itemId int, req models.CreateItem) error {
	const op = "service.UpdateItem"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.UpdateItem(ctx, itemId, req)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			log.Error("Item not found", sl.Err(err))
			return storage.ErrCollectionNotFound
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			log.Error("Required 3 languages", sl.Err(err))
			return storage.ErrRequiredLanguage
		}
		if errors.Is(err, storage.ErrInvalidLanguageCode) {
			log.Error("Required 3 languages kgz, ru, en", sl.Err(err))
			return storage.ErrInvalidLanguageCode
		}
		log.Error("Failed to update item", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (pr *Plumping) RemoveItem(ctx context.Context, token string, req models.ItemRequest) error {
	const op = "service.RemoveItem"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.RemoveItem(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			log.Error("item not found", sl.Err(err))
			return storage.ErrItemNotFound
		}

		log.Error("Failed to remove collection", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (pr *Plumping) GetItemID(ctx context.Context, itemId int) (*models.ItemResponseForAdmin, error) {
	const op = "service.GetItemID"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collection, err := pr.plumpingRepository.GetItem(ctx, itemId)
	if err != nil {
		if errors.Is(err, storage.ErrItemNotFound) {
			log.Error("Failed to found collection", sl.Err(err))
			return nil, storage.ErrItemNotFound
		}
		log.Error("Failed to get collection", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return collection, nil
}

func (pr *Plumping) GetItems(ctx context.Context) ([]*models.ItemResponses, error) {
	const op = "service.GetItems"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetAllItems(ctx)
	if err != nil {
		log.Error("Failed to get all items", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if items == nil {
		items = []*models.ItemResponses{}
	}
	return items, nil
}
