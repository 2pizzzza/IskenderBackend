package plumbing

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/schemas"
	"log/slog"
	"net/http"
)

type Service interface {
	//Catalog
	CreateCatalog(ctx context.Context, req *schemas.CreateCatalogRequest) (res *schemas.CreateCatalogResponse, err error)
	AddNewCatalogLocalization(ctx context.Context, req *schemas.CatalogLocalizationRequest) (*schemas.CatalogLocalization, error)
	GetCatalogsByLangCode(ctx context.Context, req *schemas.CatalogsByLanguageRequest) ([]*schemas.CatalogResponse, error)
	RemoveCatalog(ctx context.Context, req *schemas.CatalogRemoveRequest) error
	UpdateCatalog(ctx context.Context, req *schemas.UpdateCatalogRequest) error
	GetCatalogById(ctx context.Context, id int) (*schemas.CatalogDetailResponse, error)
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
	//Catalog
	mux.HandleFunc("POST /catalog", s.CreateCatalog)
	mux.HandleFunc("POST /catalog-localization", s.CreateNewLocalizationForCatalog)
	mux.HandleFunc("GET /catalogs/by-language-code", s.GetAllCatalogsByLangCode)
	mux.HandleFunc("DELETE /catalog", s.RemoveCatalog)
	mux.HandleFunc("PUT /catalog", s.UpdateCatalog)
	mux.HandleFunc("GET /catalog", s.GetCatalogById)
}
