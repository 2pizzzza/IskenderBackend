package models

import "time"

type Discount struct {
	ID                 int              `json:"id"`
	LanguageCode       string           `json:"language_code"`
	Name               string           `json:"name"`
	CollectionID       *int             `json:"collection_id,omitempty"`
	Description        string           `json:"description"`
	NewPrice           float64          `json:"new_price"`
	OldPrice           float64          `json:"old_price"`
	DiscountID         int              `json:"discount_id"`
	DiscountType       string           `json:"discount_type"`
	TargetID           int              `json:"target_id"`
	DiscountPercentage float64          `json:"discount_percentage"`
	StartDate          time.Time        `json:"start_date"`
	EndDate            time.Time        `json:"end_date"`
	IsProducer         bool             `json:"is_producer,omitempty"`
	IsPainted          bool             `json:"is_painted,omitempty"`
	IsPopular          bool             `json:"is_popular,omitempty"`
	IsNew              bool             `json:"is_new,omitempty"`
	Photo              []PhotosResponse `json:"photo"`
	Color              []ColorResponse  `json:"color"`
}

type DiscountCreate struct {
	ID                 int       `json:"id"`
	DiscountType       string    `json:"discount_type"`
	TargetID           int       `json:"target_id"`
	DiscountPercentage float64   `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}

type DiscountRequest struct {
	Id int `json:"id"`
}
