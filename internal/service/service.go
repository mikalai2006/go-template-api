package service

import (
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/auths"
	"github.com/mikalai2006/go-template-api/pkg/hasher"
)

type Authorization interface {
	CreateAuth(auth *domain.SignInInput) (string, error)
	SignIn(input *domain.SignInInput) (domain.ResponseTokens, error)
	ExistAuth(auth *domain.SignInInput) (domain.Auth, error)
	CreateSession(auth *domain.Auth) (domain.ResponseTokens, error)
	VerificationCode(userID string, code string) error
	RefreshTokens(refreshToken string) (domain.ResponseTokens, error)
}

type Shop interface {
	FindShop(params domain.RequestParams) (domain.Response[domain.Shop], error)

	GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error)
	CreateShop(userID string, shop *domain.Shop) (*domain.Shop, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userID string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user *domain.User) (domain.User, error)
	Iam(userID string) (domain.User, error)
}

type Page interface {
	GetPageForRouters() (domain.Response[domain.PageRoutes], error)
	GetPage(id string) (domain.Page, error)
	GetFullPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	FindPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	CreatePage(userID string, page *domain.PageInputData) (*domain.Page, error)
	DeletePage(id string) (domain.Page, error)
	UpdatePage(id string, data interface{}) (domain.Page, error)
	UpdatePageWithContent(id string, data map[string]interface{}) (domain.Page, error)
}

type Component interface {
	GetComponent(id string) (domain.Component, error)
	FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	CreateComponent(userID string, component *domain.ComponentInput) (*domain.Component, error)
	DeleteComponent(id string) (domain.Component, error)
	UpdateComponent(id string, data interface{}) (domain.Component, error)
	FindByPopulate(params domain.RequestParams) (domain.Response[domain.Component], error)

	FindLibrarys(params domain.RequestParams) (domain.Response[domain.Library], error)
}

type ComponentGroup interface {
	FindComponentGroup() (domain.Response[domain.ComponentGroup], error)
	CreateComponentGroup(userID string, component *domain.ComponentGroup) (*domain.ComponentGroup, error)
	UpdateComponentGroup(id string, data interface{}) (domain.ComponentGroup, error)
	DeleteComponentGroup(id string) (domain.ComponentGroup, error)
}

type Product interface {
	CreateProduct(userID string, data *domain.ProductInput) (domain.Product, error)
	GetProduct(id string) (domain.Product, error)
	FindProduct(params domain.RequestParams) (domain.Response[domain.Product], error)
	UpdateProduct(id string, data interface{}) (domain.Product, error)
	DeleteProduct(id string) (domain.Product, error)
}

type Apps interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Services struct {
	Authorization
	Apps
	Component
	ComponentGroup
	Page
	Product
	Shop
	User
}

type ConfigServices struct {
	Repositories           *repository.Repositories
	Hasher                 hasher.PasswordHasher
	TokenManager           auths.TokenManager
	OtpGenerator           utils.Generator
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeLength int
	I18n                   config.I18nConfig
}

func NewServices(cfgService *ConfigServices) *Services {
	return &Services{
		Authorization: NewAuthService(
			cfgService.Repositories.Authorization,
			cfgService.Hasher,
			cfgService.TokenManager,
			cfgService.RefreshTokenTTL,
			cfgService.AccessTokenTTL,
			cfgService.OtpGenerator,
			cfgService.VerificationCodeLength,
		),
		Shop:           NewShopService(cfgService.Repositories.Shop),
		Apps:           NewAppsService(cfgService.Repositories, cfgService.I18n),
		Component:      NewComponentService(cfgService.Repositories.Component, cfgService.I18n),
		ComponentGroup: NewComponentGroupService(cfgService.Repositories.ComponentGroup, cfgService.I18n),
		Page:           NewPageService(cfgService.Repositories.Page, cfgService.I18n),
		Product:        NewProductService(cfgService.Repositories, cfgService.I18n),
		User:           NewUserService(cfgService.Repositories.User),
	}
}
