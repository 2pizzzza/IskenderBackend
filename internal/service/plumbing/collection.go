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

func (pr *Plumping) GetCollectionByCategoryId(ctx context.Context, code string) ([]*models.CollectionResponse, error) {
	const op = "service.GetCollectionByCategoryId"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collections, err := pr.plumpingRepository.GetCollectionsByLanguageCode(ctx, code)
	if err != nil {
		log.Error("Failed to get collections", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return collections, nil
}

func (pr *Plumping) GetCollectionByID(ctx context.Context, collectionId int, code string) (*models.CollectionResponse, error) {
	const op = "service.GetCollectionByID"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collection, err := pr.plumpingRepository.GetCollectionByID(ctx, collectionId, code)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			log.Error("Failed to found collection", sl.Err(err))
			return nil, storage.ErrCollectionNotFound
		}
		log.Error("Failed to get collection", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return collection, nil
}
