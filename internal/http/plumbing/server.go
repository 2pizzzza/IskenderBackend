package plumbing

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"log/slog"
	"net/http"
)

type PlumbingService interface {
	// Category
	CreateCategory(ctx context.Context, req *schemas.CreateCategoryRequest) (res *schemas.CreateCategoryResponse, err error)

	//Item
	CreateItem(ctx context.Context, req *schemas.CreateItemRequest) (res *schemas.CreateItemResponse, err error)
}

type Server struct {
	log     *slog.Logger
	service PlumbingService
}

func New(log *slog.Logger, service PlumbingService) *Server {
	return &Server{
		log:     log,
		service: service,
	}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /create-item", s.CreateItem)
	mux.HandleFunc("POST /create-category", s.CreateCategory)
}
