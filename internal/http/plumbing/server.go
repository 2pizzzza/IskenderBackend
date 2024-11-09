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

	//Category
	GetCategoriesByCode(ctx context.Context, languageCode string) ([]*models.Category, error)

	//Collection
	GetCollectionByCategoryId(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionId int, code string) (*models.CollectionResponse, error)

	//Popular and New
	GetPopular(ctx context.Context, code string) (*models.PopularResponse, error)
	GetNew(ctx context.Context, code string) (*models.NewResponse, error)

	//Items
	GetItemsByCategoryId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemById(ctx context.Context, id int, code string) (*models.ItemResponse, error)
	GetItemsByCollectionId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
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

	//Category
	mux.HandleFunc("GET /category", s.GetAllCategoriesByCode)

	//Collection
	mux.HandleFunc("GET /collections", s.GetCollectionsByCategoryId)
	mux.HandleFunc("GET /collection", s.GetCollectionById)

	//Popular and New
	mux.HandleFunc("GET /popular", s.GetPopular)
	mux.HandleFunc("GET /new", s.GetNew)

	//Item
	mux.HandleFunc("GET /items", s.GetItemsByCategoryId)
	mux.HandleFunc("GET /item", s.GetItemsById)
	mux.HandleFunc("GET /items-collection", s.GetItemsByCollectionId)

}
