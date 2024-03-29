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
	RemoveRefreshTokens(refreshToken string) (string, error)
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
	// GetPage(id string) (domain.Page, error)
	GetStory(params domain.RequestParams) (domain.Page, error)
	FindPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	CreatePage(userID string, page *domain.PageInputData) (*domain.Page, error)
	DeletePage(id string) (domain.Page, error)
	UpdatePage(id string, data interface{}) (domain.Page, error)
	UpdatePageWithContent(id string, data map[string]interface{}) (domain.Page, error)
}

type Story interface {
	PublishStory(id string, data domain.StoryInputData) (domain.Story, error)
	GetStory(params domain.RequestParams) (domain.Story, error)
}

type Component interface {
	GetComponent(id string) (domain.Component, error)
	FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	CreateComponent(userID string, component *domain.ComponentInput) (*domain.Component, error)
	DeleteComponent(id string) (domain.Component, error)
	UpdateComponent(id string, data interface{}) (domain.Component, error)
	FindByPopulate(params domain.RequestParams) (domain.Response[domain.Component], error)
	FindGroupComponent(params domain.RequestParams) (domain.Response[domain.Component], error)

	FindLibrarys(params domain.RequestParams) (domain.Response[domain.Library], error)
}

type ComponentGroup interface {
	FindComponentGroup(params domain.RequestParams) (domain.Response[domain.ComponentGroup], error)
	CreateComponentGroup(userID string, component *domain.ComponentGroup) (*domain.ComponentGroup, error)
	UpdateComponentGroup(id string, data interface{}) (domain.ComponentGroup, error)
	DeleteComponentGroup(id string) (domain.ComponentGroup, error)
}

type ComponentPreset interface {
	FindComponentPreset(params domain.RequestParams) (domain.Response[domain.ComponentPreset], error)
	CreateComponentPreset(userID string, preset *domain.ComponentPresetInput) (*domain.ComponentPreset, error)
	UpdateComponentPreset(id string, data interface{}) (domain.ComponentPreset, error)
	DeleteComponentPreset(id string) (domain.ComponentPreset, error)
}

type Partner interface {
	CreatePartner(userID string, data *domain.PartnerInput) (domain.Partner, error)
	GetPartner(id string) (domain.Partner, error)
	FindPartner(params domain.RequestParams) (domain.Response[domain.PartnerPopulate], error)
	UpdatePartner(id string, data *domain.PartnerInput) (domain.Partner, error)
	DeletePartner(id string) (domain.Partner, error)
}

type Product interface {
	CreateProduct(userID string, data *domain.ProductInput) (domain.Product, error)
	GetProduct(id string) (domain.Product, error)
	FindProduct(params domain.RequestParams) (domain.Response[domain.Product], error)
	UpdateProduct(id string, data interface{}) (domain.Product, error)
	DeleteProduct(id string) (domain.Product, error)
}

type Image interface {
	CreateImage(userID string, data *domain.ImageInput) (domain.Image, error)
	GetImage(id string) (domain.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[domain.Image], error)
	DeleteImage(id string) (domain.Image, error)
}

type Apps interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Space interface {
	CreateSpace(userID string, space *domain.SpaceInput) (*domain.Space, error)
	GetSpace(id string) (domain.Space, error)
	FindSpace(params domain.RequestParams) (domain.Response[domain.Space], error)
	UpdateSpace(id string, data interface{}) (domain.Space, error)
	DeleteSpace(id string) (domain.Space, error)
}

type Plugin interface {
	CreatePlugin(userID string, plugin *domain.PluginInput) (*domain.Plugin, error)
	GetPlugin(id string) (domain.Plugin, error)
	FindPlugin(params domain.RequestParams) (domain.Response[domain.Plugin], error)
	UpdatePlugin(id string, data interface{}) (domain.Plugin, error)
	DeletePlugin(id string) (domain.Plugin, error)
}

type Services struct {
	Authorization
	Apps
	Component
	ComponentGroup
	ComponentPreset
	Image
	Page
	Partner
	Product
	Plugin
	Shop
	Space
	Story
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
		Shop:            NewShopService(cfgService.Repositories.Shop),
		Apps:            NewAppsService(cfgService.Repositories, cfgService.I18n),
		Component:       NewComponentService(cfgService.Repositories.Component, cfgService.I18n),
		ComponentGroup:  NewComponentGroupService(cfgService.Repositories.ComponentGroup, cfgService.I18n),
		ComponentPreset: NewComponentPresetService(cfgService.Repositories.ComponentPreset, cfgService.I18n),
		Image:           NewImageService(cfgService.Repositories.Image),
		Page:            NewPageService(cfgService.Repositories.Page, cfgService.I18n),
		Partner:         NewPartnerService(cfgService.Repositories.Partner, cfgService.I18n),
		Product:         NewProductService(cfgService.Repositories, cfgService.I18n),
		Plugin:          NewPluginService(cfgService.Repositories.Plugin, cfgService.I18n),
		Space:           NewSpaceService(cfgService.Repositories.Space, cfgService.I18n),
		Story:           NewStoryService(cfgService.Repositories.Story, cfgService.I18n),
		User:            NewUserService(cfgService.Repositories.User),
	}
}
