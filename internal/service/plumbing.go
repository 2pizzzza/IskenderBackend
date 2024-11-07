package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

const imageDir = "media/image"

func (rp *Plumping) CreateItem(ctx context.Context, req *schemas.CreateItemRequest) (res *schemas.CreateItemResponse, err error) {
	const op = "service.CreateItem"

	log := rp.log.With(
		slog.String("op", op),
	)

	item, err := rp.plumpingRepository.SaveItem(ctx, req.Name, req.Description, req.CategoryID, req.Price)

	log.Debug("value: ", req.Name, req.Description, req.CategoryID, req.Price)
	if err != nil {
		log.Error("Failed Create Song", sl.Err(err))

		return &schemas.CreateItemResponse{}, err
	}
	return &schemas.CreateItemResponse{
		ItemID:      item.ItemID,
		Name:        item.Name,
		Description: item.Description,
		CategoryID:  item.CategoryID,
		Price:       item.Price,
		IsProduced:  item.IsProduced,
		Colors:      item.Colors,
		Photos:      item.Photos,
	}, nil
}

func (rp *Plumping) GetItemById(ctx context.Context, req *schemas.GetItemByIdRequest) (res *models.Item, err error) {
	const op = "service.GetItemById"

	log := rp.log.With(slog.String(
		"op: ", op),
	)

	item, err := rp.plumpingRepository.GetItemById(ctx, req.ItemID)

	if err != nil {
		if errors.Is(err, schemas.ErrItemNotFound) {
			return nil, schemas.ErrItemNotFound
		}
		log.Error("Err get item by id: ", sl.Err(err))

		return nil, fmt.Errorf("%s, %w", op, err)
	}

	res = &models.Item{
		ItemID:      item.ItemID,
		Name:        item.Name,
		Description: item.Description,
		CategoryID:  item.CategoryID,
		Price:       item.Price,
		IsProduced:  item.IsProduced,
		Colors:      item.Colors,
		Photos:      item.Photos,
	}

	return res, nil
}

func (rp *Plumping) SaveItemWithDetails(ctx context.Context, req *schemas.CreateItemWithDetailsRequest) (schemas.CreateItemResponse, error) {
	const op = "service.SaveItemWithDetails"

	log := rp.log.With(slog.String("op", op))

	item, err := rp.plumpingRepository.SaveItemWithDetails(
		ctx,
		req.Name,
		req.Description,
		req.CategoryID,
		req.Price,
		req.Colors,
		req.Photos,
	)
	if err != nil {
		log.Error("Failed to save item with details", sl.Err(err))
		return schemas.CreateItemResponse{}, err
	}

	return schemas.CreateItemResponse{
		ItemID:      item.ItemID,
		Name:        item.Name,
		Description: item.Description,
		CategoryID:  item.CategoryID,
		Price:       item.Price,
		IsProduced:  item.IsProduced,
		Colors:      item.Colors,
		Photos:      item.Photos,
	}, nil
}
