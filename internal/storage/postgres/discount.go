package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
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
		err := db.fillDiscountDetails(ctx, &discount, languageCode)
		if err != nil {
			slog.Error("Failed", sl.Err(err))
			return nil, storage.ErrDiscountExists
		}
		discounts = append(discounts, discount)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return discounts, nil
}

func (db *DB) CreateDiscount(ctx context.Context, discount models.DiscountCreate) (*models.DiscountCreate, error) {
	const op = "postgres.CreateDiscount"

	existsQuery := `
    SELECT EXISTS (
        SELECT 1 
        FROM Discount 
        WHERE discount_type = $1 AND target_id = $2 AND end_date >= NOW()
    )`

	var exists bool
	err := db.Pool.QueryRow(ctx, existsQuery, discount.DiscountType, discount.TargetID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check if discount exists: %w", op, err)
	}
	if exists {
		return nil, storage.ErrDiscountExists
	}

	insertQuery := `
    INSERT INTO Discount (discount_type, target_id, discount_percentage, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, discount_type, target_id, discount_percentage, start_date, end_date`

	var createdDiscount models.DiscountCreate
	err = db.Pool.QueryRow(ctx, insertQuery,
		discount.DiscountType,
		discount.TargetID,
		discount.DiscountPercentage,
		discount.StartDate,
		discount.EndDate,
	).Scan(
		&createdDiscount.ID,
		&createdDiscount.DiscountType,
		&createdDiscount.TargetID,
		&createdDiscount.DiscountPercentage,
		&createdDiscount.StartDate,
		&createdDiscount.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert new discount: %w", op, err)
	}

	return &createdDiscount, nil
}

func (db *DB) DeleteDiscount(ctx context.Context, id int) error {
	const op = "postgres.DeleteDiscount"

	existsQuery := `
    SELECT 1 
    FROM Discount 
    WHERE id = $1`

	var exists int
	err := db.Pool.QueryRow(ctx, existsQuery, id).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return storage.ErrDiscountNotFound
		}
		return fmt.Errorf("%s: failed to check if discount exists: %w", op, err)
	}

	deleteQuery := `
    DELETE FROM Discount 
    WHERE id = $1`

	_, err = db.Pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete discount with ID %d: %w", op, id, err)
	}

	return nil
}

func (db *DB) fillDiscountDetails(ctx context.Context, discount *models.Discount, languageCode string) error {
	const op = "postgres.FillDiscountDetails"

	discount.LanguageCode = languageCode

	switch discount.DiscountType {
	case "collection":
		collection, err := db.GetCollectionByID(ctx, discount.TargetID, languageCode)
		if err != nil {
			return fmt.Errorf("%s: failed to get collection by ID: %w", op, err)
		}

		discount.ID = collection.ID
		discount.Name = collection.Name
		discount.Description = collection.Description
		discount.IsPopular = collection.IsPopular
		discount.IsNew = collection.IsNew
		discount.IsProducer = collection.IsProducer
		discount.OldPrice = collection.Price
		discount.NewPrice = collection.Price - (collection.Price * discount.DiscountPercentage / 100)
		discount.Photo = collection.Photos
		discount.Color = collection.Colors

	case "item":
		item, err := db.GetItemByID(ctx, discount.TargetID, languageCode)
		if err != nil {
			return fmt.Errorf("%s: failed to get item by ID: %w", op, err)
		}

		discount.ID = item.ID
		discount.Name = item.Name
		discount.Description = item.Description
		discount.IsPopular = item.IsPopular
		discount.IsNew = item.IsNew
		discount.IsProducer = item.IsProducer
		discount.OldPrice = item.Price
		discount.NewPrice = item.Price - (item.Price * discount.DiscountPercentage / 100)
		discount.Photo = item.Photos
		discount.Color = item.Colors

	default:
		return fmt.Errorf("%s: unsupported discount type: %s", op, discount.DiscountType)
	}

	return nil
}
