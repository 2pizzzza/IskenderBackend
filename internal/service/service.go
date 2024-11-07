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
	SaveItem(ctx context.Context, name, description string, categoryId int, price float64) (models.Item, error)
	GetItemById(ctx context.Context, id int) (models.Item, error)
	SaveItemWithDetails(ctx context.Context, name, description string, categoryId int, price float64, colors, photos []string) (models.Item, error)
	UpdateItem(ctx context.Context, itemID int, name, description string, categoryId int, price float64, isProduced bool, colors, photos []string) (models.Item, error)
	RemoveItem(ctx context.Context, itemID int) error
	GetAllItemsByCategory(ctx context.Context, categoryID int) ([]models.Item, error)
	GetAllItems(ctx context.Context) ([]models.Item, error)
	// Category
	SaveCategory(ctx context.Context, name string) (models.Category, error)
	GetAllCategories(ctx context.Context) (*[]models.Category, error)
	GetCategoryByID(ctx context.Context, categoryID int) (models.Category, error)
	UpdateCategory(ctx context.Context, categoryID int, name string) error
	RemoveCategory(ctx context.Context, categoryID int) error
}

func New(log *slog.Logger, baseDir string, repository PlumpingRepository) *Plumping {
	return &Plumping{
		log:                log,
		baseDir:            baseDir,
		plumpingRepository: repository,
	}
}
