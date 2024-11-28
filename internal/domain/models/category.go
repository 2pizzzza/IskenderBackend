package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type CategoryTranslation struct {
	CategoryID   int    `json:"category_id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
}

type UpdateCategoryRequest struct {
	CategoryID   int    `json:"category_id"`
	Name         string `json:"name"`
	LanguageCode string `json:"language_code"`
}

type RemoveCategoryRequest struct {
	ID int `json:"id"`
}
type CreateCategoryRequest struct {
	Categories []CategoriesRequest `json:"categories"`
}

type GetCategoryRequest struct {
	ID         int                 `json:"id"`
	Categories []CategoriesRequest `json:"categories"`
}

type CategoriesRequest struct {
	Name         string `json:"name"`
	LanguageCode string `json:"language_code"`
}

type CreateCategoryResponse struct {
	Categories []CategoriesResponse `json:"categories"`
}

type CategoriesResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LanguageCode string `json:"language_code"`
}

type Message struct {
	Message string `json:"message"`
}

type UpdateCategoriesResponse struct {
	Name         string `json:"name"`
	LanguageCode string `json:"language_code"`
}
