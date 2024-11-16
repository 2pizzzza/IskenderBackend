package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
)

func (db *DB) GetAllDiscount(ctx context.Context) ([]models.Discount, error) {
	const op = "postgres.GetAllDiscount"

	query := `
    SELECT id, discount_type, target_id, discount_percentage, start_date, end_date
    FROM Discount
    WHERE start_date <= NOW() AND end_date >= NOW()`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query discounts: %w", op, err)
	}
	defer rows.Close()

	var discounts []models.Discount

	for rows.Next() {
		var discount models.Discount
		if err := rows.Scan(
			&discount.ID,
			&discount.DiscountType,
			&discount.TargetID,
			&discount.DiscountPercentage,
			&discount.StartDate,
			&discount.EndDate,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan discount row: %w", op, err)
		}

		discounts = append(discounts, discount)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return discounts, nil
}
