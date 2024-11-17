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

func (pr *Plumping) GetAllDiscounts(ctx context.Context, languageCode string) ([]models.Discount, error) {

	const op = "service.GetAllDiscounts"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	discounts, err := pr.plumpingRepository.GetAllDiscount(ctx, languageCode)
	if err != nil {
		if errors.Is(err, storage.ErrLanguageNotFound) {
			log.Error("Language with this code not found", sl.Err(err))
			return nil, storage.ErrLanguageNotFound
		}
		log.Error("Error get all discount", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return discounts, nil
}
