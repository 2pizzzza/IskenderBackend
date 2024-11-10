package plumbing

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

	//Starter
	CreateStarter(ctx context.Context) error

	//Language
	GetLanguages(ctx context.Context) ([]*models.Language, error)

	//Category
	GetCategoriesByLanguageCode(ctx context.Context, languageCode string) ([]*models.Category, error)
	UpdateCategory(ctx context.Context, categoryID int, name string, languageCode string) error
	CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.CreateCategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryID int) error

	//Collection
	GetCollectionsByLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionID int, languageCode string) (*models.CollectionResponse, error)

	//Popular and new
	GetPopularCollections(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetPopularItems(ctx context.Context, languageCode string) ([]*models.ItemResponse, error)
	GetNewCollections(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetNewItems(ctx context.Context, languageCode string) ([]*models.ItemResponse, error)

	//Search and filtr
	SearchCollections(ctx context.Context, languageCode string, isProducer *bool, searchQuery string) ([]*models.CollectionResponse, error)
	SearchItems(ctx context.Context, languageCode string, isProducer *bool, searchQuery string) ([]*models.ItemResponse, error)

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
