package models

type Catalog struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Category struct {
	ID          int    `json:"id"`
	CatalogID   int    `json:"catalog_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Collection struct {
	ID          int     `json:"id"`
	CatalogID   int     `json:"catalog_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsPainted   bool    `json:"is_painted"`
	IsProducer  bool    `json:"is_producer"`
	Price       float64 `json:"price"`
}

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	IsProducer  bool    `json:"is_producer"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Size        string  `json:"size"`
}

type Color struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	HashColor string `json:"hash_color"`
}

type ItemColor struct {
	ItemID  int `json:"item_id"`
	ColorID int `json:"color_id"`
}

type CatalogColor struct {
	CatalogID int `json:"catalog_id"`
	ColorID   int `json:"color_id"`
}

type CollectionColor struct {
	CollectionID int `json:"collection_id"`
	ColorID      int `json:"color_id"`
}

type Photo struct {
	ID     int    `json:"id"`
	ItemID int    `json:"item_id"`
	URL    string `json:"url"`
	IsMain bool   `json:"is_main"`
}

type CategoryItem struct {
	CategoryID int `json:"category_id"`
	ItemID     int `json:"item_id"`
}

type CollectionItem struct {
	CollectionID int `json:"collection_id"`
	ItemID       int `json:"item_id"`
}
