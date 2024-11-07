package schemas

import (
	"errors"
	"github.com/2pizzzza/plumbing/internal/domain/models"
)

type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	IsProduced  bool    `json:"is_produced"`
}

type CreateItemResponse struct {
	ItemID      int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	Price       float64  `json:"price"`
	IsProduced  bool     `json:"is_produced"`
	Colors      []string `json:"colors"`
	Photos      []string `json:"photos"`
}

type CreateCategoryRequest struct {
	CategoryID int    `json:"id"`
	Name       string `json:"name"`
}

type CreateCategoryResponse struct {
	CategoryID int    `json:"id"`
	Name       string `json:"name"`
}

type GetItemByIdRequest struct {
	ItemID int `json:"id"`
}

var ErrItemNotFound = errors.New("item not found")

type CategoriesResponse struct {
	Categories []models.Category `json:"categories"`
}

type CategoryByIdRequest struct {
	Id int `json:"id"`
}

type CategoryResponse struct {
	CategoryID int    `json:"id"`
	Name       string `json:"name"`
}
