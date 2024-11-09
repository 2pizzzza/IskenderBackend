package service

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
)

type Plumping struct {
	log                *slog.Logger
	baseDir            string
	plumpingRepository PlumpingRepository
}

type PlumpingRepository interface {
	//Language
	GetLanguages(ctx context.Context) ([]*models.Language, error)

	//Category
	GetCategoriesByLanguageCode(ctx context.Context, languageCode string) ([]*models.Category, error)

	//Collection
	GetCollectionsByLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionID int, languageCode string) (*models.CollectionResponse, error)

	//Item
	GetItemsByCategoryID(ctx context.Context, categoryID int, languageCode string) ([]*models.ItemResponse, error)
	GetItemByID(ctx context.Context, itemID int, languageCode string) (*models.ItemResponse, error)
	GetItemsByCollectionID(ctx context.Context, collectionID int, languageCode string) ([]*models.ItemResponse, error)
}

func New(log *slog.Logger, baseDir string, repository PlumpingRepository) *Plumping {
	return &Plumping{
		log:                log,
		baseDir:            baseDir,
		plumpingRepository: repository,
	}
}
