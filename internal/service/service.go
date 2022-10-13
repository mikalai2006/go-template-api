package service

import (
	"time"

	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/auths"
	"github.com/mikalai2006/go-template-api/pkg/hasher"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Authorization interface {
	CreateAuth(auth domain.SignInInput) (primitive.ObjectID, error)
	SignIn(input domain.SignInInput) (domain.ResponseTokens, error)
	ExistAuth(auth domain.SignInInput) (domain.Auth, error)
	CreateSession(auth domain.Auth) (domain.ResponseTokens, error)
	VerificationCode(userId string, code string) error
	RefreshTokens(refreshToken string) (domain.ResponseTokens, error)
}

type Shop interface {
	Find(params domain.RequestParams) (domain.Response[domain.Shop], error)

	GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userId string, user domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user domain.User) (domain.User, error)
}

type Page interface {
	GetPage(id string) (domain.Page, error)
	FindPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	CreatePage(userId string, page domain.Page) (*domain.Page, error)
	DeletePage(id string) (domain.Page, error)
	UpdatePage(id string, user domain.Page) (domain.Page, error)
}

type Component interface {
	GetComponent(id string) (domain.Component, error)
	FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	CreateComponent(userId string, component domain.Component) (*domain.Component, error)
	DeleteComponent(id string) (domain.Component, error)
	UpdateComponent(id string, user domain.Component) (domain.Component, error)
}

type Services struct {
	Authorization
	Shop
	User
	Page
	Component
}

type ConfigServices struct {
	Repositories *repository.Repositories
	Hasher hasher.PasswordHasher
	TokenManager auths.TokenManager
	OtpGenerator utils.Generator
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
	VerificationCodeLength int
}

func NewServices(cfgService *ConfigServices) *Services {
	return &Services{
		Authorization: NewAuthService(cfgService.Repositories.Authorization, cfgService.Hasher, cfgService.TokenManager, cfgService.RefreshTokenTTL, cfgService.AccessTokenTTL, cfgService.OtpGenerator, cfgService.VerificationCodeLength),
		Shop: NewShopService(cfgService.Repositories.Shop),
		User: NewUserService(cfgService.Repositories.User),
		Page: NewPageService(cfgService.Repositories.Page),
		Component: NewComponentService(cfgService.Repositories.Component),
	}
}