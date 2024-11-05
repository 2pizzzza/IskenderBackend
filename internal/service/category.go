package service

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

func (rp *Plumping) CreateCategory(ctx context.Context, req *schemas.CreateCategoryRequest) (res *schemas.CreateCategoryResponse, err error) {
	const op = "service.CreateItem"

	log := rp.log.With(
		slog.String("op", op),
	)

	item, err := rp.plumpingRepository.SaveCategory(ctx, req.Name)

	log.Debug("value: ", req.Name)
	if err != nil {
		log.Error("Failed Create Song", sl.Err(err))

		return &schemas.CreateCategoryResponse{}, err
	}
	return &schemas.CreateCategoryResponse{
		CategoryID: item.CategoryID,
		Name:       item.Name,
	}, nil
}
