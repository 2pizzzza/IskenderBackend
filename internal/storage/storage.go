package storage

import "errors"

var (
	ErrLanguageExists     = errors.New("language already exists")
	ErrLanguageNotFound   = errors.New("language not found")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrItemNotFound       = errors.New("item not found")
)
