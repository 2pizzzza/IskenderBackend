package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
)

func (db *DB) GetCollectionsByLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error) {
	const op = "postgres.GetCollectionsByLanguageCode"

	query := `
    SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew, 
           COALESCE(ct.name, ''), COALESCE(ct.description, '')
    FROM Collection c
    LEFT JOIN CollectionTranslation ct 
    ON c.id = ct.collection_id AND ct.language_code = $1`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query collections: %w", op, err)
	}
	defer rows.Close()

	var collections []*models.CollectionResponse

	for rows.Next() {
		var collection models.CollectionResponse
		if err := rows.Scan(&collection.ID, &collection.Price, &collection.IsProducer, &collection.IsPainted, &collection.IsPopular, &collection.IsNew, &collection.Name, &collection.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan collection row: %w", op, err)
		}

		photos, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		collection.Photos = photos

		colors, err := db.getCollectionColors(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for collection %d: %w", op, collection.ID, err)
		}
		collection.Colors = colors

		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return collections, nil
}

func (db *DB) GetPopularCollections(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error) {
	const op = "postgres.GetPopularCollections"

	query := `
		SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew,
		       COALESCE(ct.name, ''), COALESCE(ct.description, '')
		FROM Collection c
		LEFT JOIN CollectionTranslation ct ON c.id = ct.collection_id AND ct.language_code = $1
		WHERE c.isPopular = TRUE`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query popular collections: %w", op, err)
	}
	defer rows.Close()

	var collections []*models.CollectionResponse

	for rows.Next() {
		var collection models.CollectionResponse
		if err := rows.Scan(&collection.ID, &collection.Price, &collection.IsProducer, &collection.IsPainted, &collection.IsPopular,
			&collection.IsNew, &collection.Name, &collection.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan collection row: %w", op, err)
		}

		photos, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		collection.Photos = photos

		colors, err := db.getCollectionColors(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for collection %d: %w", op, collection.ID, err)
		}
		collection.Colors = colors

		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return collections, nil
}

func (db *DB) GetNewCollections(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error) {
	const op = "postgres.GetNewCollections"

	query := `
		SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew,
		       COALESCE(ct.name, ''), COALESCE(ct.description, '')
		FROM Collection c
		LEFT JOIN CollectionTranslation ct ON c.id = ct.collection_id AND ct.language_code = $1
		WHERE c.isNew = TRUE`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query new collections: %w", op, err)
	}
	defer rows.Close()

	var collections []*models.CollectionResponse

	for rows.Next() {
		var collection models.CollectionResponse
		if err := rows.Scan(&collection.ID, &collection.Price, &collection.IsProducer, &collection.IsPainted, &collection.IsPopular,
			&collection.IsNew, &collection.Name, &collection.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan collection row: %w", op, err)
		}

		photos, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		collection.Photos = photos

		colors, err := db.getCollectionColors(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get colors for collection %d: %w", op, collection.ID, err)
		}
		collection.Colors = colors

		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return collections, nil
}

func (db *DB) GetCollectionByID(ctx context.Context, collectionID int, languageCode string) (*models.CollectionResponse, error) {
	const op = "postgres.GetCollectionByID"

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
    SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew, 
           COALESCE(ct.name, ''), COALESCE(ct.description, '')
    FROM Collection c
    LEFT JOIN CollectionTranslation ct 
    ON c.id = ct.collection_id AND ct.language_code = $2
    WHERE c.id = $1`

	var collection models.CollectionResponse
	err = db.Pool.QueryRow(ctx, query, collectionID, languageCode).Scan(
		&collection.ID, &collection.Price, &collection.IsProducer, &collection.IsPainted, &collection.IsPopular, &collection.IsNew,
		&collection.Name, &collection.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to retrieve collection data: %w", op, err)
	}

	photos, err := db.getCollectionPhotos(ctx, collection.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
	}
	collection.Photos = photos

	colors, err := db.getCollectionColors(ctx, collection.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get colors for collection %d: %w", op, collection.ID, err)
	}
	collection.Colors = colors

	return &collection, nil
}

func (db *DB) getCollectionPhotos(ctx context.Context, collectionID int) ([]models.PhotosResponse, error) {
	query := `
		SELECT p.id, p.url, p.isMain
		FROM CollectionPhoto cp
		JOIN Photo p ON cp.photo_id = p.id
		WHERE cp.collection_id = $1`

	rows, err := db.Pool.Query(ctx, query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("getCollectionPhotos: failed to query photos: %w", err)
	}
	defer rows.Close()

	var photos []models.PhotosResponse
	for rows.Next() {
		var photo models.PhotosResponse
		if err := rows.Scan(&photo.ID, &photo.URL, &photo.IsMain); err != nil {
			return nil, fmt.Errorf("getCollectionPhotos: failed to scan photo row: %w", err)
		}
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCollectionPhotos: row iteration error: %w", err)
	}

	return photos, nil
}

func (db *DB) getCollectionColors(ctx context.Context, collectionID int) ([]models.ColorResponse, error) {
	query := `
		SELECT c.id, c.hash_color
		FROM CollectionColor cc
		JOIN Color c ON cc.color_id = c.id
		WHERE cc.collection_id = $1`

	rows, err := db.Pool.Query(ctx, query, collectionID)
	if err != nil {
		return nil, fmt.Errorf("getCollectionColors: failed to query colors: %w", err)
	}
	defer rows.Close()

	var colors []models.ColorResponse
	for rows.Next() {
		var color models.ColorResponse
		if err := rows.Scan(&color.ID, &color.HashColor); err != nil {
			return nil, fmt.Errorf("getCollectionColors: failed to scan color row: %w", err)
		}
		colors = append(colors, color)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getCollectionColors: row iteration error: %w", err)
	}

	return colors, nil
}
