package models

type BrandResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"photo"`
}

type BrandRequest struct {
	Name string `json:"name"`
	Url  string `json:"photo"`
}

type RemoveBrandRequest struct {
	ID int `json:"id"`
}
