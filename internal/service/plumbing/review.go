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

	if reviews == nil {
		reviews = []*models.ReviewResponse{}
	}
	return reviews, nil
}

func (pr *Plumping) RemoveReview(ctx context.Context, token string, req models.RemoveReview) error {
	const op = "service.RemoveReview"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.DeleteReview(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrReviewNotFound) {
			log.Error("Review not found", sl.Err(err))
			return storage.ErrReviewNotFound

		}
		log.Error("Failed to remove review", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) SwitchIsShowReview(ctx context.Context, token string, req models.RemoveReview) error {
	const op = "service.SwitchIsShowReview"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.ToggleReviewVisibility(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrReviewNotFound) {
			log.Error("Review not found", sl.Err(err))
			return storage.ErrReviewNotFound

		}
		log.Error("Failed to remove review", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) GetAllReviewAdmin(ctx context.Context) ([]*models.ReviewResponseAdmin, error) {
	const op = "service.GetAllReviewAdmin"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	reviews, err := pr.plumpingRepository.GetAllReviewsAdmin(ctx)
	if err != nil {
		log.Error("Failed to get all reviews", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if reviews == nil {
		reviews = []*models.ReviewResponseAdmin{}
	}
	return reviews, nil
}
