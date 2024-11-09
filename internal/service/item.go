package service

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
