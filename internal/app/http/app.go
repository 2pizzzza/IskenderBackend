package http

import (
	"fmt"
	_ "github.com/2pizzzza/plumbing/cmd/plumbing/docs"
	authService "github.com/2pizzzza/plumbing/internal/http/auth"
	plumbingRouters "github.com/2pizzzza/plumbing/internal/http/plumbing"
	"github.com/2pizzzza/plumbing/internal/lib/logger/sl"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	host       string
	port       int
}

func New(log *slog.Logger, host string, port int, app *plumbingRouters.Server, auth *authService.Server) *App {
	mux := http.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "http://64.176.71.25:5174", "https://garant-admin.vercel.app", "http://127.0.0.1:5173/", "http://garant-asia.com", "http://64.176.71.25:8081/"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(mux)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	app.RegisterRoutes(mux)
	auth.RegisterRoutes(mux)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handler,
	}

	return &App{
		log:        log,
		httpServer: httpServer,
		host:       host,
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
		slog.String("port", a.host),
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
