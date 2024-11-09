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
