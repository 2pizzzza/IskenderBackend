package models

type VacancyResponse struct {
	Id               int      `json:"id"`
	LanguageCode     string   `json:"language_code"`
	Title            string   `json:"title"`
	Requirements     []string `json:"requirements"`
	Responsibilities []string `json:"responsibilities"`
	Conditions       []string `json:"conditions"`
	Information      []string `json:"information"`
	IsActive         bool     `json:"isActive"`
	Salary           int      `json:"salary"`
}

type VacancyResponses struct {
	Salary   int              `json:"salary"`
	IsActive bool             `json:"is_active"`
	Vacancy  []*CreateVacancy `json:"vacancy"`
}

type VacancyUpdateRequest struct {
	Id       int              `json:"id,omitempty"`
	Salary   int              `json:"salary" json:"salary,omitempty"`
	IsActive bool             `json:"is_active" json:"isActive,omitempty"`
	Vacancy  []*CreateVacancy `json:"vacancy" json:"vacancy,omitempty"`
}

type UpdateVacancy struct {
	LanguageCode     string   `json:"language_code"`
	Title            string   `json:"title"`
	Requirements     []string `json:"requirements"`
	Responsibilities []string `json:"responsibilities"`
	Conditions       []string `json:"conditions"`
	Information      []string `json:"information"`
}

type CreateVacancy struct {
	Id               int      `json:"id,omitempty"`
	LanguageCode     string   `json:"language_code"`
	Title            string   `json:"title"`
	Requirements     []string `json:"requirements"`
	Responsibilities []string `json:"responsibilities"`
	Conditions       []string `json:"conditions"`
	Information      []string `json:"information"`
}

type RemoveVacancyRequest struct {
	ID int `json:"id"`
}
