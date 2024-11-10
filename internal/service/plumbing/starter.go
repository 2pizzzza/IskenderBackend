package plumbing

import (
	"context"
	"errors"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage"
	"log/slog"
)

func (pr *Plumping) Starter(ctx context.Context) error {
	const op = "service.Starter"

	log := pr.log.With(
		slog.String("op: ", op),
	)

	err := pr.plumpingRepository.CreateStarter(ctx)

	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			log.Error("Already exist", sl.Err(err))
			return storage.ErrCategoryNotFound
		}
		log.Error("Failed to create starter", sl.Err(err))
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}
