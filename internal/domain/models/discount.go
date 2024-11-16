package models

import "time"

type Discount struct {
	ID                 int       `json:"id"`
	DiscountType       string    `json:"discount_type"`
	TargetID           int       `json:"target_id"`
	DiscountPercentage float64   `json:"discount_percentage"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}
