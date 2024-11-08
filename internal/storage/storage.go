package storage

import "errors"

var (
	ErrCatalogExists   = errors.New("catalog already exists")
	ErrCatalogNotFound = errors.New("catalog not found")
)
