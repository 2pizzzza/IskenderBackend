package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
)

func (db *DB) GetItemsByCategoryID(ctx context.Context, categoryID int, languageCode string) ([]*models.ItemResponse, error) {
	const op = "postgres.GetItemsByCategoryIDWithDetails"

	var exists bool
	categoryQuery := `SELECT EXISTS(SELECT 1 FROM Category WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, categoryQuery, categoryID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check category existence: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: category with id %d not found: %w", op, categoryID, storage.ErrCategoryNotFound)
	}

	query := `
		SELECT i.id, COALESCE(it.name, ''), COALESCE(it.description, ''), i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $2
		WHERE i.category_id = $1`

	rows, err := db.Pool.Query(ctx, query, categoryID, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponse

	for rows.Next() {
		var item models.ItemResponse
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer, &item.IsPainted); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		photos, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		item.Photos = photos

		colors, err := db.getItemColors(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for item %d: %w", op, item.ID, err)
		}
		item.Colors = colors

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}
func (db *DB) GetItemByID(ctx context.Context, itemID int, languageCode string) (*models.ItemResponse, error) {
	const op = "postgres.GetItemByID"

	var exists bool
	existenceQuery := `SELECT EXISTS(SELECT 1 FROM Item WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, existenceQuery, itemID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check item existence: %w", op, err)
	}
	if !exists {
		return nil, storage.ErrItemNotFound
	}

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted,
		       COALESCE(it.name, ''), COALESCE(it.description, '')
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $2
		WHERE i.id = $1`

	var item models.ItemResponse
	err = db.Pool.QueryRow(ctx, query, itemID, languageCode).Scan(
		&item.ID, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer,
		&item.IsPainted, &item.Name, &item.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to retrieve item data: %w", op, err)
	}

	photos, err := db.getItemPhotos(ctx, item.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
	}
	item.Photos = photos

	colors, err := db.getItemColors(ctx, item.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get colors for item %d: %w", op, item.ID, err)
	}
	item.Colors = colors

	return &item, nil
}

func (db *DB) GetItemsByCollectionID(ctx context.Context, collectionID int, languageCode string) ([]*models.ItemResponse, error) {
	const op = "postgres.GetItemsByCollectionID"

	var exists bool
	existenceQuery := `SELECT EXISTS(SELECT 1 FROM Collection WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, existenceQuery, collectionID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check collection existence: %w", op, err)
	}
	if !exists {
		return nil, storage.ErrCollectionNotFound
	}

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew,
		       COALESCE(it.name, ''), COALESCE(it.description, '')
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $2
		WHERE i.collection_id = $1`

	rows, err := db.Pool.Query(ctx, query, collectionID, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponse

	for rows.Next() {
		var item models.ItemResponse
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer,
			&item.IsPainted, &item.IsPopular, &item.IsNew, &item.Name, &item.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		photos, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		item.Photos = photos

		colors, err := db.getItemColors(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for item %d: %w", op, item.ID, err)
		}
		item.Colors = colors

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}

func (db *DB) GetPopularItems(ctx context.Context, languageCode string) ([]*models.ItemResponse, error) {
	const op = "postgres.GetPopularItems"

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew,
		       COALESCE(it.name, ''), COALESCE(it.description, '')
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $1
		WHERE i.isPopular = TRUE`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query popular items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponse

	for rows.Next() {
		var item models.ItemResponse
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer,
			&item.IsPainted, &item.IsPopular, &item.IsNew, &item.Name, &item.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		photos, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		item.Photos = photos

		colors, err := db.getItemColors(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for item %d: %w", op, item.ID, err)
		}
		item.Colors = colors

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}

func (db *DB) GetNewItems(ctx context.Context, languageCode string) ([]*models.ItemResponse, error) {
	const op = "postgres.GetNewItems"

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew,
		       COALESCE(it.name, ''), COALESCE(it.description, '')
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $1
		WHERE i.isNew = TRUE`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query new items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponse

	for rows.Next() {
		var item models.ItemResponse
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer,
			&item.IsPainted, &item.IsPopular, &item.IsNew, &item.Name, &item.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		photos, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		item.Photos = photos

		colors, err := db.getItemColors(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for item %d: %w", op, item.ID, err)
		}
		item.Colors = colors

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}

func (db *DB) SearchItems(ctx context.Context, languageCode string, isProducer *bool, searchQuery string) ([]*models.ItemResponse, error) {
	const op = "postgres.SearchItems"

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew,
		       COALESCE(it.name, ''), COALESCE(it.description, '')
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $1
		WHERE 1=1`

	var args []interface{}
	args = append(args, languageCode)

	if isProducer != nil {
		query += ` AND i.isProducer = $2`
		args = append(args, *isProducer)
	}

	if searchQuery != "" {
		query += ` AND (it.name ILIKE $` + fmt.Sprintf("%d", len(args)+1) + ` OR it.description ILIKE $` + fmt.Sprintf("%d", len(args)+2) + `)`
		args = append(args, "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to search items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponse

	for rows.Next() {
		var item models.ItemResponse
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer,
			&item.IsPainted, &item.IsPopular, &item.IsNew, &item.Name, &item.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		photos, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		item.Photos = photos

		colors, err := db.getItemColors(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for item %d: %w", op, item.ID, err)
		}
		item.Colors = colors

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}

func (db *DB) getItemPhotos(ctx context.Context, itemID int) ([]models.PhotosResponse, error) {
	query := `
		SELECT p.id, p.url, p.isMain
		FROM ItemPhoto ip
		JOIN Photo p ON ip.photo_id = p.id
		WHERE ip.item_id = $1`

	rows, err := db.Pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("getItemPhotos: failed to query photos: %w", err)
	}
	defer rows.Close()

	var photos []models.PhotosResponse
	for rows.Next() {
		var photo models.PhotosResponse
		if err := rows.Scan(&photo.ID, &photo.URL, &photo.IsMain); err != nil {
			return nil, fmt.Errorf("getItemPhotos: failed to scan photo row: %w", err)
		}
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getItemPhotos: row iteration error: %w", err)
	}

	return photos, nil
}

func (db *DB) getItemColors(ctx context.Context, itemID int) ([]models.ColorResponse, error) {
	query := `
		SELECT c.id, c.hash_color
		FROM ItemColor ic
		JOIN Color c ON ic.color_id = c.id
		WHERE ic.item_id = $1`

	rows, err := db.Pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("getItemColors: failed to query colors: %w", err)
	}
	defer rows.Close()

	var colors []models.ColorResponse
	for rows.Next() {
		var color models.ColorResponse
		if err := rows.Scan(&color.ID, &color.HashColor); err != nil {
			return nil, fmt.Errorf("getItemColors: failed to scan color row: %w", err)
		}
		colors = append(colors, color)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getItemColors: row iteration error: %w", err)
	}

	return colors, nil
}
