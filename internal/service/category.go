package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

func (rp *Plumping) CreateCategory(ctx context.Context, req *schemas.CreateCategoryRequest) (res *schemas.CreateCategoryResponse, err error) {
	const op = "service.CreateItem"

	log := rp.log.With(
		slog.String("op", op),
	)

	item, err := rp.plumpingRepository.SaveCategory(ctx, req.Name)

	log.Debug("value: ", req.Name)
	if err != nil {
		log.Error("Failed Create Song", sl.Err(err))

		return &schemas.CreateCategoryResponse{}, fmt.Errorf("%s, %w", op, err)
	}
	return &schemas.CreateCategoryResponse{
		CategoryID: item.CategoryID,
		Name:       item.Name,
	}, nil
}

func (rp *Plumping) GetAllCategory(ctx context.Context) (res *schemas.CategoriesResponse, err error) {
	const op = "service.GetAllCategory"

	log := rp.log.With(
		slog.String("op: ", op),
	)

	categoriesRaw, err := rp.plumpingRepository.GetAllCategories(ctx)

	if err != nil {
		log.Error("Failed get all categories: ", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	categories := make([]models.Category, len(*categoriesRaw))
	for i, v := range *categoriesRaw {
		categories[i] = models.Category{
			CategoryID: v.CategoryID,
			Name:       v.Name,
		}
	}
	res = &schemas.CategoriesResponse{
		Categories: categories,
	}

	return res, nil
}

func (pr *Plumping) GetCategoryById(ctx context.Context, req *schemas.CategoryByIdRequest) (res *schemas.CategoryResponse, err error) {
	const op = "service.GetCategoryById"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	categoryRaw, err := pr.plumpingRepository.GetCategoryByID(ctx, req.Id)

	if err != nil {
		if errors.Is(err, schemas.ErrItemNotFound) {
			return nil, schemas.ErrItemNotFound
		}
		log.Error("Failed to get category by id", sl.Err(err))
	}

	category := &schemas.CategoryResponse{
		CategoryID: categoryRaw.CategoryID,
		Name:       categoryRaw.Name,
	}

	return category, nil
}
