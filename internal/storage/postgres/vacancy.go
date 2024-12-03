package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/storage"
)

func (db *DB) GetAllActiveVacanciesByLanguage(ctx context.Context, languageCode string) ([]models.VacancyResponse, error) {
	const op = "postgres.GetAllActiveVacanciesByLanguage"

	query := `
		SELECT v.id, vt.language_code, vt.title, vt.requirements, 
		       vt.responsibilities, vt.conditions, vt.information, v.salary
		FROM Vacancy v
		JOIN VacancyTranslation vt ON v.id = vt.vacancy_id
		WHERE v.isActive = TRUE AND vt.language_code = $1`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query active vacancies: %w", op, err)
	}
	defer rows.Close()

	var vacancies []models.VacancyResponse

	for rows.Next() {
		var vacancy models.VacancyResponse
		if err := rows.Scan(
			&vacancy.Id,
			&vacancy.LanguageCode,
			&vacancy.Title,
			&vacancy.Requirements,
			&vacancy.Responsibilities,
			&vacancy.Conditions,
			&vacancy.Information,
			&vacancy.Salary,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row into vacancy struct: %w", op, err)
		}
		vacancies = append(vacancies, vacancy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return vacancies, nil
}

func (db *DB) UpdateVacancy(ctx context.Context, req models.VacancyUpdateRequest) error {
	const op = "postgres.UpdateVacancy"

	var exists bool
	checkVacancyQuery := `SELECT EXISTS(SELECT 1 FROM Vacancy WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkVacancyQuery, req.Id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check vacancy existence: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: vacancy with id %d not found", op, req.Id)
	}

	if len(req.Vacancy) != 3 {
		return fmt.Errorf("%s: all translations for vacancy must be provided (3 languages required)", op)
	}
	languageCodes := map[string]bool{"ru": false, "kgz": false, "en": false}
	for _, translation := range req.Vacancy {
		if _, ok := languageCodes[translation.LanguageCode]; !ok {
			return fmt.Errorf("%s: invalid language code %s", op, translation.LanguageCode)
		}
		languageCodes[translation.LanguageCode] = true
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	updateVacancyQuery := `
		UPDATE Vacancy
		SET salary = $1, isActive = $2
		WHERE id = $3
	`
	_, err = tx.Exec(ctx, updateVacancyQuery, req.Salary, req.IsActive, req.Id)
	if err != nil {
		return fmt.Errorf("%s: failed to update vacancy: %w", op, err)
	}

	updateTranslationQuery := `
		UPDATE VacancyTranslation
		SET title = $1, requirements = $2, responsibilities = $3, conditions = $4, information = $5
		WHERE vacancy_id = $6 AND language_code = $7
	`
	insertTranslationQuery := `
		INSERT INTO VacancyTranslation (vacancy_id, language_code, title, requirements, responsibilities, conditions, information)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	for _, translation := range req.Vacancy {
		result, err := tx.Exec(ctx, updateTranslationQuery,
			translation.Title,
			translation.Requirements,
			translation.Responsibilities,
			translation.Conditions,
			translation.Information,
			req.Id,
			translation.LanguageCode,
		)
		if err != nil {
			return fmt.Errorf("%s: failed to update translation for language %s: %w", op, translation.LanguageCode, err)
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			_, err = tx.Exec(ctx, insertTranslationQuery,
				req.Id,
				translation.LanguageCode,
				translation.Title,
				translation.Requirements,
				translation.Responsibilities,
				translation.Conditions,
				translation.Information,
			)
			if err != nil {
				return fmt.Errorf("%s: failed to insert translation for language %s: %w", op, translation.LanguageCode, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) RemoveVacancy(ctx context.Context, id int) error {
	const op = "postgres.RemoveVacancy"

	var exists bool
	checkVacancyQuery := `SELECT EXISTS(SELECT 1 FROM Vacancy WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkVacancyQuery, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check vacancy existence: %w", op, err)
	}
	if !exists {
		return storage.ErrVacancyNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	deleteTranslationQuery := `DELETE FROM VacancyTranslation WHERE vacancy_id = $1`
	_, err = tx.Exec(ctx, deleteTranslationQuery, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete vacancy translations: %w", op, err)
	}

	deleteVacancyQuery := `DELETE FROM Vacancy WHERE id = $1`
	_, err = tx.Exec(ctx, deleteVacancyQuery, id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete vacancy: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (db *DB) GetAllVacanciesByLanguage(ctx context.Context, languageCode string) ([]models.VacancyResponse, error) {
	const op = "postgres.GetAllVacanciesByLanguage"

	query := `
		SELECT v.id, vt.language_code, vt.title, vt.requirements, 
		       vt.responsibilities, vt.conditions, vt.information, v.isActive, v.salary
		FROM Vacancy v
		JOIN VacancyTranslation vt ON v.id = vt.vacancy_id
		WHERE vt.language_code = $1`

	rows, err := db.Pool.Query(ctx, query, languageCode)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query vacancies: %w", op, err)
	}
	defer rows.Close()

	var vacancies []models.VacancyResponse

	for rows.Next() {
		var vacancy models.VacancyResponse
		if err := rows.Scan(
			&vacancy.Id,
			&vacancy.LanguageCode,
			&vacancy.Title,
			&vacancy.Requirements,
			&vacancy.Responsibilities,
			&vacancy.Conditions,
			&vacancy.Information,
			&vacancy.IsActive,
			&vacancy.Salary,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row into vacancy struct: %w", op, err)
		}
		vacancies = append(vacancies, vacancy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	return vacancies, nil
}

func (db *DB) GetVacancyById(ctx context.Context, id int) (*models.VacancyResponses, error) {
	const op = "postgres.GetVacancyById"

	var exists bool
	checkVacancyQuery := `SELECT EXISTS(SELECT 1 FROM Vacancy WHERE id = $1)`
	err := db.Pool.QueryRow(ctx, checkVacancyQuery, id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to check vacancy existence: %w", op, err)
	}
	if !exists {
		return nil, storage.ErrVacancyNotFound
	}

	query := `
    SELECT v.salary, v.isActive, vt.language_code, vt.title, 
           vt.requirements, vt.responsibilities, vt.conditions, vt.information
    FROM Vacancy v
    JOIN VacancyTranslation vt ON v.id = vt.vacancy_id
    WHERE v.id = $1`

	rows, err := db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to query vacancy translations: %w", op, err)
	}
	defer rows.Close()

	var response models.VacancyResponses
	var translations []*models.CreateVacancy

	for rows.Next() {
		var translation models.CreateVacancy
		if err := rows.Scan(
			&response.Salary,
			&response.IsActive,
			&translation.LanguageCode,
			&translation.Title,
			&translation.Requirements,
			&translation.Responsibilities,
			&translation.Conditions,
			&translation.Information,
		); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row into struct: %w", op, err)
		}
		translations = append(translations, &translation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: row iteration error: %w", op, err)
	}

	// Заполнение ответа
	response.Vacancy = translations
	return &response, nil
}

func (db *DB) CreateVacancy(ctx context.Context, req *models.VacancyResponses) (*models.VacancyResponses, error) {
	const op = "postgres.CreateVacancy"

	if len(req.Vacancy) == 0 {
		return nil, storage.ErrRequiredLanguage
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var vacancyID int
	insertVacancy := `INSERT INTO Vacancy (salary, isActive) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(ctx, insertVacancy, req.Salary, req.IsActive).Scan(&vacancyID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to insert vacancy: %w", op, err)
	}

	insertTranslation := `
		INSERT INTO VacancyTranslation (
			vacancy_id, language_code, title, requirements, 
			responsibilities, conditions, information
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, translation := range req.Vacancy {
		var exists bool
		checkLanguageQuery := `SELECT EXISTS(SELECT 1 FROM Language WHERE code = $1)`
		err = tx.QueryRow(ctx, checkLanguageQuery, translation.LanguageCode).Scan(&exists)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to check language existence for code %s: %w", op, translation.LanguageCode, err)
		}
		if !exists {
			return nil, storage.ErrLanguageNotFound
		}

		_, err = tx.Exec(ctx, insertTranslation, vacancyID, translation.LanguageCode, translation.Title,
			translation.Requirements, translation.Responsibilities, translation.Conditions, translation.Information)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to insert vacancy translation for language %s: %w", op, translation.LanguageCode, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	response := &models.VacancyResponses{
		Salary:   req.Salary,
		IsActive: req.IsActive,
		Vacancy:  make([]*models.CreateVacancy, len(req.Vacancy)),
	}

	for i, translation := range req.Vacancy {
		response.Vacancy[i] = &models.CreateVacancy{
			Id:               vacancyID,
			LanguageCode:     translation.LanguageCode,
			Title:            translation.Title,
			Requirements:     translation.Requirements,
			Responsibilities: translation.Responsibilities,
			Conditions:       translation.Conditions,
			Information:      translation.Information,
		}
	}

	return response, nil
}

func (db *DB) SearchVacancies(ctx context.Context, query string) ([]models.VacancyResponse, error) {
	const op = "postgres.SearchVacancies"

	getVacancyIdsQuery := `
		SELECT 
			vacancy_id
		FROM 
			VacancyTranslation
		WHERE 
			language_code = 'ru' 
		  AND title ILIKE $1 
	`

	rows, err := db.Pool.Query(ctx, getVacancyIdsQuery, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get vacancy IDs: %w", op, err)
	}
	defer rows.Close()

	var vacancyIds []int
	for rows.Next() {
		var vacancyId int
		if err := rows.Scan(&vacancyId); err != nil {
			return nil, fmt.Errorf("%s: failed to scan vacancy ID: %w", op, err)
		}
		vacancyIds = append(vacancyIds, vacancyId)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	if len(vacancyIds) == 0 {
		return []models.VacancyResponse{}, nil
	}

	getVacanciesQuery := `
		SELECT 
			id, isActive, salary
		FROM 
			Vacancy
		WHERE 
			id = ANY($1)
	`

	rows, err = db.Pool.Query(ctx, getVacanciesQuery, vacancyIds)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get vacancies: %w", op, err)
	}
	defer rows.Close()

	vacancyMap := make(map[int]*models.VacancyResponse)
	for rows.Next() {
		var id int
		var isActive bool
		var salary int
		if err := rows.Scan(&id, &isActive, &salary); err != nil {
			return nil, fmt.Errorf("%s: failed to scan vacancy: %w", op, err)
		}
		vacancyMap[id] = &models.VacancyResponse{
			Id:       id,
			IsActive: isActive,
			Salary:   salary,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	getTranslationsQuery := `
		SELECT 
			vacancy_id, title, requirements, responsibilities, conditions, information
		FROM 
			VacancyTranslation
		WHERE 
			language_code = 'ru' AND vacancy_id = ANY($1)
	`

	rows, err = db.Pool.Query(ctx, getTranslationsQuery, vacancyIds)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get translations: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var requirements, responsibilities, conditions, information []string
		if err := rows.Scan(&id, &title, &requirements, &responsibilities, &conditions, &information); err != nil {
			return nil, fmt.Errorf("%s: failed to scan translation: %w", op, err)
		}

		if vacancy, ok := vacancyMap[id]; ok {
			vacancy.LanguageCode = "ru"
			vacancy.Title = title
			vacancy.Requirements = requirements
			vacancy.Responsibilities = responsibilities
			vacancy.Conditions = conditions
			vacancy.Information = information
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", op, err)
	}

	var result []models.VacancyResponse
	for _, vacancy := range vacancyMap {
		result = append(result, *vacancy)
	}

	return result, nil
}
