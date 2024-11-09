package service

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

func (pr *Plumping) GetCollectionByCategoryId(ctx context.Context, code string) ([]*models.CollectionResponse, error) {
	const op = "service.GetCollectionByCategoryId"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collection, err := pr.plumpingRepository.GetCollectionsByLanguageCode(ctx, code)
	if err != nil {
		log.Error("Failed to found collection", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return collection, nil
}
