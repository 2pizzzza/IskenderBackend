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

func (db *DB) CreateCatalog(ctx context.Context, name, description string, languageID int, price float64, colorsReq []models.Color) (*schemas.CreateCatalogResponse, error) {
	const op = "postgres.CreateCatalog"

	var langExists bool
	languageCheckQuery := `SELECT EXISTS (SELECT 1 FROM Language WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, languageCheckQuery, languageID).Scan(&langExists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check language existence: %w", op, err)
	}
	if !langExists {
		return nil, fmt.Errorf("%s: language with ID %d does not exist", op, languageID)
	}

	var existingID int
	checkQuery := `
		SELECT c.id FROM Catalogs c
		JOIN Catalogs_Localization cl ON c.id = cl.catalog_id
		WHERE cl.name = $1 AND cl.language_id = $2`
	err = db.Pool.QueryRow(ctx, checkQuery, name, languageID).Scan(&existingID)

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

	_, err = db.InsertCatalogLocalization(ctx, catalogID, languageID, name, description)
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
		ID:          catalogID,
		Name:        name,
		Description: description,
		Price:       price,
		Color:       colors,
	}

	return response, nil
}

func (db *DB) InsertCatalogLocalization(ctx context.Context, catalogID int, languageID int, name, description string) (*schemas.CatalogLocalization, error) {
	const op = "postgres.InsertCatalogLocalization"

	var langExists bool
	languageCheckQuery := `SELECT EXISTS (SELECT 1 FROM Language WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, languageCheckQuery, languageID).Scan(&langExists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check language existence: %w", op, err)
	}
	if !langExists {
		return nil, fmt.Errorf("%s: language with ID %d does not exist", op, languageID)
	}

	checkLocalizationQuery := `
		SELECT id, catalog_id, language_id, name, description 
		FROM Catalogs_Localization 
		WHERE catalog_id = $1 AND language_id = $2 AND name = $3`
	var existingLocalization schemas.CatalogLocalization
	err = db.Pool.QueryRow(ctx, checkLocalizationQuery, catalogID, languageID, name).Scan(
		&existingLocalization.ID,
		&existingLocalization.CatalogID,
		&existingLocalization.LanguageID,
		&existingLocalization.Name,
		&existingLocalization.Description,
	)

	if err == nil {
		return &existingLocalization, storage.ErrCatalogExists
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: failed to check existing localization: %w", op, err)
	}

	insertLocalizationQuery := `
		INSERT INTO Catalogs_Localization (catalog_id, language_id, name, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, catalog_id, language_id, name, description`
	err = db.Pool.QueryRow(ctx, insertLocalizationQuery, catalogID, languageID, name, description).Scan(
		&existingLocalization.ID,
		&existingLocalization.CatalogID,
		&existingLocalization.LanguageID,
		&existingLocalization.Name,
		&existingLocalization.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert catalog localization: %w", op, err)
	}

	return &existingLocalization, nil
}

func (db *DB) GetCatalogsByLanguageCode(ctx context.Context, languageCode string) ([]*schemas.CatalogResponse, error) {
	const op = "postgres.GetCatalogsByLanguage"

	var languageID int
	languageQuery := `SELECT id FROM Language WHERE code = $1`
	err := db.Pool.QueryRow(ctx, languageQuery, languageCode).Scan(&languageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: language with code %s not found", op, languageCode)
		}
		return nil, fmt.Errorf("%s: failed to get language id: %w", op, err)
	}

	query := `
		SELECT c.id, cl.name, cl.description, c.price
		FROM Catalogs c
		JOIN Catalogs_Localization cl ON c.id = cl.catalog_id
		WHERE cl.language_id = $1`

	rows, err := db.Pool.Query(ctx, query, languageID)
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

func (db *DB) UpdateCatalog(ctx context.Context, catalogID int, languageID int, newName, newDescription string, newPrice float64) error {
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

	checkLocalizationQuery := `SELECT id FROM Catalogs_Localization WHERE catalog_id = $1 AND language_id = $2`
	var localizationID int
	err = tx.QueryRow(ctx, checkLocalizationQuery, catalogID, languageID).Scan(&localizationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrCatalogNotFound
		}
		return fmt.Errorf("%s: failed to check localization existence: %w", op, err)
	}

	updateLocalizationQuery := `
		UPDATE Catalogs_Localization
		SET name = COALESCE($1, name), description = COALESCE($2, description)
		WHERE catalog_id = $3 AND language_id = $4`
	_, err = tx.Exec(ctx, updateLocalizationQuery, newName, newDescription, catalogID, languageID)
	if err != nil {
		return fmt.Errorf("%s: failed to update catalog localization: %w", op, err)
	}

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
		SELECT c.id, c.price, l.id AS language_id, l.code AS language_code, cl.name, cl.description, col.id AS color_id, col.name AS color_name, col.hash_color
		FROM Catalogs c
		JOIN Catalogs_Localization cl ON c.id = cl.catalog_id
		JOIN Language l ON cl.language_id = l.id
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

	for rows.Next() {
		var languageID, colorID int
		var languageCode, name, description, colorName, hashColor string
		var price float64

		err := rows.Scan(&catalog.ID, &price, &languageID, &languageCode, &name, &description, &colorID, &colorName, &hashColor)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		catalog.Price = price

		var langResponse *schemas.CatalogLocalizationResponse
		for i := range catalog.Languages {
			if catalog.Languages[i].LanguageCode == languageCode {
				langResponse = &catalog.Languages[i]
				break
			}
		}

		if langResponse == nil {
			langResponse = &schemas.CatalogLocalizationResponse{
				LanguageCode: languageCode,
				Name:         name,
				Description:  description,
				Colors:       make([]schemas.ColorResponse, 0),
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
