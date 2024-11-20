package plumbing

import (
	"context"
	"github.com/2pizzzza/plumbing/internal/domain/models"
	"log/slog"
)

type Plumping struct {
	log                *slog.Logger
	baseDir            string
	plumpingRepository PlumpingRepository
}

type PlumpingRepository interface {

	//Discount
	GetAllDiscount(ctx context.Context, languageCode string) ([]models.Discount, error)
	CreateDiscount(ctx context.Context, discount models.DiscountCreate) (*models.DiscountCreate, error)
	DeleteDiscount(ctx context.Context, id int) error

	//Brand
	CreateBrand(ctx context.Context, name, url string) (*models.BrandResponse, error)
	GetAllBrand(ctx context.Context) ([]*models.BrandResponse, error)
	RemoveBrand(ctx context.Context, id int) error
	UpdateBrand(ctx context.Context, id int, name, url string) (*models.BrandResponse, error)
	GetBrandByID(ctx context.Context, id int) (*models.BrandResponse, error)

	//Vacancy
	GetAllActiveVacanciesByLanguage(ctx context.Context, languageCode string) ([]models.VacancyResponse, error)
	UpdateVacancy(ctx context.Context, req models.VacancyResponse) error
	RemoveVacancy(ctx context.Context, id int) error
	GetAllVacanciesByLanguage(ctx context.Context, languageCode string) ([]models.VacancyResponse, error)
	GetVacancyById(ctx context.Context, id int) (*models.VacancyResponses, error)
	CreateVacancy(ctx context.Context, req *models.VacancyResponses) (*models.VacancyResponses, error)

	//Review
	CreateReview(ctx context.Context, username string, rating int, text string) error
	GetAllReviews(ctx context.Context) ([]*models.ReviewResponse, error)

	//Starter
	CreateStarter(ctx context.Context) error

	//Language
	GetLanguages(ctx context.Context) ([]*models.Language, error)

	//Category
	GetCategoriesByLanguageCode(ctx context.Context, languageCode string) ([]*models.Category, error)
	UpdateCategory(ctx context.Context, categoryID int, name string, languageCode string) error
	CreateCategory(ctx context.Context, req models.CreateCategoryRequest) (*models.CreateCategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryID int) error
	GetCategoryByID(ctx context.Context, id int) (*models.GetCategoryRequest, error)

	//Collection
	GetCollectionsByLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetCollectionByID(ctx context.Context, collectionID int, languageCode string) (*models.CollectionResponse, error)
	DeleteCollection(ctx context.Context, collectionID int) error
	UpdateCollection(ctx context.Context, collectionID int, req models.CreateCollectionRequest) error
	GetRandomCollectionsWithPopularity(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetCollectionsByIsProducerSLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetCollectionsByIsProducerPLanguageCode(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	CreateCollection(ctx context.Context, req models.CreateCollectionRequest) (*models.CreateCollectionResponse, error)

	//Popular and new
	GetPopularCollections(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetPopularItems(ctx context.Context, languageCode string) ([]*models.ItemResponse, error)
	GetNewCollections(ctx context.Context, languageCode string) ([]*models.CollectionResponse, error)
	GetNewItems(ctx context.Context, languageCode string) ([]*models.ItemResponse, error)

	//Search and filtr
	SearchCollections(ctx context.Context, languageCode string, isProducer *bool, isPainted *bool, searchQuery string) ([]*models.CollectionResponse, error)
	SearchItems(ctx context.Context, languageCode string, isProducer *bool, isPainted *bool, searchQuery string) ([]*models.ItemResponse, error)

	//Item
	GetItemsByCategoryID(ctx context.Context, categoryID int, languageCode string) ([]*models.ItemResponse, error)
	GetItemByID(ctx context.Context, itemID int, languageCode string) (*models.ItemResponse, error)
	GetItemsByCollectionID(ctx context.Context, collectionID int, languageCode string) ([]*models.ItemResponse, error)
	GetRandomItemsWithPopularity(ctx context.Context, languageCode string, itemID int) ([]*models.ItemResponse, error)
	CreateItem(ctx context.Context, req models.CreateItem) (*models.CreateItemResponse, error)
	UpdateItem(ctx context.Context, itemID int, req models.CreateItem) error
	RemoveItem(ctx context.Context, itemID int) error
	GetAllItems(ctx context.Context) ([]*models.ItemResponses, error)
}

func New(log *slog.Logger, baseDir string, repository PlumpingRepository) *Plumping {
	return &Plumping{
		log:                log,
		baseDir:            baseDir,
		plumpingRepository: repository,
	}
}
