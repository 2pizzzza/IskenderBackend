package app

import (
	httpapp "github.com/2pizzzza/plumbing/internal/app/http"
	"github.com/2pizzzza/plumbing/internal/config"
	authService "github.com/2pizzzza/plumbing/internal/http/auth"
	handlerPlumbing "github.com/2pizzzza/plumbing/internal/http/plumbing"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/2pizzzza/plumbing/internal/service/auth"
	service "github.com/2pizzzza/plumbing/internal/service/plumbing"
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
	plumbingService := handlerPlumbing.New(log, plumbingRepository)

	authRepository := auth.New(log, cfg.DBHost, db)
	authServ := authService.New(log, authRepository)

	httpApp := httpapp.New(log, cfg.HttpHost, cfg.HttpPort, plumbingService, authServ)

	return &App{
		HTTPserv: httpApp,
	}
}
