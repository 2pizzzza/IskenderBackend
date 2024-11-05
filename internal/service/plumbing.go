package service

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

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
