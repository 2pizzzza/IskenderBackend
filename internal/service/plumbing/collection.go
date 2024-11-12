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

func (pr *Plumping) RemoveCollection(ctx context.Context, token string, req *models.RemoveCollectionRequest) error {
	const op = "service.RemoveCollection"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.DeleteCollection(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			log.Error("Collection not found", sl.Err(err))
			return storage.ErrCollectionNotFound
		}

		log.Error("Failed to remove collection", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (pr *Plumping) UpdateCollection(ctx context.Context, token string, req *models.UpdateCollectionRequest) error {
	const op = "service.UpdateCollection"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.UpdateCollection(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrCollectionNotFound) {
			log.Error("Collection not found", sl.Err(err))
			return storage.ErrCollectionNotFound
		}
		log.Error("Failed to update collection", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (pr *Plumping) GetCollectionRec(ctx context.Context, language string) ([]*models.CollectionResponse, error) {
	const op = "service.GetCollectionRec"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collection, err := pr.plumpingRepository.GetRandomCollectionsWithPopularity(ctx, language)
	if err != nil {
		log.Error("Failed to get collection", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return collection, nil
}

func (pr *Plumping) GetCollectionByStadart(ctx context.Context, code string) ([]*models.CollectionResponse, error) {
	const op = "service.GetCollectionByCategoryId"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collections, err := pr.plumpingRepository.GetCollectionsByIsProducerSLanguageCode(ctx, code)
	if err != nil {
		log.Error("Failed to get collections", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return collections, nil
}

func (pr *Plumping) GetCollectionByPainted(ctx context.Context, code string) ([]*models.CollectionResponse, error) {
	const op = "service.GetCollectionByCategoryId"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	collections, err := pr.plumpingRepository.GetCollectionsByIsProducerPLanguageCode(ctx, code)
	if err != nil {
		log.Error("Failed to get collections", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return collections, nil
}
