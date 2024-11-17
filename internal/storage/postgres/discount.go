package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"log/slog"
)

func (db *DB) GetAllDiscount(ctx context.Context, languageCode string) ([]models.Discount, error) {
	const op = "postgres.GetAllDiscount"

	query := `
    SELECT EXISTS (
        SELECT 1
        FROM Language
        WHERE code = $1
    )`

	var exists bool
	err := db.Pool.QueryRow(ctx, query, languageCode).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check language existence: %w", op, err)
	}

	if !exists {
		return nil, storage.ErrLanguageNotFound
	}

	query = `
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
		discount.LanguageCode = languageCode
		if discount.DiscountType == "collection" {
			collection, err := db.GetCollectionByID(ctx, discount.TargetID, languageCode)
			if err != nil {
				return nil, fmt.Errorf("%s, %w", op, err)
			}
			slog.Info("collection", collection)
			discount.ID = collection.ID
			discount.Name = collection.Name
			discount.Description = collection.Description
			discount.IsPopular = collection.IsPopular
			discount.IsNew = collection.IsNew
			discount.IsProducer = collection.IsProducer
			discount.OldPrice = collection.Price
			discount.NewPrice = collection.Price
			discount.Photo = collection.Photos
			discount.Color = collection.Colors
		}

		if discount.DiscountType == "item" {
			item, err := db.GetItemByID(ctx, discount.TargetID, languageCode)
			slog.Info("collection", item)
			if err != nil {
				return nil, fmt.Errorf("%s, %w", op, err)
			}
			discount.ID = item.ID
			discount.Name = item.Name
			discount.Description = item.Description
			discount.IsPopular = item.IsPopular
			discount.IsNew = item.IsNew
			discount.IsProducer = item.IsProducer
			discount.OldPrice = item.Price
			discount.NewPrice = item.Price
			discount.Photo = item.Photos
			discount.Color = item.Colors
		}
		discounts = append(discounts, discount)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return discounts, nil
}
