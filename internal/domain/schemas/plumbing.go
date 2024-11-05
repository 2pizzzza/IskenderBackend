package schemas

type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	Price       float64 `json:"price"`
	IsProduced  bool    `json:"is_produced"`
}

type CreateItemResponse struct {
	ItemID      int      `json:"item_id"`
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
