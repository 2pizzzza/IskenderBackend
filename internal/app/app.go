package app

import (
	httpapp "github.com/2pizzzza/plumbing/internal/app/http"
	"github.com/2pizzzza/plumbing/internal/config"
	"github.com/2pizzzza/plumbing/internal/http/plumbing"
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
	baseDir string,
	cfg *config.Config,
) *App {
	db, err := postgres.New(cfg)

	if err != nil {
		log.Error("Failed connect db err: %s", sl.Err(err))
	}

	plumbingRepository := service.New(log, baseDir, db)
	plumbingService := plumbing.New(log, plumbingRepository)

	httpApp := httpapp.New(log, cfg.HttpPort, plumbingService)

	return &App{
		HTTPserv: httpApp,
	}
}
