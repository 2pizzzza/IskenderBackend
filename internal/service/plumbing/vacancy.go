package plumbing

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	token2 "github.com/2pizzzza/plumbing/internal/lib/jwt"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"log/slog"
)

func (pr *Plumping) GetAllActiveVacancyByLang(ctx context.Context, code string) ([]models.VacancyResponse, error) {
	const op = "service.GetAllActiveVacancyByLang"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	vacancy, err := pr.plumpingRepository.GetAllActiveVacanciesByLanguage(ctx, code)
	if err != nil {
		log.Error("Failed to get all vacancy", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if vacancy == nil {
		vacancy = []models.VacancyResponse{}
	}
	return vacancy, nil
}

func (pr *Plumping) UpdateVacancy(ctx context.Context, token string, req models.VacancyUpdateRequest) error {
	const op = "service.UpdateVacancy"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.UpdateVacancy(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrVacancyNotFound) {
			log.Error("Vacancy not found", sl.Err(err))
			return storage.ErrVacancyNotFound
		}
		if errors.Is(err, storage.ErrRequiredLanguage) {
			log.Error("Required 3 languages", sl.Err(err))
			return storage.ErrRequiredLanguage
		}
		if errors.Is(err, storage.ErrInvalidLanguageCode) {
			log.Error("Required 3 languages kgz, ru, en", sl.Err(err))
			return storage.ErrInvalidLanguageCode
		}
		log.Error("Failed to update brand", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) RemoveVacancy(ctx context.Context, token string, req *models.RemoveVacancyRequest) error {
	const op = "service.RemoveVacancy"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return storage.ErrToken
	}

	err = pr.plumpingRepository.RemoveVacancy(ctx, req.ID)
	if err != nil {
		if errors.Is(err, storage.ErrVacancyNotFound) {
			log.Error("Brand not found", sl.Err(err))
			return storage.ErrVacancyNotFound

		}
		log.Error("Failed to remove vacancy", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (pr *Plumping) GetAllVacancyByLang(ctx context.Context, code string) ([]models.VacancyResponse, error) {
	const op = "service.GetAllActiveVacancyByLang"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	vacancy, err := pr.plumpingRepository.GetAllVacanciesByLanguage(ctx, code)
	if err != nil {
		log.Error("Failed to get all vacancy", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	if vacancy == nil {
		vacancy = []models.VacancyResponse{}
	}
	return vacancy, nil
}

func (pr Plumping) GetVacancyById(ctx context.Context, id int) (*models.VacancyResponses, error) {
	const op = "service.GetVacancyById"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	brand, err := pr.plumpingRepository.GetVacancyById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrVacancyNotFound) {
			log.Error("Vacancy not found", sl.Err(err))
			return nil, storage.ErrVacancyNotFound

		}
		log.Error("Failed to get brand", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return brand, nil
}

func (pr *Plumping) CreateVacancy(ctx context.Context, token string, req *models.VacancyResponses) (*models.VacancyResponses, error) {
	const op = "service.RemoveVacancy"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	_, err := token2.ValidateToken(token)
	if err != nil {
		log.Error("Failed validate token")
		return nil, storage.ErrToken
	}
	log.Info("amit", req.Vacancy)
	vacancy, err := pr.plumpingRepository.CreateVacancy(ctx, req)
	if err != nil {
		if errors.Is(err, storage.ErrRequiredLanguage) {
			log.Error("Required 3 language", sl.Err(err))
			return nil, storage.ErrRequiredLanguage

		}
		if errors.Is(err, storage.ErrLanguageNotFound) {
			log.Error("Language not found", sl.Err(err))
			return nil, storage.ErrLanguageNotFound

		}
		log.Error("Failed to create vacancy", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return vacancy, nil
}
