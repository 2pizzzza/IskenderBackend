package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
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

		photos, colors, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
		item.NewPrice = newPrice
		item.Photos = photos
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

	photos, colors, err := db.getItemPhotos(ctx, item.ID)
	newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
	}
	item.NewPrice = newPrice
	item.Photos = photos
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

		photos, colors, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
		item.NewPrice = newPrice
		item.Photos = photos
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

		photos, colors, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
		item.NewPrice = newPrice
		item.Photos = photos
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

		photos, colors, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
		item.NewPrice = newPrice
		item.Photos = photos
		item.Colors = colors

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}

func (db *DB) SearchItems(ctx context.Context, languageCode string, isProducer *bool, isPainted *bool, searchQuery string) ([]*models.ItemResponse, error) {
	const op = "postgres.SearchItems"

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew,
		COALESCE(it.name, ''), COALESCE(it.description, '')
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $1
		WHERE it.name IS NOT NULL`

	var args []interface{}
	args = append(args, languageCode)

	if isProducer != nil {
		query += ` AND i.isProducer = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *isProducer)
	}

	if isPainted != nil {
		query += ` AND i.isPainted = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *isPainted)
	}

	if searchQuery != "" {
		query += ` AND it.name ILIKE $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+searchQuery+"%")
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

		photos, colors, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
		item.NewPrice = newPrice
		item.Photos = photos
		item.Colors = colors

		items = append(items, &item)
	}

	if len(items) == 0 {
		return nil, storage.ErrCollectionNotFound
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
}

func (db *DB) GetRandomItemsWithPopularity(ctx context.Context, languageCode string, itemID int) ([]*models.ItemResponse, error) {
	const op = "postgres.GetRandomItemsWithPopularity"

	var categoryID int
	getCategoryQuery := `SELECT category_id FROM Item WHERE id = $1`
	err := db.Pool.QueryRow(ctx, getCategoryQuery, itemID).Scan(&categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: item with id %d not found", op, itemID)
		}
		return nil, fmt.Errorf("%s: failed to get category id for item %d: %w", op, itemID, err)
	}

	query := `
		SELECT i.id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew, it.name, it.description
		FROM Item i
		LEFT JOIN ItemTranslation it ON i.id = it.item_id AND it.language_code = $1
		WHERE i.category_id = $2
		ORDER BY RANDOM() 
		LIMIT 7`

	rows, err := db.Pool.Query(ctx, query, languageCode, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query random items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponse

	for rows.Next() {
		var item models.ItemResponse
		var name sql.NullString
		var description sql.NullString

		if err := rows.Scan(
			&item.ID,
			&item.Size,
			&item.Price,
			&item.IsProducer,
			&item.IsPainted,
			&item.IsPopular,
			&item.IsNew,
			&name,
			&description,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		if name.Valid {
			item.Name = name.String
		} else {
			item.Name = ""
		}

		if description.Valid {
			item.Description = description.String
		} else {
			item.Description = ""
		}

		photos, colors, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for item %d: %w", op, item.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "item", item.ID, item.Price)
		item.NewPrice = newPrice
		item.Photos = photos
		item.Colors = colors

		items = append(items, &item)
	}

	var popularItems []*models.ItemResponse
	var regularItems []*models.ItemResponse

	for _, item := range items {
		if item.IsPopular {
			popularItems = append(popularItems, item)
		} else {
			regularItems = append(regularItems, item)
		}
	}

	items = append(popularItems, regularItems...)

	return items, nil
}

func (db *DB) getItemPhotos(ctx context.Context, itemID int) ([]models.PhotosResponse, []models.ColorResponse, error) {
	query := `
		SELECT p.id, p.url, p.isMain, p.hash_color
		FROM ItemPhoto ip
		JOIN Photo p ON ip.photo_id = p.id
		WHERE ip.item_id = $1`

	baseURL := fmt.Sprintf("http://%s:%d", db.Config.HttpHost, db.Config.HttpPort)
	rows, err := db.Pool.Query(ctx, query, itemID)
	if err != nil {
		return nil, nil, fmt.Errorf("getItemPhotos: failed to query photos: %w", err)
	}
	defer rows.Close()
	var colors []models.ColorResponse
	var photos []models.PhotosResponse
	for rows.Next() {
		var photo models.PhotosResponse
		if err := rows.Scan(&photo.ID, &photo.URL, &photo.IsMain, &photo.HashColor); err != nil {
			return nil, nil, fmt.Errorf("getItemPhotos: failed to scan photo row: %w", err)
		}
		colors = append(colors, models.ColorResponse{HashColor: photo.HashColor})
		photo.URL = fmt.Sprintf("%s/%s", baseURL, photo.URL)
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("getItemPhotos: row iteration error: %w", err)
	}

	return photos, colors, nil
}

func (db *DB) GetDiscountedPrice(ctx context.Context, targetType string, targetID int, oldPrice float64) (float64, error) {
	const op = "postgres.GetDiscountedPrice"

	if targetType != "item" && targetType != "collection" {
		return 0.0, fmt.Errorf("%s: invalid targetType %s, must be 'item' or 'collection'", op, targetType)
	}

	var discountPercentage float64

	var query string
	if targetType == "item" {
		query = `SELECT discount_percentage FROM Discount d
                 JOIN Item i ON d.target_id = i.id
                 WHERE d.discount_type = 'item' AND d.target_id = $1`
	} else if targetType == "collection" {
		query = `SELECT discount_percentage FROM Discount d
                 JOIN Collection c ON d.target_id = c.id
                 WHERE d.discount_type = 'collection' AND d.target_id = $1`
	}

	err := db.Pool.QueryRow(ctx, query, targetID).Scan(&discountPercentage)
	if err != nil {
		if err == sql.ErrNoRows {
			return oldPrice, nil
		}
		return 0.0, fmt.Errorf("%s: failed to get discount percentage for targetID %d: %w", op, targetID, err)
	}
	if oldPrice <= 0 {
		return 0.0, nil
	}
	discountFactor := 1 - (discountPercentage / 100)
	newPrice := oldPrice * discountFactor

	return newPrice, nil
}
