package plumbing

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

func (pr *Plumping) CreateReview(ctx context.Context, req *models.CreateReviewRequest) error {
	const op = "service.GetCategoryByCode"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	err := pr.plumpingRepository.CreateReview(ctx, req.Username, req.Rating, req.Text)
	if err != nil {
		log.Error("Failed to create review", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) GetAllReview(ctx context.Context) ([]*models.ReviewResponse, error) {
	const op = "service.GetCategoryByCode"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	reviews, err := pr.plumpingRepository.GetAllReviews(ctx)
	if err != nil {
		log.Error("Failed to get all reviews", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return reviews, nil
}
