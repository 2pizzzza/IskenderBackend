package models

type Item struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"category_id"`
	CollectionID int     `json:"collection_id"`
	Size         string  `json:"size"`
	Price        float64 `json:"price"`
	IsProducer   bool    `json:"isProducer"`
	IsPainted    bool    `json:"isPainted"`
}

type ItemTranslation struct {
	ItemID       int    `json:"item_id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type ItemPhoto struct {
	ItemID  int `json:"item_id"`
	PhotoID int `json:"photo_id"`
}

type ItemColor struct {
	ItemID  int `json:"item_id"`
	ColorID int `json:"color_id"`
}
