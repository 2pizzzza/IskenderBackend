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

func (pr *Plumping) GetPopular(ctx context.Context, code string) (*models.PopularResponse, error) {
	const op = "service.GetPopular"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetPopularItems(ctx, code)
	if err != nil {
		log.Error("Failed to get populars items", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	collection, err := pr.plumpingRepository.GetPopularCollections(ctx, code)
	if err != nil {
		log.Error("Failed to get populars collections", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	res := &models.PopularResponse{
		Collections: collection,
		Items:       items,
	}

	return res, nil
}

func (pr *Plumping) GetNew(ctx context.Context, code string) (*models.PopularResponse, error) {
	const op = "service.GetNew"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	items, err := pr.plumpingRepository.GetNewItems(ctx, code)
	if err != nil {
		log.Error("Failed to get new items", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	collection, err := pr.plumpingRepository.GetNewCollections(ctx, code)
	if err != nil {
		log.Error("Failed to get new collections", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	res := &models.PopularResponse{
		Collections: collection,
		Items:       items,
	}

	return res, nil
}

func (pr *Plumping) Search(ctx context.Context, code string, isProducer *bool, isPainted *bool, searchQuery string) (*models.PopularResponse, error) {
	const op = "service.Search"

	log := pr.log.With(
		slog.String("op: ", op),
	)
	count := 0
	items, err := pr.plumpingRepository.SearchItems(ctx, code, isProducer, isPainted, searchQuery)
	if errors.Is(err, storage.ErrCollectionNotFound) {
		count++
	} else if err != nil {
		log.Error("Failed to get items", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	collection, err := pr.plumpingRepository.SearchCollections(ctx, code, isProducer, isPainted, searchQuery)
	if errors.Is(err, storage.ErrCollectionNotFound) {
		count++
		if count == 2 {
			return nil, storage.ErrCollectionNotFound
		}
	} else if err != nil {
		log.Error("Failed to get collections", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	if collection == nil {
		collection = []*models.CollectionResponse{}
	}
	if items == nil {
		items = []*models.ItemResponse{}
	}
	res := &models.PopularResponse{
		Collections: collection,
		Items:       items,
	}

	return res, nil
}
