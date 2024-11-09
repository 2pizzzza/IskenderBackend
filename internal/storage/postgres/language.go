package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
)

func (db *DB) GetLanguages(ctx context.Context) ([]*models.Language, error) {
	const op = "postgres.GetLanguages"

	query := `SELECT code, name FROM Language`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query to get languages: %w", op, err)
	}
	defer rows.Close()

	var languages []*models.Language

	for rows.Next() {
		var language models.Language
		if err := rows.Scan(&language.Code, &language.Name); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row into language struct: %w", op, err)
		}
		languages = append(languages, &language)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return languages, nil
}
