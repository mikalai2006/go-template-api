package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`

	SpaceID     primitive.ObjectID `json:"spaceId" bson:"space_id" primitive:"true"`
	ComponentID primitive.ObjectID `json:"componentId" bson:"component_id" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layoutId" bson:"layout_id" primitive:"true"`
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slugFull" bson:"slug_full"`
	Title       string             `json:"title" bson:"title"`
	Path        string             `json:"path" bson:"path"`
	Name        string             `json:"name" bson:"name"`
	Publish     bool               `json:"publish" bson:"publish"`
	SortOrder   int                `json:"sortOrder" bson:"sort_order"`
	ContentType string             `json:"contentType" bson:"content_type"`
	Setting     interface{}        `json:"setting" bson:"setting"`

	Component     Component       `json:"component" bson:"component"`
	Layout        Component       `json:"layout" bson:"layout"`
	XXX           any             `json:"-" bson:"xxx"`
	Content       interface{}     `json:"content" bson:"content"`
	ComponentData []ComponentData `json:"-" bson:"component_data"`
	CreatedAt     time.Time       `json:"createdAt" bson:"created_at"`
	UpdatedAt     time.Time       `json:"updatedAt" bson:"updated_at"`
}

type PageWithContent struct {
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	SpaceID     primitive.ObjectID `json:"spaceId" bson:"space_id" primitive:"true"`
	ComponentID primitive.ObjectID `json:"componentId" bson:"component_id" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layoutId" bson:"layout_id" primitive:"true"`
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slugFull" bson:"slug_full"`
	Title       string             `json:"title" bson:"title"`
	Path        string             `json:"path" bson:"path"`
	Name        string             `json:"name" bson:"name"`
	Publish     bool               `json:"publish" bson:"publish"`
	SortOrder   int                `json:"sortOrder" bson:"sort_order"`
	Setting     interface{}        `json:"setting" bson:"setting"`
	ContentType string             `json:"contentType" bson:"content_type"`

	Content   any       `json:"content" bson:"content" form:"content"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

// type PageContent struct {
// }

type PageInputData struct {
	SpaceID     primitive.ObjectID `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	UserID      primitive.ObjectID `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slugFull" bson:"slug_full" form:"slugFull"`
	Title       string             `json:"title" bson:"title" form:"title"`
	Path        string             `json:"path" bson:"path" form:"path"`
	Name        string             `json:"name" bson:"name" form:"name"`
	Publish     bool               `json:"publish" bson:"publish" form:"publish"`
	SortOrder   int                `json:"sortOrder" bson:"sort_order" form:"sort_order"`
	Setting     interface{}        `json:"setting" bson:"setting" form:"setting"`
	ComponentID primitive.ObjectID `json:"componentId" bson:"component_id" form:"componentId" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layoutId" bson:"layout_id" form:"layoutId" primitive:"true"`
	ContentType string             `json:"contentType" bson:"content_type" form:"contentType"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type PageFilterData struct {
	SpaceID     string      `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	ID          string      `json:"id" bson:"_id" form:"id" primitive:"true"`
	UserID      string      `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	Slug        string      `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string      `json:"slugFull" bson:"slug_full" form:"slugFull"`
	Title       string      `json:"title" bson:"title" form:"title"`
	Path        string      `json:"path" bson:"path" form:"path"`
	Name        string      `json:"name" bson:"name" form:"name"`
	Publish     bool        `json:"publish" bson:"publish" form:"publish"`
	SortOrder   int         `json:"sortOrder" bson:"sort_order" form:"sort_order"`
	Setting     interface{} `json:"setting" bson:"setting" form:"setting"`
	ComponentID string      `json:"componentId" bson:"component_id" form:"componentId" primitive:"true"`
	LayoutID    string      `json:"layoutId" bson:"layout_id" form:"layoutId" primitive:"true"`
	ContentType string      `json:"contentType" bson:"content_type" form:"contentType"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

// type PageQuery struct {
// 	PaginationQuery
// 	Options
// 	PageFilter
// }

// type PageFilter struct {
// 	Slug        string      `json:"slug" bson:"slug" form:"slug"`
// 	SlugFull    string      `json:"slugFull" bson:"slug_full" form:"slugFull"`
// 	Title       string      `json:"title" bson:"title" form:"title"`
// 	Path        string      `json:"path" bson:"path" form:"path"`
// 	Name        string      `json:"name" bson:"name" form:"name"`
// 	ContentType string      `json:"contentType" bson:"content_type" form:"contentType"`
// 	Publish     bool        `json:"publish" bson:"publish" form:"publish"`
// 	SortOrder   int         `json:"sortOrder" bson:"sort_order" form:"sort_order"`
// 	Setting     interface{} `json:"setting" bson:"setting" form:"setting"`
// 	ComponentID string      `json:"componentId" bson:"componentId" form:"componentId" primitive:"true"`
// 	LayoutID    string      `json:"layoutId" bson:"layoutId" form:"layoutId" primitive:"true"`
// }

type PageRoutes struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Slug     string             `json:"slug" bson:"slug"`
	Path     string             `json:"path" bson:"path"`
	SlugFull string             `json:"slugFull" bson:"slug_full"`
	Publish  bool               `json:"publish" bson:"publish"`
}

// func (page *Page) BodyToData() (interface{}, error) {
// 	result := bson.M{}
// 	var tagValue string
// 	elementsFilter := reflect.ValueOf(Page{})

// 	for i := 0; i < elementsFilter.NumField(); i++ {
// 		typeField := elementsFilter.Type().Field(i)
// 		tag := typeField.Tag

// 		tagValue = tag.Get("bson")

// 		if tagValue == "-" {
// 			continue
// 		}

// 		if elementsFilter.Field(i).Interface() == "" {
// 			continue
// 		}

// 		switch elementsFilter.Field(i).Kind() {
// 		case reflect.String:
// 			value := elementsFilter.Field(i).String()
// 			result[tagValue] = value

// 		case reflect.Bool:
// 			value := elementsFilter.Field(i).Bool()
// 			result[tagValue] = value

// 		case reflect.Int:
// 			value := elementsFilter.Field(i).Int()
// 			result[tagValue] = value
// 		}
// 	}
// 	// fmt.Println("result ", result)
// 	result["updated_at"] = time.Now()
// 	return result, nil
// }
