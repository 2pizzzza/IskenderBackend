package app

import (
	httpapp "github.com/2pizzzza/plumbing/internal/app/http"
	"github.com/2pizzzza/plumbing/internal/config"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/service"
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

	if err != nil {
		log.Error("Failed connect db err: %s", sl.Err(err))
	}

	plumbingService := service.New(log, db)
	_ = plumbingService

	httpApp := httpapp.New(log, cfg.HttpPort)

	return &App{
		HTTPserv: httpApp,
	}
}
