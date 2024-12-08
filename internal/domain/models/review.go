package models

import "time"

type CreateReviewRequest struct {
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Text     string `json:"text"`
}

type ReviewResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type RemoveReview struct {
	ID int `json:"id"`
}

type ReviewResponseAdmin struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text"`
	IsShow    bool      `json:"is_show"`
	CreatedAt time.Time `json:"created_at"`
}
