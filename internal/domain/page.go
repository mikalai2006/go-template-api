package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty" form:"-"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id" primitive:"true"`

	ComponentID primitive.ObjectID `json:"component_id" bson:"component_id" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layout_id" bson:"layout_id" primitive:"true"`
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slug_full" bson:"slug_full" form:"slug_full"`
	Title       string             `json:"title" bson:"title" form:"title"`
	Path        string             `json:"path" bson:"path" form:"path"`
	Name        string             `json:"name" bson:"name" form:"name"`
	Publish     bool               `json:"publish" bson:"publish" form:"publish"`
	SortOrder   int                `json:"sort_order" bson:"sort_order" form:"sort_order"`
	Setting     interface{}        `json:"setting" bson:"setting" form:"setting"`

	Component     Component       `json:"component" bson:"component"`
	Layout        Component       `json:"layout" bson:"layout"`
	XXX           any             `json:"xxx" bson:"xxx"`
	Content       any             `json:"content" bson:"content"`
	ComponentData []ComponentData `json:"component_data" bson:"component_data"`
	CreatedAt     time.Time       `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

// type PageContent struct {
// }

type PageInputData struct {
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slug_full" bson:"slug_full" form:"slug_full"`
	Title       string             `json:"title" bson:"title" form:"title"`
	Path        string             `json:"path" bson:"path" form:"path"`
	Name        string             `json:"name" bson:"name" form:"name"`
	Publish     bool               `json:"publish" bson:"publish" form:"publish"`
	SortOrder   int                `json:"sort_order" bson:"sort_order" form:"sort_order"`
	Setting     interface{}        `json:"setting" bson:"setting" form:"setting"`
	ComponentID primitive.ObjectID `json:"component_id" bson:"component_id" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layout_id" bson:"layout_id" primitive:"true"`
}

type PageRoutes struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Slug     string             `json:"slug" bson:"slug"`
	Path     string             `json:"path" bson:"path"`
	SlugFull string             `json:"slug_full" bson:"slug_full"`
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
