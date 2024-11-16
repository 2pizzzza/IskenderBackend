package plumbing

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

func (pr *Plumping) GetAllDiscounts(ctx context.Context) ([]models.Discount, error) {

	const op = "service.GetAllDiscounts"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	discounts, err := pr.plumpingRepository.GetAllDiscount(ctx)
	if err != nil {
		log.Error("Error get all discount", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return discounts, nil
}
