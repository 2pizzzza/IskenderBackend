package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
)

func (db *DB) CreateReview(ctx context.Context, username string, rating int, text string) error {
	const op = "postgres.CreateReview"

	query := `
		INSERT INTO Review (username, rating, text)
		VALUES ($1, $2, $3)
	`
	_, err := db.Pool.Exec(ctx, query, username, rating, text)
	if err != nil {
		return fmt.Errorf("%s: failed to create review: %w", op, err)
	}

	return nil
}
func (db *DB) GetAllReviews(ctx context.Context) ([]*models.ReviewResponse, error) {
	const op = "postgres.GetAllReviews"

	query := `
		SELECT id, username, rating, text, created_at
		FROM Review
		ORDER BY created_at DESC
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query reviews: %w", op, err)
	}
	defer rows.Close()

	var reviews []*models.ReviewResponse
	for rows.Next() {
		var review models.ReviewResponse
		if err := rows.Scan(&review.ID, &review.Username, &review.Rating, &review.Text, &review.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: failed to scan review row: %w", op, err)
		}
		reviews = append(reviews, &review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return reviews, nil
}
