package service

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
)

type Plumping struct {
	log                *slog.Logger
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
}

func New(log *slog.Logger, repository PlumpingRepository) *Plumping {
	return &Plumping{
		log:                log,
		plumpingRepository: repository,
	}
}
