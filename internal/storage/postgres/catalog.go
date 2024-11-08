package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/storage"
	"github.com/jackc/pgx/v5"
)

func (db *DB) CreateCatalog(ctx context.Context, name, description, languageCode string, price float64, colorsReq []models.Color) (*schemas.CreateCatalogResponse, error) {
	const op = "postgres.CreateCatalog"

	var existingID int
	checkQuery := `
		SELECT c.id FROM Catalogs c
		JOIN Catalogs_Localization cl ON c.id = cl.catalog_id
		WHERE cl.name = $1 AND cl.language_code = $2`
	err := db.Pool.QueryRow(ctx, checkQuery, name, languageCode).Scan(&existingID)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: failed to check existing catalog: %w", op, err)
	}

	if err == nil {
		return nil, storage.ErrCatalogExists
	}

	createQuery := `INSERT INTO Catalogs (price) VALUES ($1) RETURNING id`
	var catalogID int
	err = db.Pool.QueryRow(ctx, createQuery, price).Scan(&catalogID)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to create catalog: %w", op, err)
	}

	_, err = db.InsertCatalogLocalization(ctx, catalogID, languageCode, name, description)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create catalog localization: %w", op, err)
	}

	var colors []models.Color
	for _, colorReq := range colorsReq {
		var color models.Color
		colorCheckQuery := `SELECT id, name, hash_color FROM Color WHERE name = $1 AND hash_color = $2`
		err := db.Pool.QueryRow(ctx, colorCheckQuery, colorReq.Name, colorReq.HashColor).Scan(
			&color.ID, &color.Name, &color.HashColor,
		)

		if errors.Is(err, pgx.ErrNoRows) {
			colorInsertQuery := `
				INSERT INTO Color (name, hash_color)
				VALUES ($1, $2)
				RETURNING id, name, hash_color`
			err = db.Pool.QueryRow(ctx, colorInsertQuery, colorReq.Name, colorReq.HashColor).Scan(
				&color.ID, &color.Name, &color.HashColor,
			)

			if err != nil {
				return nil, fmt.Errorf("%s: failed to create color: %w", op, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("%s: failed to check color: %w", op, err)
		}

		colors = append(colors, color)

		catalogColorInsertQuery := `
			INSERT INTO Catalog_Color (catalog_id, color_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING`
		_, err = db.Pool.Exec(ctx, catalogColorInsertQuery, catalogID, color.ID)

		if err != nil {
			return nil, fmt.Errorf("%s: failed to link color with catalog: %w", op, err)
		}
	}

	response := &schemas.CreateCatalogResponse{
		Id:          catalogID,
		Name:        name,
		Description: description,
		Price:       price,
		Color:       colors,
	}

	return response, nil
}

func (db *DB) InsertCatalogLocalization(ctx context.Context, catalogID int, languageCode, name, description string) (*schemas.CatalogLocalization, error) {
	const op = "postgres.InsertLocalization"

	checkLocalizationQuery := `
		SELECT id, catalog_id, language_code, name, description 
		FROM Catalogs_Localization 
		WHERE catalog_id = $1 AND language_code = $2 AND name = $3`
	var existingLocalization schemas.CatalogLocalization
	err := db.Pool.QueryRow(ctx, checkLocalizationQuery, catalogID, languageCode, name).Scan(
		&existingLocalization.ID,
		&existingLocalization.CatalogID,
		&existingLocalization.LanguageCode,
		&existingLocalization.Name,
		&existingLocalization.Description,
	)

	if err == nil {
		return &existingLocalization, storage.ErrCatalogExists
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: failed to check existing localization: %w", op, err)
	}

	insertLocalizationQuery := `
		INSERT INTO Catalogs_Localization (catalog_id, language_code, name, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, catalog_id, language_code, name, description`
	err = db.Pool.QueryRow(ctx, insertLocalizationQuery, catalogID, languageCode, name, description).Scan(
		&existingLocalization.ID,
		&existingLocalization.CatalogID,
		&existingLocalization.LanguageCode,
		&existingLocalization.Name,
		&existingLocalization.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert catalog localization: %w", op, err)
	}

	return &existingLocalization, nil
}

func (db *DB) GetCatalogsByLanguage(ctx context.Context, languageCode string) ([]*schemas.CatalogResponse, error) {
	const op = "postgres.GetCatalogsByLanguage"

	query := `
		SELECT c.id, cl.name, cl.description, c.price
		FROM Catalogs c
		JOIN Catalogs_Localization cl ON c.id = cl.catalog_id
		WHERE cl.language_code = $1`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get catalogs by language: %w", op, err)
	}
	defer rows.Close()

	var catalogs []*schemas.CatalogResponse
	for rows.Next() {
		var catalog schemas.CatalogResponse
		if err := rows.Scan(&catalog.ID, &catalog.Name, &catalog.Description, &catalog.Price); err != nil {
			return nil, fmt.Errorf("%s: failed to scan catalog row: %w", op, err)
		}
		catalogs = append(catalogs, &catalog)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return catalogs, nil
}

func (db *DB) DeleteCatalog(ctx context.Context, catalogID int) error {
	const op = "postgres.DeleteCatalog"

	checkCatalogQuery := `SELECT id FROM Catalogs WHERE id = $1`
	var existingID int
	err := db.Pool.QueryRow(ctx, checkCatalogQuery, catalogID).Scan(&existingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrCatalogNotFound
		}
		return fmt.Errorf("%s: failed to check catalog existence: %w", op, err)
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	deleteLocalizationsQuery := `DELETE FROM Catalogs_Localization WHERE catalog_id = $1`
	_, err = tx.Exec(ctx, deleteLocalizationsQuery, catalogID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete catalog localizations: %w", op, err)
	}

	deleteCatalogQuery := `DELETE FROM Catalogs WHERE id = $1`
	_, err = tx.Exec(ctx, deleteCatalogQuery, catalogID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete catalog: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) UpdateCatalog(ctx context.Context, catalogID int, languageCode, newName, newDescription string, newPrice float64) error {
	const op = "postgres.UpdateCatalog"

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	checkCatalogQuery := `SELECT id FROM Catalogs WHERE id = $1`
	var existingID int
	err = tx.QueryRow(ctx, checkCatalogQuery, catalogID).Scan(&existingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrCatalogNotFound
		}
		return fmt.Errorf("%s: failed to check catalog existence: %w", op, err)
	}

	updateCatalogQuery := `
		UPDATE Catalogs
		SET price = COALESCE($1, price)  
		WHERE id = $2`
	_, err = tx.Exec(ctx, updateCatalogQuery, newPrice, catalogID)
	if err != nil {
		return fmt.Errorf("%s: failed to update catalog price: %w", op, err)
	}

	checkLocalizationQuery := `SELECT id FROM Catalogs_Localization WHERE catalog_id = $1 AND language_code = $2`
	var localizationID int
	err = tx.QueryRow(ctx, checkLocalizationQuery, catalogID, languageCode).Scan(&localizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrCatalogNotFound
		}
		return fmt.Errorf("%s: failed to check localization existence: %w", op, err)
	}

	updateLocalizationQuery := `
		UPDATE Catalogs_Localization
		SET name = COALESCE($1, name), description = COALESCE($2, description)  -- если не передано, сохраняем старое
		WHERE catalog_id = $3 AND language_code = $4`
	_, err = tx.Exec(ctx, updateLocalizationQuery, newName, newDescription, catalogID, languageCode)
	if err != nil {
		return fmt.Errorf("%s: failed to update catalog localization: %w", op, err)
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) GetCatalogByID(ctx context.Context, catalogID int) (*schemas.CatalogDetailResponse, error) {
	const op = "postgres.GetCatalogByID"

	checkCatalogQuery := `SELECT id FROM Catalogs WHERE id = $1`
	var existingID int
	err := db.Pool.QueryRow(ctx, checkCatalogQuery, catalogID).Scan(&existingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrCatalogNotFound
		}
		return nil, fmt.Errorf("%s: failed to check catalog existence: %w", op, err)
	}

	query := `
		SELECT c.id, c.price, cl.language_code, cl.name, cl.description, col.id AS color_id, col.name AS color_name, col.hash_color
		FROM Catalogs c
		JOIN Catalogs_Localization cl ON c.id = cl.catalog_id
		LEFT JOIN Catalog_Color cc ON c.id = cc.catalog_id
		LEFT JOIN Color col ON cc.color_id = col.id
		WHERE c.id = $1`

	rows, err := db.Pool.Query(ctx, query, catalogID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get catalog by ID: %w", op, err)
	}
	defer rows.Close()

	var catalog schemas.CatalogDetailResponse
	catalog.Languages = make([]schemas.CatalogLocalizationResponse, 0)

	//colorsByLanguage := make(map[string][]schemas.ColorResponse)

	for rows.Next() {
		var langCode, name, description, colorName, hashColor string
		var price float64
		var colorID int
		err := rows.Scan(&catalog.ID, &price, &langCode, &name, &description, &colorID, &colorName, &hashColor)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		catalog.Price = price

		var langResponse *schemas.CatalogLocalizationResponse
		for i := range catalog.Languages {
			if catalog.Languages[i].LanguageCode == langCode {
				langResponse = &catalog.Languages[i]
				break
			}
		}

		if langResponse == nil {
			langResponse = &schemas.CatalogLocalizationResponse{
				LanguageCode: langCode,
				Name:         name,
				Description:  description,
			}
			catalog.Languages = append(catalog.Languages, *langResponse)
		}

		if colorID != 0 {
			langResponse.Colors = append(langResponse.Colors, schemas.ColorResponse{
				ID:        colorID,
				Name:      colorName,
				HashColor: hashColor,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return &catalog, nil
}
