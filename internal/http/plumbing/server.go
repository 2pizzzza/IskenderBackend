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
	UpdateVacancy(ctx context.Context, token string, req models.VacancyUpdateRequest) error
	RemoveVacancy(ctx context.Context, token string, req *models.RemoveVacancyRequest) error
	GetAllVacancyByLang(ctx context.Context, code string) ([]models.VacancyResponse, error)
	GetVacancyById(ctx context.Context, id int) (*models.VacancyResponses, error)
	CreateVacancy(ctx context.Context, token string, req *models.VacancyResponses) (*models.VacancyResponses, error)
	SearchVacancy(ctx context.Context, query string) ([]models.VacancyResponse, error)

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
	UpdateCategory(ctx context.Context, token string, categoryID int, req []models.UpdateCategoriesResponse) error
	RemoveCategory(ctx context.Context, token string, req *models.RemoveCategoryRequest) error
	GetCategoryById(ctx context.Context, id int) (*models.GetCategoryRequest, error)

	//Collection
	GetCollectionByCategoryId(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionId int, code string) (*models.CollectionResponse, error)
	RemoveCollection(ctx context.Context, token string, req *models.RemoveCollectionRequest) error
	UpdateCollection(ctx context.Context, token string, collectionId int, req models.CreateCollectionRequest) error
	GetCollectionRec(ctx context.Context, language string) ([]*models.CollectionResponse, error)
	GetCollectionByStadart(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	GetCollectionByPainted(ctx context.Context, code string) ([]*models.CollectionResponse, error)
	CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.CreateCollectionResponse, error)
	GetCollection(ctx context.Context) ([]*models.CollectionResponses, error)
	GetCollectionID(ctx context.Context, collectionId int) (*models.CollectionResponseForAdmin, error)

	//Popular and New
	GetPopular(ctx context.Context, code string) (*models.PopularResponse, error)
	GetNew(ctx context.Context, code string) (*models.PopularResponse, error)

	//Items
	GetItemsByCategoryId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemById(ctx context.Context, id int, code string) (*models.ItemResponse, error)
	GetItemsByCollectionId(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	GetItemsRec(ctx context.Context, id int, code string) ([]*models.ItemResponse, error)
	CreateItem(ctx context.Context, req models.CreateItem) (*models.CreateItemResponse, error)
	UpdateItem(ctx context.Context, token string, itemId int, req models.CreateItem) error
	RemoveItem(ctx context.Context, token string, req models.ItemRequest) error
	GetItems(ctx context.Context) ([]*models.ItemResponses, error)
	GetItemID(ctx context.Context, itemId int) (*models.ItemResponseForAdmin, error)

	//Seach
	Search(ctx context.Context, code string, isProducer *bool, isPainted *bool, searchQuery string, minPrice, maxPrice *float64) (*models.PopularResponse, error)
	SearchCollection(ctx context.Context, code string, isProducer *bool, isPainted *bool, searchQuery string, minPrice, maxPrice *float64) ([]*models.CollectionResponse, error)
	SearchItem(ctx context.Context, code string, isProducer *bool, isPainted *bool, searchQuery string, minPrice, maxPrice *float64) ([]*models.ItemResponse, error)

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
	mux.HandleFunc("POST /api/starter", s.Starter)

	//Discount
	mux.HandleFunc("GET /api/discounts", s.GetAllDiscount)
	mux.HandleFunc("POST /api/discount", s.CreateDiscount)
	mux.HandleFunc("DELETE /api/discount", s.RemoveDiscount)

	//Vacancy
	mux.HandleFunc("GET /api/vacancies/activ", s.GetAllVacancyActive)
	mux.HandleFunc("DELETE /api/vacancy", s.RemoveVacancy)
	mux.HandleFunc("PUT /api/vacancy", s.UpdateVacancy)
	mux.HandleFunc("GET /api/vacancies", s.GetAllVacancy)
	mux.HandleFunc("GET /api/vacancy", s.GetVacancyById)
	mux.HandleFunc("GET /api/searchVacancy", s.SearchVacancy)
	mux.HandleFunc("POST /api/vacancy", s.CreateVacancy)

	//Brand
	mux.HandleFunc("GET /api/brands", s.GetAllBrands)
	mux.HandleFunc("POST /api/brand", s.CreateBrand)
	mux.HandleFunc("DELETE /api/brand", s.RemoveBrand)
	mux.HandleFunc("PUT /api/brand", s.UpdateBrand)
	mux.HandleFunc("GET /api/brand", s.GetBrandById)

	//Review
	mux.HandleFunc("GET /api/reviews", s.GetAllReviews)
	mux.HandleFunc("POST /api/reviews", s.CreateReview)

	//Language
	mux.HandleFunc("GET /api/languages", s.GetAllLanguages)

	//Category
	mux.HandleFunc("GET /api/category", s.GetAllCategoriesByCode)
	mux.HandleFunc("POST /api/category", s.CreateCategory)
	mux.HandleFunc("PUT /api/category", s.UpdateCategory)
	mux.HandleFunc("DELETE /api/category", s.RemoveCategory)
	mux.HandleFunc("GET /api/category/by/id", s.GetCategoryById)

	//Collection
	mux.HandleFunc("GET /api/collections", s.GetCollectionsByCategoryId)
	mux.HandleFunc("GET /api/collection", s.GetCollectionById)
	mux.HandleFunc("DELETE /api/collection", s.RemoveCollection)
	mux.HandleFunc("GET /api/collections/rec", s.GetCollectionsRec)
	mux.HandleFunc("GET /api/collections/standart", s.GetCollectionsStandart)
	mux.HandleFunc("GET /api/collections/painted", s.GetCollectionsByPainted)
	mux.HandleFunc("POST /api/collection", s.CreateCollection)
	mux.HandleFunc("PUT /api/collection", s.UpdateCollection)
	mux.HandleFunc("GET /api/getAllCollection", s.GetAllCollection)
	mux.HandleFunc("GET /api/getCollectionById", s.GetCollectionId)

	//Popular and New
	mux.HandleFunc("GET /api/popular", s.GetPopular)
	mux.HandleFunc("GET /api/new", s.GetNew)

	//Item
	mux.HandleFunc("GET /api/items", s.GetItemsByCategoryId)
	mux.HandleFunc("GET /api/item", s.GetItemsById)
	mux.HandleFunc("GET /api/items/collection", s.GetItemsByCollectionId)
	mux.HandleFunc("GET /api/items/rec", s.GetItemsRec)
	mux.HandleFunc("POST /api/items", s.CreateItem)
	mux.HandleFunc("PUT /api/items", s.UpdateItem)
	mux.HandleFunc("DELETE /api/items", s.RemoveItem)
	mux.HandleFunc("GET /api/getAllItems", s.GetAllItems)
	mux.HandleFunc("GET /api/getItemById", s.GetItemId)

	//Search
	mux.HandleFunc("GET /api/search", s.Search)
	mux.HandleFunc("GET /api/searchItems", s.SearchItems)
	mux.HandleFunc("GET /api/searchCollections", s.SearchCollections)

	//Photo
	mux.HandleFunc("GET /media/images/", s.GetImage)

}
