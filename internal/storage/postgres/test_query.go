package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/storage"
)

func (db *DB) CreateStarter(ctx context.Context) error {
	const op = "postgres.CreateStarter"

	var exists bool
	checkLanguageQuery := `SELECT EXISTS(SELECT 1 FROM Language WHERE code IN ('ru', 'kgz', 'en'))`
	err := db.Pool.QueryRow(ctx, checkLanguageQuery).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check if languages already exist: %w", op, err)
	}
	if exists {
		return storage.ErrCategoryNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO Language (code, name) VALUES 
		('ru', 'Русский'),
		('kgz', 'Кырғызча'),
		('en', 'English')`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert languages: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
