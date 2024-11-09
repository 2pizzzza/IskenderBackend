package service

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
)

type Plumping struct {
	log                *slog.Logger
	baseDir            string
	plumpingRepository PlumpingRepository
}

type PlumpingRepository interface {
	//Language
	GetLanguages(ctx context.Context) ([]*models.Language, error)
}

func New(log *slog.Logger, baseDir string, repository PlumpingRepository) *Plumping {
	return &Plumping{
		log:                log,
		baseDir:            baseDir,
		plumpingRepository: repository,
	}
}
