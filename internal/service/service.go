package service

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"log/slog"
)

type Plumping struct {
	log                *slog.Logger
	baseDir            string
	plumpingRepository PlumpingRepository
}

type PlumpingRepository interface {
	//Catalog
	CreateCatalog(ctx context.Context, name, description string, languageID int, price float64, colorsReq []models.Color) (*schemas.CreateCatalogResponse, error)
	InsertCatalogLocalization(ctx context.Context, catalogID int, languageID int, name, description string) (*schemas.CatalogLocalization, error)
	GetCatalogsByLanguageCode(ctx context.Context, languageCode string) ([]*schemas.CatalogResponse, error)
	DeleteCatalog(ctx context.Context, catalogID int) error
	UpdateCatalog(ctx context.Context, catalogID int, languageID int, newName, newDescription string, newPrice float64) error
	GetCatalogByID(ctx context.Context, catalogID int) (*schemas.CatalogDetailResponse, error)
}

func New(log *slog.Logger, baseDir string, repository PlumpingRepository) *Plumping {
	return &Plumping{
		log:                log,
		baseDir:            baseDir,
		plumpingRepository: repository,
	}
}
