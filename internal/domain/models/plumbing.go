package models

type Category struct {
	CategoryID int    `json:"id"`
	Name       string `json:"name"`
}

type Item struct {
	ItemID      int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CategoryID  int      `json:"category_id"`
	Price       float64  `json:"price"`
	IsProduced  bool     `json:"is_produced"`
	Colors      []string `json:"colors"`
	Photos      []string `json:"photos"`
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
