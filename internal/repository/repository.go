package repository

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateAuth(auth *domain.Auth) (string, error)
	GetAuth(auth *domain.Auth) (domain.Auth, error)
	CheckExistAuth(auth *domain.SignInInput) (domain.Auth, error)
	GetByCredentials(auth *domain.SignInInput) (domain.Auth, error)
	SetSession(authID primitive.ObjectID, session domain.Session) error
	VerificationCode(userID string, code string) error
	RefreshToken(refreshToken string) (domain.Auth, error)
}

type Shop interface {
	FindShop(params domain.RequestParams) (domain.Response[domain.Shop], error)
	GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error)
	CreateShop(userID string, shop *domain.Shop) (*domain.Shop, error)
}

type Page interface {
	GetPageForRouters() (domain.Response[domain.PageRoutes], error)
	GetFullPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	FindPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	GetPage(id string) (domain.Page, error)
	CreatePage(userID string, page *domain.Page) (*domain.Page, error)
	DeletePage(id string) (domain.Page, error)
	UpdatePage(id string, data interface{}) (domain.Page, error)
}

type Component interface {
	GetComponent(id string) (domain.Component, error)
	FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	CreateComponent(userID string, component *domain.ComponentCreate) (*domain.Component, error)
	DeleteComponent(id string) (domain.Component, error)
	UpdateComponent(id string, data interface{}) (domain.Component, error)

	FindByPopulate(params domain.RequestParams) (domain.Response[domain.Component], error)
	FindLibrarys(params domain.RequestParams) (domain.Response[domain.Library], error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userID string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user *domain.User) (domain.User, error)
	Iam(userID string) (domain.User, error)
}

type Apps interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Repositories struct {
	Authorization
	Apps
	Component
	Page
	Shop
	User
}

func NewRepositories(mongodb *mongo.Database, i18n config.I18nConfig) *Repositories {
	return &Repositories{
		Authorization: NewAuthMongo(mongodb),
		Shop:          NewShopMongo(mongodb),
		User:          NewUserMongo(mongodb),
		Page:          NewPageMongo(mongodb, i18n),
		Component:     NewComponentMongo(mongodb, i18n),
		Apps:          NewAppsMongo(mongodb, i18n),
	}
}
