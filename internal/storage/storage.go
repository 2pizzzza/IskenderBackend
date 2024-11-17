package storage

import "errors"

var (
	ErrLanguageExists             = errors.New("language already exists")
	ErrBrandExists                = errors.New("brand already exists")
	ErrBrandNotFound              = errors.New("brand not found")
	ErrLanguageNotFound           = errors.New("language not found")
	ErrCategoryNotFound           = errors.New("category not found")
	ErrCollectionNotFound         = errors.New("collection not found")
	ErrItemNotFound               = errors.New("item not found")
	ErrCategoryExists             = errors.New("category already exists")
	ErrRequiredLanguage           = errors.New("exactly 3 languages are required")
	ErrToken                      = errors.New("Token not valid")
	ErrVacancyNotFound            = errors.New("Vacancy not found")
	ErrVacancyTranslationNotFound = errors.New("translation not found")
	ErrDiscountExists             = errors.New("discount already exists")
	ErrDiscountNotFound           = errors.New("discount already exists")
	ErrInvalidLanguageCode        = errors.New("invalid language code ")
	ErrItemExists                 = errors.New("item already exists")
)
