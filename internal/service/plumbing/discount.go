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
	if discounts == nil {
		discounts = []models.Discount{}
	}
	return discounts, nil
}

func (pr *Plumping) CreateDiscount(ctx context.Context, token string, discount models.DiscountCreate) (*models.DiscountCreate, error) {
	const op = "service.CreateDiscount"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return nil, storage.ErrToken
	}

	discout, err := pr.plumpingRepository.CreateDiscount(ctx, discount)
	if err != nil {
		if errors.Is(err, storage.ErrDiscountExists) {
			log.Error("Discount exist", sl.Err(err))
			return nil, storage.ErrDiscountExists
		}
		slog.Error("Failed to create discount", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return discout, nil
}

func (pr *Plumping) DeleteDiscount(ctx context.Context, token string, discount models.DiscountRequest) error {
	const op = "service.DeleteDiscount"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.DeleteDiscount(ctx, discount.Id)
	if err != nil {
		if errors.Is(err, storage.ErrDiscountNotFound) {
			log.Error("Failed to found discount", sl.Err(err))
			return storage.ErrDiscountNotFound
		}
		log.Error("Failed to remove discount", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil

}
