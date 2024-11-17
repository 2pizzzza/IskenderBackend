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

	//Discount
	GetAllDiscounts(ctx context.Context, languageCode string) ([]models.Discount, error)
	CreateDiscount(ctx context.Context, token string, discount models.DiscountCreate) (*models.DiscountCreate, error)
	DeleteDiscount(ctx context.Context, token string, discount models.DiscountRequest) error

	//Vacancy
	GetAllActiveVacancyByLang(ctx context.Context, code string) ([]models.VacancyResponse, error)
	UpdateVacancy(ctx context.Context, token string, req models.VacancyResponse) error
	RemoveVacancy(ctx context.Context, token string, req *models.RemoveVacancyRequest) error
	GetAllVacancyByLang(ctx context.Context, code string) ([]models.VacancyResponse, error)
	GetVacancyById(ctx context.Context, id int) (*models.VacancyResponses, error)
	CreateVacancy(ctx context.Context, token string, req *models.VacancyResponses) (*models.VacancyResponses, error)

	//Brand
	CreateBrand(ctx context.Context, token string, req *models.BrandRequest) (*models.BrandResponse, error)
	GetAllBrand(ctx context.Context) ([]*models.BrandResponse, error)
	RemoveBrand(ctx context.Context, token string, req *models.RemoveBrandRequest) error
	UpdateBrand(ctx context.Context, token string, id int, name, url string) (*models.BrandResponse, error)
	GetBrandById(ctx context.Context, id int) (*models.BrandResponse, error)

	//Review
	CreateReview(ctx context.Context, req *models.CreateReviewRequest) error
	GetAllReview(ctx context.Context) ([]*models.ReviewResponse, error)

	//Language
	GetAllLanguages(ctx context.Context) ([]*models.Language, error)

	//Category
	GetCategoriesByCode(ctx context.Context, languageCode string) ([]*models.Category, error)
	CreateCategory(ctx context.Context, token string, req models.CreateCategoryRequest) (*models.CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, token string, req *models.UpdateCategoryRequest) error
	RemoveCategory(ctx context.Context, token string, req *models.RemoveCategoryRequest) error
	GetCategoryById(ctx context.Context, id int) (*models.GetCategoryRequest, error)

	//Collection
	GetCollectionByCategoryId(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionId int, code string) (*models.CollectionResponse, error)
	RemoveCollection(ctx context.Context, token string, req *models.RemoveCollectionRequest) error
	UpdateCollection(ctx context.Context, token string, req *models.UpdateCollectionRequest) error
	GetCollectionRec(ctx context.Context, language string) ([]*models.CollectionResponse, error)
	GetCollectionByStadart(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	GetCollectionByPainted(ctx context.Context, code string) ([]*models.CollectionResponse, error)

	//Popular and New
	GetPopular(ctx context.Context, code string) (*models.PopularResponse, error)
	GetNew(ctx context.Context, code string) (*models.PopularResponse, error)

	//Items
	GetItemsByCategoryId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemById(ctx context.Context, id int, code string) (*models.ItemResponse, error)
	GetItemsByCollectionId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemsRec(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)

	//Seach
	Search(ctx context.Context, code string, isProducer *bool, isPainted *bool, searchQuery string) (*models.PopularResponse, error)

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

	//Discount
	mux.HandleFunc("GET /discounts", s.GetAllDiscount)
	mux.HandleFunc("POST /discount", s.CreateDiscount)
	mux.HandleFunc("DELETE /discount", s.RemoveDiscount)

	//Vacancy
	mux.HandleFunc("GET /vacancies/activ", s.GetAllVacancyActive)
	mux.HandleFunc("DELETE /vacancy", s.RemoveVacancy)
	mux.HandleFunc("PUT /vacancy", s.UpdateVacancy)
	mux.HandleFunc("GET /vacancies", s.GetAllVacancy)
	mux.HandleFunc("GET /vacancy", s.GetVacancyById)
	mux.HandleFunc("POST /vacancy", s.CreateVacancy)

	//Brand
	mux.HandleFunc("GET /brands", s.GetAllBrands)
	mux.HandleFunc("POST /brand", s.CreateBrand)
	mux.HandleFunc("DELETE /brand", s.RemoveBrand)
	mux.HandleFunc("PUT /brand", s.UpdateBrand)
	mux.HandleFunc("GET /brand", s.GetBrandById)

	//Review
	mux.HandleFunc("GET /reviews", s.GetAllReviews)
	mux.HandleFunc("POST /reviews", s.CreateReview)

	//Language
	mux.HandleFunc("GET /languages", s.GetAllLanguages)

	//Category
	mux.HandleFunc("GET /category", s.GetAllCategoriesByCode)
	mux.HandleFunc("POST /category", s.CreateCategory)
	mux.HandleFunc("PUT /category", s.UpdateCategory)
	mux.HandleFunc("DELETE /category", s.RemoveCategory)
	mux.HandleFunc("GET /category/by/id", s.GetCategoryById)

	//Collection
	mux.HandleFunc("GET /collections", s.GetCollectionsByCategoryId)
	mux.HandleFunc("GET /collection", s.GetCollectionById)
	mux.HandleFunc("DELETE /collection", s.RemoveCollection)
	mux.HandleFunc("GET /collections/rec", s.GetCollectionsRec)
	mux.HandleFunc("GET /collections/standart", s.GetCollectionsStandart)
	mux.HandleFunc("GET /collections/painted", s.GetCollectionsByPainted)

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
