package app

import (
	httpapp "github.com/2pizzzza/plumbing/internal/app/http"
	"github.com/2pizzzza/plumbing/internal/config"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/storage/postgres"
	"log/slog"
)

type App struct {
	HTTPserv *httpapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	db, err := postgres.New(cfg)

	_ = db
	if err != nil {
		log.Error("Failed to connect to database", sl.Err(err))
		return nil
	}

	httpApp := httpapp.New(log, cfg.HttpPort)

	return &App{
		HTTPserv: httpApp,
	}
}
