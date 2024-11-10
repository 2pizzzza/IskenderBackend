package plumbing

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
)

func (pr *Plumping) GetAllLanguages(ctx context.Context) ([]*models.Language, error) {
	const op = "service.GetAllLanguages"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	languages, err := pr.plumpingRepository.GetLanguages(ctx)
	if err != nil {
		log.Error("Failed to get all languages", sl.Err(err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return languages, nil
}
