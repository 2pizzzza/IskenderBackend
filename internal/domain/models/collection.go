package models

type Collection struct {
	ID         int     `json:"id"`
	CategoryID int     `json:"category_id"`
	Price      float64 `json:"price"`
	IsProducer bool    `json:"isProducer"`
	IsPainted  bool    `json:"isPainted"`
}

type CollectionTranslation struct {
	CollectionID int    `json:"collection_id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type CollectionPhoto struct {
	CollectionID int `json:"collection_id"`
	PhotoID      int `json:"photo_id"`
}

type CollectionColor struct {
	CollectionID int `json:"collection_id"`
	ColorID      int `json:"color_id"`
}

type CollectionResponse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Price       float64          `json:"price"`
	IsProducer  bool             `json:"isProducer"`
	IsPainted   bool             `json:"isPainted"`
	IsPopular   bool             `json:"is_popular"`
	IsNew       bool             `json:"is_new"`
	Photos      []PhotosResponse `json:"photos"`
	Colors      []ColorResponse  `json:"colors"`
}

type PopularResponse struct {
	Collections []*CollectionResponse `json:"collections"`
	Items       []*ItemResponse       `json:"items"`
}

type UpdateCollectionRequest struct {
	CollectionID int                `json:"collection_id"`
	Price        float64            `json:"price,omitempty"`
	IsProducer   bool               `json:"isProducer,omitempty"`
	IsPainted    bool               `json:"isPainted,omitempty"`
	IsPopular    bool               `json:"isPopular,omitempty"`
	IsNew        bool               `json:"isNew,omitempty"`
	Photos       []PhotosResponse   `json:"photos"`
	Colors       []ColorResponse    `json:"colors"`
	Collections  []UpdateCollection `json:"collections"`
}
type UpdateCollection struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	LanguageCode string `json:"language_code"`
}

type RemoveCollectionRequest struct {
	ID int `json:"id"`
}
