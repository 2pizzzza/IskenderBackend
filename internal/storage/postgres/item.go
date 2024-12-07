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
	item.NewPrice = newPrice
	item.Photos = photos
	item.Colors = colors
	return &item, nil
}

func (db *DB) GetAllItems(ctx context.Context) ([]*models.ItemResponses, error) {
	const op = "postgres.GetAllItems"

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, i.isProducer, i.isPainted, i.isPopular, i.isNew
		FROM Item i`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query items: %w", op, err)
	}
	defer rows.Close()

	var items []*models.ItemResponses

	for rows.Next() {
		var item models.ItemResponses
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.CollectionID, &item.Size, &item.Price, &item.IsProducer,
			&item.IsPainted, &item.IsPopular, &item.IsNew); err != nil {
			return nil, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		transQuery := `SELECT language_code, name, description FROM ItemTranslation WHERE item_id = $1`
		transRows, err := db.Pool.Query(ctx, transQuery, item.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query translations for item %d: %w", op, item.ID, err)
		}
		defer transRows.Close()

		var translations []models.CreateItemTranslation
		for transRows.Next() {
			var trans models.CreateItemTranslation
			if err := transRows.Scan(&trans.LanguageCode, &trans.Name, &trans.Description); err != nil {
				return nil, fmt.Errorf("%s: failed to scan translation row: %w", op, err)
			}
			if trans.LanguageCode == "ru" {
				item.Name = trans.Name
			}
			translations = append(translations, trans)
		}

		photos, color, err := db.getItemPhotos(ctx, item.ID)
		if err != nil {
		}

		item.Items = translations
		item.Photos = photos
		item.Color = color

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return items, nil
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

func (db *DB) GetItem(ctx context.Context, itemID int) (*models.ItemResponseForAdmin, error) {
	const op = "postgres.GetItem"

	var exist bool
	checkItemQuery := `SELECT EXISTS(SELECT 1 FROM Item WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkItemQuery, itemID).Scan(&exist)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check item existence: %w", op, err)
	}
	if !exist {
		return nil, storage.ErrItemNotFound
	}

	query := `
		SELECT i.id, i.category_id, i.collection_id, i.size, i.price, 
			i.isProducer, i.isPainted, i.isPopular, i.isNew
		FROM Item i WHERE i.id = $1`
	var item models.ItemResponseForAdmin
	err = db.Pool.QueryRow(ctx, query, itemID).Scan(
		&item.ID,
		&item.CategoryID,
		&item.CollectionID,
		&item.Size,
		&item.Price,
		&item.IsProducer,
		&item.IsPainted,
		&item.IsPopular,
		&item.IsNew,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get item details: %w", op, err)
	}

	// Запрос для получения переводов товара
	transQuery := `SELECT language_code, name, description FROM ItemTranslation WHERE item_id = $1`
	rows, err := db.Pool.Query(ctx, transQuery, itemID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get item translations: %w", op, err)
	}
	defer rows.Close()

	var translations []models.CreateItemTranslation
	for rows.Next() {
		var translation models.CreateItemTranslation
		if err := rows.Scan(&translation.LanguageCode, &translation.Name, &translation.Description); err != nil {
			return nil, fmt.Errorf("%s: failed to scan translation: %w", op, err)
		}
		translations = append(translations, translation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating over translations: %w", op, err)
	}
	item.Items = translations

	photosQuery := `
		SELECT p.id, p.url, p.isMain, p.hash_color
		FROM Photo p
		JOIN ItemPhoto ip ON p.id = ip.photo_id
		WHERE ip.item_id = $1`
	photoRows, err := db.Pool.Query(ctx, photosQuery, itemID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get item photos: %w", op, err)
	}
	defer photoRows.Close()

	var photos []models.PhotosResponse
	for photoRows.Next() {
		var photo models.PhotosResponse
		if err := photoRows.Scan(&photo.ID, &photo.URL, &photo.IsMain, &photo.HashColor); err != nil {
			return nil, fmt.Errorf("%s: failed to scan photo: %w", op, err)
		}
		photos = append(photos, photo)
	}
	if err := photoRows.Err(); err != nil {
		return nil, fmt.Errorf("%s: error iterating over photos: %w", op, err)
	}
	item.Photos = photos

	return &item, nil
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

func (db *DB) SearchItems(ctx context.Context, languageCode string, isProducer *bool, isPainted *bool, searchQuery string, minPrice, maxPrice *float64) ([]*models.ItemResponse, error) {
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

	if minPrice != nil {
		query += ` AND i.price >= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *minPrice)
	}

	if maxPrice != nil {
		query += ` AND i.price <= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *maxPrice)
	}
	fmt.Printf("Executing query: %s\nWith arguments: %v\n", query, args)

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
		photo.URL = fmt.Sprintf("%s/%s", db.Config.BaseUrl, photo.URL)
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

func (db *DB) CreateItem(ctx context.Context, req models.CreateItem) (*models.CreateItemResponse, error) {
	const op = "postgres.CreateItem"

	if len(req.Items) != 3 {
		return nil, storage.ErrRequiredLanguage
	}

	languageCodes := map[string]bool{"ru": false, "kgz": false, "en": false}
	for _, translation := range req.Items {
		if _, ok := languageCodes[translation.LanguageCode]; !ok {
			return nil, storage.ErrInvalidLanguageCode
		}
		languageCodes[translation.LanguageCode] = true
	}

	for _, translation := range req.Items {
		var exists bool
		checkItemQuery := `SELECT EXISTS(
			SELECT 1 FROM ItemTranslation 
			WHERE name = $1 AND language_code = $2
		)`
		err := db.Pool.QueryRow(ctx, checkItemQuery, translation.Name, translation.LanguageCode).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to check item existence for language %s: %w", op, translation.LanguageCode, err)
		}
		if exists {
			return nil, storage.ErrItemExists
		}
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var itemID int
	insertItemQuery := `
		INSERT INTO Item (category_id, collection_id, size, price, isProducer, isPainted, isPopular, isNew)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err = tx.QueryRow(ctx, insertItemQuery,
		req.CategoryID,
		req.CollectionID,
		req.Size,
		req.Price,
		req.IsProducer,
		req.IsPainted,
		req.IsPopular,
		req.IsNew,
	).Scan(&itemID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert item: %w", op, err)
	}

	insertTranslationQuery := `
		INSERT INTO ItemTranslation (item_id, language_code, name, description)
		VALUES ($1, $2, $3, $4)
	`
	for _, translation := range req.Items {
		_, err = tx.Exec(ctx, insertTranslationQuery, itemID, translation.LanguageCode, translation.Name, translation.Description)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert item translation for language %s: %w", op, translation.LanguageCode, err)
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

	insertItemPhotoQuery := `
		INSERT INTO ItemPhoto (item_id, photo_id)
		VALUES ($1, $2)
	`
	for _, photoID := range photoIDs {
		_, err = tx.Exec(ctx, insertItemPhotoQuery, itemID, photoID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert photo association for photo id %d: %w", op, photoID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	var response models.CreateItemResponse
	response.ID = itemID
	response.CategoryID = req.CategoryID
	response.CollectionID = req.CollectionID
	response.Size = req.Size
	response.Price = req.Price
	response.IsProducer = req.IsProducer
	response.IsPainted = req.IsPainted
	response.IsPopular = req.IsPopular
	response.IsNew = req.IsNew
	response.Items = req.Items

	for _, photo := range req.Photos {
		response.Photos = append(response.Photos, models.PhotosResponse{
			URL:       photo.URL,
			IsMain:    photo.IsMain,
			HashColor: photo.HashColor,
		})
	}

	return &response, nil
}

func (db *DB) UpdateItem(ctx context.Context, itemID int, req models.CreateItem) error {
	const op = "postgres.UpdateItem"

	var exists bool
	checkItemQuery := `SELECT EXISTS(SELECT 1 FROM Item WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkItemQuery, itemID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check item existence: %w", op, err)
	}
	if !exists {
		return storage.ErrItemNotFound
	}

	if len(req.Items) != 3 {
		return storage.ErrRequiredLanguage
	}

	languageCodes := map[string]bool{"ru": false, "kgz": false, "en": false}
	for _, translation := range req.Items {
		if _, ok := languageCodes[translation.LanguageCode]; !ok {
			return storage.ErrInvalidLanguageCode
		}
		languageCodes[translation.LanguageCode] = true
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	updateItemQuery := `
		UPDATE Item
		SET category_id = $1, collection_id = $2, size = $3, price = $4, 
		    isProducer = $5, isPainted = $6, isPopular = $7, isNew = $8
		WHERE id = $9
	`
	_, err = tx.Exec(ctx, updateItemQuery,
		req.CategoryID,
		req.CollectionID,
		req.Size,
		req.Price,
		req.IsProducer,
		req.IsPainted,
		req.IsPopular,
		req.IsNew,
		itemID,
	)
	if err != nil {
		return fmt.Errorf("%s: failed to update item: %w", op, err)
	}

	updateTranslationQuery := `
		UPDATE ItemTranslation
		SET name = $1, description = $2
		WHERE item_id = $3 AND language_code = $4
	`
	insertTranslationQuery := `
		INSERT INTO ItemTranslation (item_id, language_code, name, description)
		VALUES ($1, $2, $3, $4)
	`
	for _, translation := range req.Items {
		result, err := tx.Exec(ctx, updateTranslationQuery, translation.Name, translation.Description, itemID, translation.LanguageCode)
		if err != nil {
			return fmt.Errorf("%s: failed to update item translation for language %s: %w", op, translation.LanguageCode, err)
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(ctx, insertTranslationQuery, itemID, translation.LanguageCode, translation.Name, translation.Description)
			if err != nil {
				return fmt.Errorf("%s: failed to insert item translation for language %s: %w", op, translation.LanguageCode, err)
			}
		}
	}

	deleteItemPhotoQuery := `DELETE FROM ItemPhoto WHERE item_id = $1`
	_, err = tx.Exec(ctx, deleteItemPhotoQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete old item photos: %w", op, err)
	}

	insertPhotoQuery := `
		INSERT INTO Photo (url, isMain, hash_color)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	insertItemPhotoQuery := `
		INSERT INTO ItemPhoto (item_id, photo_id)
		VALUES ($1, $2)
	`
	for _, photo := range req.Photos {
		var photoID int
		err = tx.QueryRow(ctx, insertPhotoQuery, photo.URL, photo.IsMain, photo.HashColor).Scan(&photoID)
		if err != nil {
			return fmt.Errorf("%s: failed to insert photo with url %s: %w", op, photo.URL, err)
		}

		_, err = tx.Exec(ctx, insertItemPhotoQuery, itemID, photoID)
		if err != nil {
			return fmt.Errorf("%s: failed to insert photo association for photo id %d: %w", op, photoID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) RemoveItem(ctx context.Context, itemID int) error {
	const op = "postgres.RemoveItem"

	var exists bool
	checkItemQuery := `SELECT EXISTS(SELECT 1 FROM Item WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkItemQuery, itemID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check item existence: %w", op, err)
	}
	if !exists {
		return storage.ErrItemNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	deleteItemTranslationQuery := `DELETE FROM ItemTranslation WHERE item_id = $1`
	_, err = tx.Exec(ctx, deleteItemTranslationQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item translations: %w", op, err)
	}

	deleteItemPhotoQuery := `DELETE FROM ItemPhoto WHERE item_id = $1`
	_, err = tx.Exec(ctx, deleteItemPhotoQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item photos: %w", op, err)
	}

	deleteItemQuery := `DELETE FROM Item WHERE id = $1`
	_, err = tx.Exec(ctx, deleteItemQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) GetItemsWithoutDiscount(ctx context.Context) ([]models.ItemWithoutDiscount, error) {
	const op = "postgres.GetItemsWithoutDiscount"

	query := `
		SELECT item_id, name
		FROM ItemTranslation
		WHERE language_code = 'ru'
		AND item_id IN (
			SELECT i.id
			FROM Item i
			WHERE NOT EXISTS (
				SELECT 1 
				FROM Discount d 
				WHERE d.discount_type = 'item' AND d.target_id = i.id
			)
		)
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query items without discount: %w", op, err)
	}
	defer rows.Close()

	var items []models.ItemWithoutDiscount
	var item models.ItemWithoutDiscount
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		item.ID = id
		item.Name = name
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	return items, nil
}
