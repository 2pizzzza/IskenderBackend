package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/jackc/pgx/v5"
)

func (db *DB) GetAllCategories(ctx context.Context) (*[]models.Category, error) {
	const op = "postgres.GetAllCategories"

	categoriesQuery := `SELECT category_id, name FROM Category`
	rows, err := db.Pool.Query(ctx, categoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query categories: %w", op, err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.CategoryID, &category.Name); err != nil {
			return nil, fmt.Errorf("%s: failed to scan category: %w", op, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return &categories, nil
}

func (db *DB) GetCategoryByID(ctx context.Context, categoryID int) (models.Category, error) {
	const op = "postgres.GetCategoryByID"

	categoryQuery := `SELECT category_id, name FROM Category WHERE category_id = $1`
	var category models.Category
	err := db.Pool.QueryRow(ctx, categoryQuery, categoryID).Scan(&category.CategoryID, &category.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Category{}, schemas.ErrItemNotFound
		}
		return models.Category{}, fmt.Errorf("%s: failed to scan category: %w", op, err)
	}

	return category, nil
}

func (db *DB) SaveCategory(ctx context.Context, name string) (models.Category, error) {
	const op = "postgres.SaveCategory"

	categoryQuery := `INSERT INTO Category (name) VALUES ($1) RETURNING category_id`
	var categoryID int
	err := db.Pool.QueryRow(ctx, categoryQuery, name).Scan(&categoryID)
	if err != nil {
		return models.Category{}, fmt.Errorf("%s: failed to create category: %w", op, err)
	}

	return models.Category{CategoryID: categoryID, Name: name}, nil
}

func (db *DB) UpdateCategory(ctx context.Context, categoryID int, name string) error {
	const op = "postgres.UpdateCategory"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Category WHERE category_id = $1)`
	err := db.Pool.QueryRow(ctx, checkQuery, categoryID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check category existence: %w", op, err)
	}

	if !exists {
		return schemas.ErrItemNotFound
	}

	updateQuery := `UPDATE Category SET name = $1 WHERE category_id = $2`
	_, err = db.Pool.Exec(ctx, updateQuery, name, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to update category: %w", op, err)
	}

	return nil
}

func (db *DB) RemoveCategory(ctx context.Context, categoryID int) error {
	const op = "postgres.RemoveCategory"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Category WHERE category_id = $1)`
	err := db.Pool.QueryRow(ctx, checkQuery, categoryID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check category existence: %w", op, err)
	}

	if !exists {
		return schemas.ErrItemNotFound
	}

	deleteQuery := `DELETE FROM Category WHERE category_id = $1`
	_, err = db.Pool.Exec(ctx, deleteQuery, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to remove category: %w", op, err)
	}

	return nil
}
