package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
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
		WHERE isShow=TRUE
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

func (db *DB) DeleteReview(ctx context.Context, id int) error {
	const op = "postgres.DeleteReview"

	existsQuery := `
    SELECT 1 
    FROM Review 
    WHERE id = $1`

	var exists int
	err := db.Pool.QueryRow(ctx, existsQuery, id).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return storage.ErrReviewNotFound
		}
		return fmt.Errorf("%s: failed to check if review exists: %w", op, err)
	}

	// Удаляем отзыв
	deleteQuery := `
    DELETE FROM Review 
    WHERE id = $1`

	_, err = db.Pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete review with ID %d: %w", op, id, err)
	}

	return nil
}

func (db *DB) ToggleReviewVisibility(ctx context.Context, id int) error {
	const op = "postgres.ToggleReviewVisibility"

	existsQuery := `
    SELECT isShow 
    FROM Review 
    WHERE id = $1`

	var isShow bool
	err := db.Pool.QueryRow(ctx, existsQuery, id).Scan(&isShow)
	if err != nil {
		if err == pgx.ErrNoRows {
			return storage.ErrReviewNotFound
		}
		return fmt.Errorf("%s: failed to check if review exists: %w", op, err)
	}

	updateQuery := `
    UPDATE Review 
    SET isShow = NOT isShow 
    WHERE id = $1`

	_, err = db.Pool.Exec(ctx, updateQuery, id)
	if err != nil {
		return fmt.Errorf("%s: failed to toggle visibility for review with ID %d: %w", op, id, err)
	}

	return nil
}

func (db *DB) GetAllReviewsAdmin(ctx context.Context) ([]*models.ReviewResponseAdmin, error) {
	const op = "postgres.GetAllReviews"

	query := `
		SELECT id, username, rating, text, created_at, isShow
		FROM Review
		ORDER BY created_at DESC 
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query reviews: %w", op, err)
	}
	defer rows.Close()

	var reviews []*models.ReviewResponseAdmin
	for rows.Next() {
		var review models.ReviewResponseAdmin
		if err := rows.Scan(&review.ID, &review.Username, &review.Rating, &review.Text, &review.CreatedAt, &review.IsShow); err != nil {
			return nil, fmt.Errorf("%s: failed to scan review row: %w", op, err)
		}
		reviews = append(reviews, &review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return reviews, nil
}
