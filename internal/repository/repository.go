package repository

import (
	"reflect"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
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
	CreatePage(userID string, page *domain.PageInputData) (*domain.Page, error)
	DeletePage(id string) (domain.Page, error)
	UpdatePage(id string, data interface{}) (domain.Page, error)
	UpdatePageWithContent(id string, data map[string]interface{}) (domain.Page, error)
}

type Component interface {
	GetComponent(id string) (domain.Component, error)
	FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	FindGroupComponent(params domain.RequestParams) (domain.Response[domain.Component], error)
	CreateComponent(userID string, component *domain.ComponentInput) (*domain.Component, error)
	DeleteComponent(id string) (domain.Component, error)
	UpdateComponent(id string, data interface{}) (domain.Component, error)

	FindByPopulate(params domain.RequestParams) (domain.Response[domain.Component], error)
	FindLibrarys(params domain.RequestParams) (domain.Response[domain.Library], error)
}

type ComponentGroup interface {
	FindComponentGroup() (domain.Response[domain.ComponentGroup], error)
	CreateComponentGroup(userID string, componentGroup *domain.ComponentGroup) (*domain.ComponentGroup, error)
	UpdateComponentGroup(id string, data interface{}) (domain.ComponentGroup, error)
	DeleteComponentGroup(id string) (domain.ComponentGroup, error)
}

type ComponentPreset interface {
	FindComponentPreset(params domain.RequestParams) (domain.Response[domain.ComponentPreset], error)
	CreateComponentPreset(userID string, ComponentPreset *domain.ComponentPresetInput) (*domain.ComponentPreset, error)
	UpdateComponentPreset(id string, data interface{}) (domain.ComponentPreset, error)
	DeleteComponentPreset(id string) (domain.ComponentPreset, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userID string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user *domain.User) (domain.User, error)
	Iam(userID string) (domain.User, error)
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

type Repositories struct {
	Authorization
	Apps
	Component
	ComponentGroup
	ComponentPreset
	Image
	Page
	Product
	Shop
	User
}

func NewRepositories(mongodb *mongo.Database, i18n config.I18nConfig) *Repositories {
	return &Repositories{
		Authorization:   NewAuthMongo(mongodb),
		Apps:            NewAppsMongo(mongodb, i18n),
		Component:       NewComponentMongo(mongodb, i18n),
		ComponentGroup:  NewComponentGroupMongo(mongodb, i18n),
		ComponentPreset: NewComponentPresetMongo(mongodb, i18n),
		Image:           NewImageMongo(mongodb, i18n),
		Page:            NewPageMongo(mongodb, i18n),
		Product:         NewProductMongo(mongodb, i18n),
		Shop:            NewShopMongo(mongodb, i18n),
		User:            NewUserMongo(mongodb, i18n),
	}
}

// func getPaginationOpts(pagination *domain.PaginationQuery) *options.FindOptions {
// 	var opts *options.FindOptions
// 	if pagination != nil {
// 		opts = &options.FindOptions{
// 			Skip:  pagination.GetSkip(),
// 			Limit: pagination.GetLimit(),
// 		}
// 	}

// 	return opts
// }

func createFilter[V any](filterData *V) any {
	var filter V

	filterReflect := reflect.ValueOf(filterData)
	// fmt.Println("========== filterReflect ===========")
	// fmt.Println("struct > ", filterReflect)
	// fmt.Println("struct type > ", filterReflect.Type())
	filterIndirectData := reflect.Indirect(filterReflect)
	// fmt.Println("filter data > ", filterIndirectData)
	// fmt.Println("filter numField > ", filterIndirectData.NumField())
	dataFilter := bson.M{}

	var tagJSON, tagPrimitive string
	for i := 0; i < filterIndirectData.NumField(); i++ {
		field := filterIndirectData.Field(i)
		if field.Kind() == reflect.Ptr {
			field = reflect.Indirect(field)
		}
		typeField := filterIndirectData.Type().Field(i)
		tag := typeField.Tag
		// tagBson = tag.Get("bson")
		tagJSON = tag.Get("json")
		tagPrimitive = tag.Get("primitive")
		switch field.Kind() {
		case reflect.String:
			value := field.String()
			if tagPrimitive == "true" {
				id, _ := primitive.ObjectIDFromHex(value)
				// fmt.Println("===== string add ", tag, value)
				dataFilter[tagJSON] = id
			} else {
				dataFilter[tagJSON] = value
			}

		case reflect.Bool:
			value := field.Bool()
			dataFilter[tagJSON] = value

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := field.Int()
			dataFilter[tagJSON] = value

		default:

		}

		// fmt.Println(tagBson, tagJSON, tagPrimitive, fmt.Sprintf("[%s]", field), field.Kind(), field)
	}

	// structure := reflect.ValueOf(&filter)
	// fmt.Println("========== filter ===========")
	// fmt.Println("struct > ", structure)
	// fmt.Println("struct type > ", structure.Type())
	// fmt.Println("filter data > ", reflect.Indirect(structure))
	// fmt.Println("filter numField > ", reflect.Indirect(structure).NumField())

	// fmt.Println("========== result ===========")
	// fmt.Println("dataFilter > ", dataFilter)
	return filter
}
