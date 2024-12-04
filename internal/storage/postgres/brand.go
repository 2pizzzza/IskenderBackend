package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
)

func (db *DB) CreateBrand(ctx context.Context, name, url string) (*models.BrandResponse, error) {
	const op = "postgres.CreateBrand"

	var exists bool
	checkBrandQuery := `SELECT EXISTS(SELECT 1 FROM Brand WHERE name = $1)`
	err := db.Pool.QueryRow(ctx, checkBrandQuery, name).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check brand existence: %w", op, err)
	}
	if exists {
		return nil, storage.ErrBrandExists
	}

	var brandID int
	insertBrandQuery := `INSERT INTO Brand (name, url) VALUES ($1, $2) RETURNING id`
	err = db.Pool.QueryRow(ctx, insertBrandQuery, name, url).Scan(&brandID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert brand: %w", op, err)
	}

	response := &models.BrandResponse{
		ID:   brandID,
		Name: name,
		Url:  fmt.Sprintf("%s/%s", db.Config.BaseUrl, url),
	}

	return response, nil
}

func (db *DB) GetAllBrand(ctx context.Context) ([]*models.BrandResponse, error) {
	const op = "postgres.GetAllBrand"

	query := `SELECT id, name, url FROM Brand`
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query brands: %w", op, err)
	}
	defer rows.Close()

	var brands []*models.BrandResponse
	for rows.Next() {
		var brand models.BrandResponse
		if err := rows.Scan(&brand.ID, &brand.Name, &brand.Url); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row into brand struct: %w", op, err)
		}
		brand.Url = fmt.Sprintf("%s/%s", db.Config.BaseUrl, brand.Url)

		brands = append(brands, &brand)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}
	return brands, nil
}

func (db *DB) RemoveBrand(ctx context.Context, id int) error {
	const op = "postgres.RemoveBrand"

	var exists bool
	checkBrandQuery := `SELECT EXISTS(SELECT 1 FROM Brand WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkBrandQuery, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check brand existence: %w", op, err)
	}
	if !exists {
		return storage.ErrBrandNotFound
	}

	deleteQuery := `DELETE FROM Brand WHERE id = $1`
	_, err = db.Pool.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete brand: %w", op, err)
	}

	return nil
}

func (db *DB) UpdateBrand(ctx context.Context, id int, name, url string) (*models.BrandResponse, error) {
	const op = "postgres.UpdateBrand"

	var exists bool
	checkBrandQuery := `SELECT EXISTS(SELECT 1 FROM Brand WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkBrandQuery, id).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to check brand existence: %w", op, err)
	}
	if !exists {
		return nil, storage.ErrBrandNotFound
	}

	updateQuery := `UPDATE Brand SET name = $1, url = $2 WHERE id = $3 RETURNING id, name, url`
	var response models.BrandResponse
	err = db.Pool.QueryRow(ctx, updateQuery, name, url, id).Scan(&response.ID, &response.Name, &response.Url)
	response.Url = fmt.Sprintf("%s/%s", db.Config.BaseUrl, response.Url)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to update brand: %w", op, err)
	}

	return &response, nil
}

func (db *DB) GetBrandByID(ctx context.Context, id int) (*models.BrandResponse, error) {
	const op = "postgres.GetBrandByID"

	var brand models.BrandResponse
	query := `SELECT id, name, url FROM Brand WHERE id = $1`

	err := db.Pool.QueryRow(ctx, query, id).Scan(&brand.ID, &brand.Name, &brand.Url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrBrandNotFound
		}
		return nil, fmt.Errorf("%s: failed to query brand: %w", op, err)
	}
	brand.Url = fmt.Sprintf("%s/%s", db.Config.BaseUrl, brand.Url)

	return &brand, nil
}
