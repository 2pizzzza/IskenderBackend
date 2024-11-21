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

type ItemRequest struct {
	ID int `json:"id"`
}

type ItemColor struct {
	ItemID  int `json:"item_id"`
	ColorID int `json:"color_id"`
}
type ItemResponse struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	CategoryID   int              `json:"category_id"`
	CollectionID int              `json:"collection_id"`
	Size         string           `json:"size"`
	Price        float64          `json:"price"`
	NewPrice     float64          `json:"new_price"`
	IsProducer   bool             `json:"isProducer"`
	IsPainted    bool             `json:"isPainted"`
	IsPopular    bool             `json:"is_popular"`
	IsNew        bool             `json:"is_new"`
	Photos       []PhotosResponse `json:"photos"`
	Colors       []ColorResponse  `json:"colors"`
}

type CreateItem struct {
	CategoryID   int                     `json:"category_id"`
	CollectionID int                     `json:"collection_id"`
	Size         string                  `json:"size"`
	Price        float64                 `json:"price"`
	IsProducer   bool                    `json:"isProducer"`
	IsPainted    bool                    `json:"isPainted"`
	IsPopular    bool                    `json:"is_popular"`
	IsNew        bool                    `json:"is_new"`
	Photos       []PhotosResponse        `json:"photos"`
	Items        []CreateItemTranslation `json:"items"`
}

type CreateItemTranslation struct {
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type CreateItemResponse struct {
	ID           int                     `json:"ID"`
	CategoryID   int                     `json:"category_id"`
	CollectionID int                     `json:"collection_id"`
	Size         string                  `json:"size"`
	Price        float64                 `json:"price"`
	IsProducer   bool                    `json:"isProducer"`
	IsPainted    bool                    `json:"isPainted"`
	IsPopular    bool                    `json:"is_popular"`
	IsNew        bool                    `json:"is_new"`
	Photos       []PhotosResponse        `json:"photos"`
	Items        []CreateItemTranslation `json:"items"`
}

type ItemResponses struct {
	ID           int                     `json:"ID"`
	CategoryID   int                     `json:"category_id"`
	CollectionID int                     `json:"collection_id"`
	Size         string                  `json:"size"`
	Price        float64                 `json:"price"`
	IsProducer   bool                    `json:"isProducer"`
	IsPainted    bool                    `json:"isPainted"`
	IsPopular    bool                    `json:"is_popular"`
	IsNew        bool                    `json:"is_new"`
	Photos       []PhotosResponse        `json:"photos"`
	Items        []CreateItemTranslation `json:"items"`
	Color        []ColorResponse         `json:"color"`
}

type ItemResponseForAdmin struct {
	ID           int                     `json:"ID"`
	CategoryID   int                     `json:"category_id"`
	CollectionID int                     `json:"collection_id"`
	Size         string                  `json:"size"`
	Price        float64                 `json:"price"`
	IsProducer   bool                    `json:"isProducer"`
	IsPainted    bool                    `json:"isPainted"`
	IsPopular    bool                    `json:"is_popular"`
	IsNew        bool                    `json:"is_new"`
	Photos       []PhotosResponse        `json:"photos"`
	Items        []CreateItemTranslation `json:"items"`
}
