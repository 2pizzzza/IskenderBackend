package http

import (
	"fmt"
	plumbingRouters "github.com/2pizzzza/plumbing/internal/http/plumbing"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	port       int
}

func New(log *slog.Logger, port int, app *plumbingRouters.Server) *App {
	mux := http.NewServeMux()

	//app.RegisterRoutes(mux)
	app.RegisterRoutes(mux)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &App{
		log:        log,
		httpServer: httpServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("starting HTTP server")

	if err := a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "httpapp.Stop"

	a.log.With(slog.String("op", op)).Info("Stopping HTTP server", slog.Int("port", a.port))

	if err := a.httpServer.Close(); err != nil {
		a.log.Error("error while closing HTTP server", sl.Err(err))
	}
}
