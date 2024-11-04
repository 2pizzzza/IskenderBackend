package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/jackc/pgx/v5"
)

func (db *DB) SaveItem(
	ctx context.Context, name, description string, categoryId int, price float64) (models.Item, error) {

	const op = "postgres.SaveItem"

	var exist bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM Item 
			WHERE name = $1 AND category_id = $2
		)
	`

	err := db.Pool.QueryRow(ctx, checkQuery, name, categoryId).Scan(&exist)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to check if item exists: %w", op, err)
	}

	if exist {
		return models.Item{}, fmt.Errorf("%s: item already exists", op)
	}

	insertQuery := `
		INSERT INTO Item (name, description, category_id, price, is_produced)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING item_id, name, description, category_id, price, is_produced
	`

	var newItem models.Item
	err = db.Pool.QueryRow(ctx, insertQuery, name, description, categoryId, price, false).Scan(
		&newItem.ItemID,
		&newItem.Name,
		&newItem.Description,
		&newItem.CategoryID,
		&newItem.Price,
		&newItem.IsProduced,
	)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to insert item: %w", op, err)
	}

	return newItem, nil
}

func (db *DB) GetItemById(ctx context.Context, id int) (models.Item, error) {
	const op = "postgres.GetItemById"

	var item models.Item
	item.Colors = []string{}
	item.Photos = []string{}

	query := `
		SELECT 
			i.item_id, i.name, i.description, i.category_id, i.price, i.is_produced,
			c.name AS color,
			p.url AS photo
		FROM 
			Item i
		LEFT JOIN 
			Item_Color ic ON i.item_id = ic.item_id
		LEFT JOIN 
			Color c ON ic.color_id = c.color_id
		LEFT JOIN 
			Photo p ON i.item_id = p.item_id
		WHERE 
			i.item_id = $1
	`

	rows, err := db.Pool.Query(ctx, query, id)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to get item by id: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var color, photo sql.NullString

		err := rows.Scan(
			&item.ItemID,
			&item.Name,
			&item.Description,
			&item.CategoryID,
			&item.Price,
			&item.IsProduced,
			&color,
			&photo,
		)
		if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to scan item row: %w", op, err)
		}

		if color.Valid {
			item.Colors = append(item.Colors, color.String)
		}

		if photo.Valid {
			item.Photos = append(item.Photos, photo.String)
		}
	}

	if item.ItemID == 0 {
		return models.Item{}, fmt.Errorf("%s: item not found with id %d", op, id)
	}

	return item, nil
}

func (db *DB) SaveItemWithDetails(
	ctx context.Context, name, description string, categoryId int, price float64, colors, photos []string) (models.Item, error) {

	const op = "postgres.SaveItemWithDetails"

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	insertItemQuery := `
		INSERT INTO Item (name, description, category_id, price, is_produced)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING item_id, name, description, category_id, price, is_produced
	`
	var newItem models.Item
	err = tx.QueryRow(ctx, insertItemQuery, name, description, categoryId, price, false).Scan(
		&newItem.ItemID,
		&newItem.Name,
		&newItem.Description,
		&newItem.CategoryID,
		&newItem.Price,
		&newItem.IsProduced,
	)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to insert item: %w", op, err)
	}

	for _, colorName := range colors {
		var colorID int
		selectColorQuery := `SELECT color_id FROM Color WHERE name = $1`
		err := tx.QueryRow(ctx, selectColorQuery, colorName).Scan(&colorID)

		if err == pgx.ErrNoRows {
			insertColorQuery := `INSERT INTO Color (name) VALUES ($1) RETURNING color_id`
			err = tx.QueryRow(ctx, insertColorQuery, colorName).Scan(&colorID)
			if err != nil {
				return models.Item{}, fmt.Errorf("%s: failed to insert color: %w", op, err)
			}
		} else if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to check color existence: %w", op, err)
		}

		insertItemColorQuery := `
			INSERT INTO Item_Color (item_id, color_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`
		_, err = tx.Exec(ctx, insertItemColorQuery, newItem.ItemID, colorID)
		if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to insert item-color relation: %w", op, err)
		}

		newItem.Colors = append(newItem.Colors, colorName)
	}

	for _, photoURL := range photos {
		insertPhotoQuery := `
			INSERT INTO Photo (item_id, url) VALUES ($1, $2)
		`
		_, err := tx.Exec(ctx, insertPhotoQuery, newItem.ItemID, photoURL)
		if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to insert photo: %w", op, err)
		}
		newItem.Photos = append(newItem.Photos, photoURL)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return newItem, nil
}

func (db *DB) UpdateItem(
	ctx context.Context, itemID int, name, description string, categoryId int, price float64, isProduced bool, colors, photos []string) (models.Item, error) {

	const op = "postgres.UpdateItem"

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	updateItemQuery := `
		UPDATE Item 
		SET name = $1, description = $2, category_id = $3, price = $4, is_produced = $5
		WHERE item_id = $6
		RETURNING item_id, name, description, category_id, price, is_produced
	`
	var updatedItem models.Item
	err = tx.QueryRow(ctx, updateItemQuery, name, description, categoryId, price, isProduced, itemID).Scan(
		&updatedItem.ItemID,
		&updatedItem.Name,
		&updatedItem.Description,
		&updatedItem.CategoryID,
		&updatedItem.Price,
		&updatedItem.IsProduced,
	)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to update item: %w", op, err)
	}

	deleteColorsQuery := `DELETE FROM Item_Color WHERE item_id = $1`
	_, err = tx.Exec(ctx, deleteColorsQuery, itemID)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to delete old colors: %w", op, err)
	}

	for _, colorName := range colors {
		var colorID int
		selectColorQuery := `SELECT color_id FROM Color WHERE name = $1`
		err := tx.QueryRow(ctx, selectColorQuery, colorName).Scan(&colorID)

		if err == pgx.ErrNoRows {
			insertColorQuery := `INSERT INTO Color (name) VALUES ($1) RETURNING color_id`
			err = tx.QueryRow(ctx, insertColorQuery, colorName).Scan(&colorID)
			if err != nil {
				return models.Item{}, fmt.Errorf("%s: failed to insert color: %w", op, err)
			}
		} else if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to check color existence: %w", op, err)
		}

		insertItemColorQuery := `
			INSERT INTO Item_Color (item_id, color_id) VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`
		_, err = tx.Exec(ctx, insertItemColorQuery, itemID, colorID)
		if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to insert item-color relation: %w", op, err)
		}

		updatedItem.Colors = append(updatedItem.Colors, colorName)
	}

	deletePhotosQuery := `DELETE FROM Photo WHERE item_id = $1`
	_, err = tx.Exec(ctx, deletePhotosQuery, itemID)
	if err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to delete old photos: %w", op, err)
	}

	for _, photoURL := range photos {
		insertPhotoQuery := `INSERT INTO Photo (item_id, url) VALUES ($1, $2)`
		_, err := tx.Exec(ctx, insertPhotoQuery, itemID, photoURL)
		if err != nil {
			return models.Item{}, fmt.Errorf("%s: failed to insert photo: %w", op, err)
		}
		updatedItem.Photos = append(updatedItem.Photos, photoURL)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Item{}, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return updatedItem, nil
}

func (db *DB) RemoveItem(ctx context.Context, itemID int) error {
	const op = "postgres.RemoveItem"

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	deletePhotosQuery := `DELETE FROM Photo WHERE item_id = $1`
	_, err = tx.Exec(ctx, deletePhotosQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete photos: %w", op, err)
	}

	deleteItemColorsQuery := `DELETE FROM Item_Color WHERE item_id = $1`
	_, err = tx.Exec(ctx, deleteItemColorsQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item-color relations: %w", op, err)
	}

	deleteItemQuery := `DELETE FROM Item WHERE item_id = $1`
	_, err = tx.Exec(ctx, deleteItemQuery, itemID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete item: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) GetAllItemsByCategory(ctx context.Context, categoryID int) ([]models.Item, error) {
	const op = "postgres.GetAllItemsByCategory"

	itemsQuery := `
		SELECT item_id, name, description, category_id, price, is_produced
		FROM Item
		WHERE category_id = $1
	`
	rows, err := db.Pool.Query(ctx, itemsQuery, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query items: %w", op, err)
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ItemID,
			&item.Name,
			&item.Description,
			&item.CategoryID,
			&item.Price,
			&item.IsProduced,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan item: %w", op, err)
		}

		colorQuery := `
			SELECT c.name
			FROM Color c
			INNER JOIN Item_Color ic ON c.color_id = ic.color_id
			WHERE ic.item_id = $1
		`
		colorRows, err := db.Pool.Query(ctx, colorQuery, item.ItemID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query colors: %w", op, err)
		}
		defer colorRows.Close()

		for colorRows.Next() {
			var color string
			err = colorRows.Scan(&color)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to scan color: %w", op, err)
			}
			item.Colors = append(item.Colors, color)
		}

		photoQuery := `
			SELECT url
			FROM Photo
			WHERE item_id = $1
		`
		photoRows, err := db.Pool.Query(ctx, photoQuery, item.ItemID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query photos: %w", op, err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photo string
			err = photoRows.Scan(&photo)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to scan photo: %w", op, err)
			}
			item.Photos = append(item.Photos, photo)
		}

		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, rows.Err())
	}

	return items, nil
}

func (db *DB) GetAllItems(ctx context.Context) ([]models.Item, error) {
	const op = "postgres.GetAllItems"

	itemsQuery := `
		SELECT item_id, name, description, category_id, price, is_produced
		FROM Item
	`
	rows, err := db.Pool.Query(ctx, itemsQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query items: %w", op, err)
	}
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ItemID,
			&item.Name,
			&item.Description,
			&item.CategoryID,
			&item.Price,
			&item.IsProduced,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan item: %w", op, err)
		}

		colorQuery := `
			SELECT c.name
			FROM Color c
			INNER JOIN Item_Color ic ON c.color_id = ic.color_id
			WHERE ic.item_id = $1
		`
		colorRows, err := db.Pool.Query(ctx, colorQuery, item.ItemID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query colors: %w", op, err)
		}
		defer colorRows.Close()

		for colorRows.Next() {
			var color string
			err = colorRows.Scan(&color)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to scan color: %w", op, err)
			}
			item.Colors = append(item.Colors, color)
		}

		photoQuery := `
			SELECT url
			FROM Photo
			WHERE item_id = $1
		`
		photoRows, err := db.Pool.Query(ctx, photoQuery, item.ItemID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to query photos: %w", op, err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photo string
			err = photoRows.Scan(&photo)
			if err != nil {
				return nil, fmt.Errorf("%s: failed to scan photo: %w", op, err)
			}
			item.Photos = append(item.Photos, photo)
		}

		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, rows.Err())
	}

	return items, nil
}
