package models

type Catalog struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type Category struct {
	Id         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
}

type Collection struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Colors      []string `json:"colors"`
	CategoryID  int      `json:"category_id"`
	IsProduced  bool     `json:"is_produced"`
	IsPained    bool     `json:"is_pained"`
}

type Item struct {
	ItemID       int     `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	CategoryID   int     `json:"category_id"`
	CollectionID int     `json:"collection_id"`
	Price        float64 `json:"price"`
	IsProduced   bool    `json:"is_produced"`
	Photos       string  `json:"photos"`
}

type Color struct {
	ColorID int    `json:"id" db:"color_id"`
	Name    string `json:"name" db:"name"`
}

type ItemColor struct {
	ItemID  int `json:"id" db:"item_id"`
	ColorID int `json:"color_id" db:"color_id"`
}

type Photo struct {
	PhotoID int    `json:"id" db:"photo_id"`
	ItemID  int    `json:"item_id" db:"item_id"`
	URL     string `json:"url" db:"url"`
	IsMain  bool   `json:"is_main" db:"is_main"`
}
