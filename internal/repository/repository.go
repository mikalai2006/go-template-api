package repository

import (
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateAuth(auth domain.Auth) (primitive.ObjectID, error)
	GetAuth(auth domain.Auth) (domain.Auth, error)
	CheckExistAuth(auth domain.SignInInput) (domain.Auth, error)
	GetByCredentials(auth domain.SignInInput) (domain.Auth, error)
	SetSession(authId primitive.ObjectID, session domain.Session) error
	VerificationCode(userId string, code string) error
	RefreshToken(refreshToken string) (domain.Auth, error)
}

type Shop interface {
	Find(params domain.RequestParams) (domain.Response[domain.Shop], error)
	GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}

type Page interface {
	FindPage(params domain.RequestParams) (domain.Response[domain.Page], error)
	GetPage(id string) (domain.Page, error)
	CreatePage(userId string, page domain.Page) (*domain.Page, error)
	DeletePage(id string) (domain.Page, error)
	UpdatePage(id string, page domain.Page) (domain.Page, error)
}

type Component interface {
	GetComponent(id string) (domain.Component, error)
	FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	CreateComponent(userId string, component domain.Component) (*domain.Component, error)
	DeleteComponent(id string) (domain.Component, error)
	UpdateComponent(id string, component domain.Component) (domain.Component, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userId string, user domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user domain.User) (domain.User, error)
}

type Repositories struct {
	Authorization
	Shop
	User
	Page
	Component
}

func NewRepositories(mongodb *mongo.Database) *Repositories {
	return &Repositories{
		Authorization: NewAuthMongo(mongodb),
		Shop: NewShopMongo(mongodb),
		User: NewUserMongo(mongodb),
		Page: NewPageMongo(mongodb),
		Component: NewComponentMongo(mongodb),
	}
}