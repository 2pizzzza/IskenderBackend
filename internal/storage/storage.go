package storage

import "errors"

var (
	ErrLanguageExists   = errors.New("language already exists")
	ErrLanguageNotFound = errors.New("language not found")
)
