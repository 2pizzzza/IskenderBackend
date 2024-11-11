package plumbing

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type Service interface {

	//Starter
	Starter(ctx context.Context) error

	//Language
	GetAllLanguages(ctx context.Context) ([]*models.Language, error)

	//Category
	GetCategoriesByCode(ctx context.Context, languageCode string) ([]*models.Category, error)
	CreateCategory(ctx context.Context, token string, req models.CreateCategoryRequest) (*models.CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, token string, req *models.UpdateCategoryRequest) error
	RemoveCategory(ctx context.Context, token string, req *models.RemoveCategoryRequest) error

	//Collection
	GetCollectionByCategoryId(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionId int, code string) (*models.CollectionResponse, error)
	RemoveCollection(ctx context.Context, token string, req *models.RemoveCollectionRequest) error
	UpdateCollection(ctx context.Context, token string, req *models.UpdateCollectionRequest) error
	GetCollectionRec(ctx context.Context, language string) ([]*models.CollectionResponse, error)

	//Popular and New
	GetPopular(ctx context.Context, code string) (*models.PopularResponse, error)
	GetNew(ctx context.Context, code string) (*models.PopularResponse, error)

	//Items
	GetItemsByCategoryId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemById(ctx context.Context, id int, code string) (*models.ItemResponse, error)
	GetItemsByCollectionId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemsRec(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)

	//Seach
	Search(ctx context.Context, code string, isProducer *bool, searchQuery string) (*models.PopularResponse, error)

	//Photo
	GetImagePath(ctx context.Context, imageName string) (string, error)
	UploadPhotos(ctx context.Context, files []*multipart.FileHeader) ([]string, error)
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
	//Starter
	mux.HandleFunc("POST /starter", s.Starter)
	//Language
	mux.HandleFunc("GET /languages", s.GetAllLanguages)

	//Category
	mux.HandleFunc("GET /category", s.GetAllCategoriesByCode)
	mux.HandleFunc("POST /category", s.CreateCategory)
	mux.HandleFunc("PUT /category", s.UpdateCategory)
	mux.HandleFunc("DELETE /category", s.RemoveCategory)

	//Collection
	mux.HandleFunc("GET /collections", s.GetCollectionsByCategoryId)
	mux.HandleFunc("GET /collection", s.GetCollectionById)
	mux.HandleFunc("DELETE /collection", s.RemoveCollection)
	mux.HandleFunc("GET /collections/rec", s.GetCollectionsRec)

	//Popular and New
	mux.HandleFunc("GET /popular", s.GetPopular)
	mux.HandleFunc("GET /new", s.GetNew)

	//Item
	mux.HandleFunc("GET /items", s.GetItemsByCategoryId)
	mux.HandleFunc("GET /item", s.GetItemsById)
	mux.HandleFunc("GET /items/collection", s.GetItemsByCollectionId)
	mux.HandleFunc("GET /items/rec", s.GetItemsRec)

	//Search
	mux.HandleFunc("GET /search", s.Search)

	//Photo
	mux.HandleFunc("GET /media/images/", s.GetImage)

}
