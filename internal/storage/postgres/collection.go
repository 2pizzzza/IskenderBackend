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

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
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

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
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

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
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

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
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

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
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
	newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
	collection.NewPrice = newPrice
	photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
	}
	collection.Photos = photos
	collection.Colors = colors

	return &collection, nil
}
func (db *DB) SearchCollections(ctx context.Context, languageCode string, isProducer *bool, isPainted *bool, searchQuery string) ([]*models.CollectionResponse, error) {
	const op = "postgres.SearchCollections"

	query := `
		SELECT c.id, c.price, c.isProducer, c.isPainted, c.isPopular, c.isNew,
		COALESCE(ct.name, ''), COALESCE(ct.description, '')
		FROM Collection c
		LEFT JOIN CollectionTranslation ct ON c.id = ct.collection_id AND ct.language_code = $1
		WHERE 1=1`

	var args []interface{}
	args = append(args, languageCode)

	if isProducer != nil {
		query += ` AND c.isProducer = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *isProducer)
	}

	if isPainted != nil {
		query += ` AND c.isPainted = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *isPainted)
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
		if err := rows.Scan(&collection.ID, &collection.Price, &collection.IsProducer, &collection.IsPainted, &collection.IsPopular,
			&collection.IsNew, &collection.Name, &collection.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan collection row: %w", op, err)
		}

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
		collection.Colors = colors

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

		photos, colors, err := db.getCollectionPhotos(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get photos for collection %d: %w", op, collection.ID, err)
		}
		newPrice, err := db.GetDiscountedPrice(ctx, "collection", collection.ID, collection.Price)
		collection.NewPrice = newPrice
		collection.Photos = photos
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

func (db *DB) getCollectionPhotos(ctx context.Context, collectionID int) ([]models.PhotosResponse, []models.ColorResponse, error) {
	query := `
		SELECT p.id, p.url, p.isMain, p.hash_color
		FROM CollectionPhoto cp
		JOIN Photo p ON cp.photo_id = p.id
		WHERE cp.collection_id = $1`

	baseURL := fmt.Sprintf("http://%s:%d", db.Config.HttpHost, db.Config.HttpPort)
	rows, err := db.Pool.Query(ctx, query, collectionID)
	if err != nil {
		return nil, nil, fmt.Errorf("getCollectionPhotos: failed to query photos: %w", err)
	}
	defer rows.Close()

	var photos []models.PhotosResponse
	var colors []models.ColorResponse
	for rows.Next() {
		var photo models.PhotosResponse
		if err := rows.Scan(&photo.ID, &photo.URL, &photo.IsMain, &photo.HashColor); err != nil {
			return nil, nil, fmt.Errorf("getCollectionPhotos: failed to scan photo row: %w", err)
		}
		colors = append(colors, models.ColorResponse{HashColor: photo.HashColor})
		photo.URL = fmt.Sprintf("%s/%s", baseURL, photo.URL)
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("getCollectionPhotos: row iteration error: %w", err)
	}

	return photos, colors, nil
}

func (db *DB) CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.CreateCollectionResponse, error) {
	const op = "postgres.CreateCollection"

	if len(req.Collections) != 3 {
		return nil, storage.ErrRequiredLanguage
	}

	languageCodes := map[string]bool{"ru": false, "kgz": false, "en": false}
	for _, translation := range req.Collections {
		if _, ok := languageCodes[translation.LanguageCode]; !ok {
			return nil, storage.ErrInvalidLanguageCode
		}
		languageCodes[translation.LanguageCode] = true
	}

	for _, translation := range req.Collections {
		var exists bool
		checkCollectionQuery := `SELECT EXISTS(
			SELECT 1 FROM CollectionTranslation 
			WHERE name = $1 AND language_code = $2
		)`
		err := db.Pool.QueryRow(ctx, checkCollectionQuery, translation.Name, translation.LanguageCode).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to check collection existence for language %s: %w", op, translation.LanguageCode, err)
		}
		if exists {
			return nil, storage.ErrCollectionExists
		}
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var collectionID int
	insertCollectionQuery := `
		INSERT INTO Collection (price, isProducer, isPainted, isPopular, isNew)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = tx.QueryRow(ctx, insertCollectionQuery,
		req.Price,
		req.IsProducer,
		req.IsPainted,
		req.IsPopular,
		req.IsNew,
	).Scan(&collectionID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert collection: %w", op, err)
	}

	insertTranslationQuery := `
		INSERT INTO CollectionTranslation (collection_id, language_code, name, description)
		VALUES ($1, $2, $3, $4)
	`
	for _, translation := range req.Collections {
		_, err = tx.Exec(ctx, insertTranslationQuery, collectionID, translation.LanguageCode, translation.Name, translation.Description)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert collection translation for language %s: %w", op, translation.LanguageCode, err)
		}
	}

	insertPhotoQuery := `
		INSERT INTO Photo (url, isMain, hash_color)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var photoIDs []int
	for _, photo := range req.Photos {
		var photoID int
		err = tx.QueryRow(ctx, insertPhotoQuery, photo.URL, photo.IsMain, photo.HashColor).Scan(&photoID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert photo with url %s: %w", op, photo.URL, err)
		}
		photoIDs = append(photoIDs, photoID)
	}

	insertCollectionPhotoQuery := `
		INSERT INTO CollectionPhoto (collection_id, photo_id)
		VALUES ($1, $2)
	`
	for _, photoID := range photoIDs {
		_, err = tx.Exec(ctx, insertCollectionPhotoQuery, collectionID, photoID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert photo association for photo id %d: %w", op, photoID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	var response models.CreateCollectionResponse
	response.ID = collectionID
	response.Price = req.Price
	response.IsProducer = req.IsProducer
	response.IsPainted = req.IsPainted
	response.IsPopular = req.IsPopular
	response.IsNew = req.IsNew
	response.Collections = req.Collections

	for _, photo := range req.Photos {
		response.Photos = append(response.Photos, models.PhotosResponse{
			URL:       photo.URL,
			IsMain:    photo.IsMain,
			HashColor: photo.HashColor,
		})
	}

	return &response, nil
}
