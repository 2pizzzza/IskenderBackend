package plumbing

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
	"net/http"
)

type Service interface {
	//Language
	GetAllLanguages(ctx context.Context) ([]*models.Language, error)
}

type Server struct {
	log     *slog.Logger
	service Service
}

func New(log *slog.Logger, service Service) *Server {
	return &Server{
		log:     log,
		service: service,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	//Language
	mux.HandleFunc("GET /language", s.GetAllLanguages)

}
