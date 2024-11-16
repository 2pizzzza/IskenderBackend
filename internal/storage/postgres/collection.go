package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
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

func (db *DB) GetCollectionsByIsProducerSLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error) {
	const op = "postgres.GetCollectionsByLanguageCode"

	query := `
    SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew, 
           COALESCE(ct.name, ''), COALESCE(ct.description, '')
    FROM Collection c
    LEFT JOIN CollectionTranslation ct 
    ON c.id = ct.collection_id AND ct.language_code = $1
    WHERE c.isProducer = true AND c.isPainted = false`

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

func (db *DB) GetCollectionsByIsProducerPLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error) {
	const op = "postgres.GetCollectionsByLanguageCode"

	query := `
    SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew, 
           COALESCE(ct.name, ''), COALESCE(ct.description, '')
    FROM Collection c
    LEFT JOIN CollectionTranslation ct 
    ON c.id = ct.collection_id AND ct.language_code = $1
    WHERE c.isProducer = true AND c.isPainted = false`

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

func (db *DB) SearchCollections(ctx context.Context, languageCode string, isProducer *bool, searchQuery string) ([]*models.CollectionResponse, error) {
	const op = "postgres.SearchCollections"

	query := `
		SELECT ct.name
		FROM Collection c
		LEFT JOIN CollectionTranslation ct ON c.id = ct.collection_id AND ct.language_code = $1
		WHERE 1=1`

	var args []interface{}
	args = append(args, languageCode)

	if isProducer != nil {
		query += ` AND c.isProducer = $2`
		args = append(args, *isProducer)
	}

	if searchQuery != "" {
		query += ` AND ct.name ILIKE $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+searchQuery+"%")
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to search collections: %w", op, err)
	}
	defer rows.Close()

	var collections []*models.CollectionResponse

	for rows.Next() {
		var collection models.CollectionResponse
		if err := rows.Scan(&collection.Name); err != nil {
			return nil, fmt.Errorf("%s: failed to scan collection row: %w", op, err)
		}

		collections = append(collections, &collection)
	}

	if len(collections) == 0 {
		return nil, storage.ErrCollectionNotFound
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return collections, nil
}

func (db *DB) DeleteCollection(ctx context.Context, collectionID int) error {
	const op = "postgres.DeleteCollection"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Collection WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkQuery, collectionID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check collection existence: %w", op, err)
	}
	if !exists {
		return storage.ErrCollectionNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	deleteColorCollection := `DELETE FROM ColorCollection WHERE collection_id = $1`
	_, err = tx.Exec(ctx, deleteColorCollection, collectionID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete color associations for collection: %w", op, err)
	}

	deletePhotoCollection := `DELETE FROM PhotoCollection WHERE collection_id = $1`
	_, err = tx.Exec(ctx, deletePhotoCollection, collectionID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete photo associations for collection: %w", op, err)
	}

	deleteCollection := `DELETE FROM Collection WHERE id = $1`
	_, err = tx.Exec(ctx, deleteCollection, collectionID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete collection: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) UpdateCollection(ctx context.Context, req *models.UpdateCollectionRequest) error {
	const op = "postgres.UpdateCollection"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM Collection WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkQuery, req.CollectionID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check collection existence: %w", op, err)
	}
	if !exists {
		return storage.ErrCollectionNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	updateCollection := `
		UPDATE Collection
		SET price = COALESCE(NULLIF($1, 0), price),
		    isProducer = COALESCE($2, isProducer),
		    isPainted = COALESCE($3, isPainted),
		    isPopular = COALESCE($4, isPopular),
		    isNew = COALESCE($5, isNew)
		WHERE id = $6`
	_, err = tx.Exec(ctx, updateCollection, req.Price, req.IsProducer, req.IsPainted, req.IsPopular, req.IsNew, req.CollectionID)
	if err != nil {
		return fmt.Errorf("%s: failed to update collection: %w", op, err)
	}

	for _, translation := range req.Collections {
		updateTranslation := `
			INSERT INTO CollectionTranslation (collection_id, language_code, name, description)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (collection_id, language_code)
			DO UPDATE SET name = COALESCE($3, CollectionTranslation.name),
			              description = COALESCE($4, CollectionTranslation.description)`
		_, err := tx.Exec(ctx, updateTranslation, req.CollectionID, translation.LanguageCode, translation.Name, translation.Description)
		if err != nil {
			return fmt.Errorf("%s: failed to update collection translation for language %s: %w", op, translation.LanguageCode, err)
		}
	}

	if err = db.updateCollectionPhotos(ctx, tx, req.CollectionID, req.Photos); err != nil {
		return fmt.Errorf("%s: failed to update collection photos: %w", op, err)
	}

	if err = db.updateCollectionColors(ctx, tx, req.CollectionID, req.Colors); err != nil {
		return fmt.Errorf("%s: failed to update collection colors: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) GetRandomCollectionsWithPopularity(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error) {
	const op = "postgres.GetRandomCollectionsWithPopularity"

	query := `
		SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew, ct.name, ct.description
		FROM Collection c
		LEFT JOIN CollectionTranslation ct ON c.id = ct.collection_id AND ct.language_code = $1
		ORDER BY RANDOM() 
		LIMIT 7`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query random collections: %w", op, err)
	}
	defer rows.Close()

	var collections []*models.CollectionResponse

	for rows.Next() {
		var collection models.CollectionResponse
		if err := rows.Scan(
			&collection.ID,
			&collection.Price,
			&collection.IsProducer,
			&collection.IsPainted,
			&collection.IsPopular,
			&collection.IsNew,
			&collection.Name,
			&collection.Description,
		); err != nil {
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

	var popularCollections []*models.CollectionResponse
	var regularCollections []*models.CollectionResponse

	for _, collection := range collections {
		if collection.IsPopular {
			popularCollections = append(popularCollections, collection)
		} else {
			regularCollections = append(regularCollections, collection)
		}
	}

	collections = append(popularCollections, regularCollections...)

	return collections, nil
}

func (db *DB) updateCollectionPhotos(ctx context.Context, tx pgx.Tx, collectionID int, photos []models.PhotosResponse) error {
	// Удаление старых фотографий
	deletePhotos := `DELETE FROM CollectionPhoto WHERE collection_id = $1`
	_, err := tx.Exec(ctx, deletePhotos, collectionID)
	if err != nil {
		return fmt.Errorf("failed to delete existing photos: %w", err)
	}

	insertPhoto := `INSERT INTO CollectionPhoto (collection_id, photo_id) VALUES ($1, $2)`
	for _, photo := range photos {
		_, err = tx.Exec(ctx, insertPhoto, collectionID, photo.ID)
		if err != nil {
			return fmt.Errorf("failed to insert new photo: %w", err)
		}
	}

	return nil
}

func (db *DB) updateCollectionColors(ctx context.Context, tx pgx.Tx, collectionID int, colors []models.ColorResponse) error {
	deleteColors := `DELETE FROM CollectionColor WHERE collection_id = $1`
	_, err := tx.Exec(ctx, deleteColors, collectionID)
	if err != nil {
		return fmt.Errorf("failed to delete existing colors: %w", err)
	}

	insertColor := `INSERT INTO CollectionColor (collection_id, color_id) VALUES ($1, $2)`
	for _, color := range colors {
		_, err = tx.Exec(ctx, insertColor, collectionID, color.ID)
		if err != nil {
			return fmt.Errorf("failed to insert new color: %w", err)
		}
	}

	return nil
}

func (db *DB) getCollectionPhotos(ctx context.Context, collectionID int) ([]models.PhotosResponse, error) {
	query := `
		SELECT p.id, p.url, p.isMain
		FROM CollectionPhoto cp
		JOIN Photo p ON cp.photo_id = p.id
		WHERE cp.collection_id = $1`

	baseURL := fmt.Sprintf("http://%s:%d", db.Config.HttpHost, db.Config.HttpPort)
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
		photo.URL = fmt.Sprintf("%s/%s", baseURL, photo.URL)
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
