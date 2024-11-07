package main

import (
	"github.com/2pizzzza/plumbing/internal/app"
	"github.com/2pizzzza/plumbing/internal/config"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.MustLoad()

	if err != nil {
		slog.Error("Failed load env", sl.Err(err))
	}

	log := setupLogger(cfg.Env)

	log.Info("Starting Apllication")

	dir, err := getBaseDir()
	if err != nil {
		log.Error("failed to get base dir")
	}
	application := app.New(log, dir, cfg)

	go application.HTTPserv.MustRun()

	stop := make(chan os.Signal)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal:", sign.String()))

	application.HTTPserv.Stop()

	log.Info("Server is dead")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(

			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func getBaseDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir, nil
}
