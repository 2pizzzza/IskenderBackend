package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
)

func (db *DB) GetCategoriesByLanguageCode(ctx context.Context, languageCode string) ([]*models.Category, error) {
	const op = "postgres.GetCategoriesByLanguageCode"

	var languageCodeCheck string
	languageQuery := `SELECT code FROM Language WHERE code = $1`
	err := db.Pool.QueryRow(ctx, languageQuery, languageCode).Scan(&languageCodeCheck)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrLanguageNotFound
		}
		return nil, fmt.Errorf("%s: failed to check language code: %w", op, err)
	}

	query := `
		SELECT c.id, ct.name
		FROM Category c
		JOIN CategoryTranslation ct ON c.id = ct.category_id
		WHERE ct.language_code = $1`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query categories: %w", op, err)
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row into category struct: %w", op, err)
		}
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return categories, nil
}

func (db *DB) UpdateCategory(ctx context.Context, categoryID int, name string, languageCode string) error {
	const op = "postgres.UpdateCategory"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Category WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkQuery, categoryID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check category existence: %w", op, err)
	}
	if !exists {
		return storage.ErrCategoryNotFound
	}

	updateQuery := `
		INSERT INTO CategoryTranslation (category_id, language_code, name)
		VALUES ($1, $2, $3)
		ON CONFLICT (category_id, language_code) DO UPDATE
		SET name = EXCLUDED.name`
	_, err = db.Pool.Exec(ctx, updateQuery, categoryID, languageCode, name)
	if err != nil {
		return fmt.Errorf("%s: failed to update category: %w", op, err)
	}

	return nil
}

func (db *DB) CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.CreateCategoryResponse, error) {
	const op = "postgres.CreateCategory"

	if len(req.Categories) != 3 {
		return nil, storage.ErrRequiredLanguage
	}

	for _, cat := range req.Categories {
		var exists bool
		checkCategoryQuery := `SELECT EXISTS(SELECT 1 FROM CategoryTranslation WHERE name = $1 AND language_code = $2)`
		err := db.Pool.QueryRow(ctx, checkCategoryQuery, cat.Name, cat.LanguageCode).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to check category existence for language %s: %w", op, cat.LanguageCode, err)
		}
		if exists {
			return nil, storage.ErrCategoryExists
		}
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var categoryID int
	insertCategory := `INSERT INTO Category DEFAULT VALUES RETURNING id`
	err = tx.QueryRow(ctx, insertCategory).Scan(&categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert category: %w", op, err)
	}

	insertTranslation := `INSERT INTO CategoryTranslation (category_id, language_code, name) VALUES ($1, $2, $3)`
	for _, cat := range req.Categories {
		_, err = tx.Exec(ctx, insertTranslation, categoryID, cat.LanguageCode, cat.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert category translation for language %s: %w", op, cat.LanguageCode, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	var response models.CreateCategoryResponse
	for _, cat := range req.Categories {
		response.Categories = append(response.Categories, models.CategoriesResponse{
			ID:           categoryID,
			Name:         cat.Name,
			LanguageCode: cat.LanguageCode,
		})
	}

	return &response, nil
}

func (db *DB) DeleteCategory(ctx context.Context, categoryID int) error {
	const op = "postgres.DeleteCategory"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Category WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkQuery, categoryID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check category existence: %w", op, err)
	}
	if !exists {
		return storage.ErrCategoryNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	deleteTranslations := `DELETE FROM CategoryTranslation WHERE category_id = $1`
	_, err = tx.Exec(ctx, deleteTranslations, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete category translations: %w", op, err)
	}

	deleteItemTranslations := `DELETE FROM ItemTranslation WHERE item_id IN (SELECT id FROM Item WHERE category_id = $1)`
	_, err = tx.Exec(ctx, deleteItemTranslations, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item translations: %w", op, err)
	}

	deleteItemColors := `DELETE FROM ItemColor WHERE item_id IN (SELECT id FROM Item WHERE category_id = $1)`
	_, err = tx.Exec(ctx, deleteItemColors, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item colors: %w", op, err)
	}

	deleteItemPhotos := `DELETE FROM ItemPhoto WHERE item_id IN (SELECT id FROM Item WHERE category_id = $1)`
	_, err = tx.Exec(ctx, deleteItemPhotos, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item photos: %w", op, err)
	}

	deleteItems := `DELETE FROM Item WHERE category_id = $1`
	_, err = tx.Exec(ctx, deleteItems, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete items associated with category: %w", op, err)
	}

	deleteCategory := `DELETE FROM Category WHERE id = $1`
	_, err = tx.Exec(ctx, deleteCategory, categoryID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete category: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
