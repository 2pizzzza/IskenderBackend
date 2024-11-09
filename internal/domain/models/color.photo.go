package models

type Photo struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	IsMain bool   `json:"isMain"`
}

type Color struct {
	ID        int    `json:"id"`
	HashColor string `json:"hash_color"`
}

type PhotosResponse struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	IsMain bool   `json:"isMain"`
}

type ColorResponse struct {
	ID        int    `json:"id"`
	HashColor string `json:"hash_color"`
}
