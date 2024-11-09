package schemas

import "github.com/2pizzzza/plumbing/internal/domain/models"

type CreateCatalogRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Color       []models.Color `json:"color"`
	LanguageID  int            `json:"language_id"`
}

type CreateCatalogResponse struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Color       []models.Color `json:"color"`
}

type CatalogLocalizationRequest struct {
	CatalogID   int    `json:"catalog_id"`
	LanguageID  int    `json:"language_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CatalogLocalization struct {
	ID          int    `json:"id"`
	CatalogID   int    `json:"catalog_id"`
	LanguageID  int    `json:"language_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CatalogResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CatalogsByLanguageRequest struct {
	LanguageCode string `json:"language_code"`
}

type CatalogRemoveRequest struct {
	ID int `json:"id"`
}

type UpdateCatalogRequest struct {
	ID             int     `json:"id"`
	NewName        string  `json:"new_name"`
	NewDescription string  `json:"new_description"`
	NewPrice       float64 `json:"new_price"`
	LanguageID     int     `json:"language_id"`
}

type CatalogDetailResponse struct {
	ID        int                           `json:"id"`
	Price     float64                       `json:"price"`
	Languages []CatalogLocalizationResponse `json:"languages"`
}

type CatalogLocalizationResponse struct {
	LanguageCode string          `json:"language_code"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Colors       []ColorResponse `json:"colors"`
}

type ColorResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	HashColor string `json:"hash_color"`
}
