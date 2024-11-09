package models

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type CategoryTranslation struct {
	CategoryID   int    `json:"category_id"`
	LanguageCode string `json:"language_code"`
	Name         string `json:"name"`
}
