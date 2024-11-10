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

func (pr *Plumping) GetCategoriesByCode(ctx context.Context, languageCode string) ([]*models.Category, error) {
	const op = "service.GetCategoryByCode"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	category, err := pr.plumpingRepository.GetCategoriesByLanguageCode(ctx, languageCode)
	if err != nil {
		if errors.Is(err, storage.ErrLanguageNotFound) {
			log.Error("Failed to found language code", sl.Err(err))
			return nil, storage.ErrLanguageNotFound
		}
		log.Error("Failed to get categories", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return category, nil
}
