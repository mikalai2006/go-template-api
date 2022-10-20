package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty" form:"-"`
	UserID primitive.ObjectID `json:"userId" bson:"userId" primitive:"true"`

	ComponentID primitive.ObjectID `json:"componentId" bson:"componentId" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layoutId" bson:"layoutId" primitive:"true"`
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slugFull" bson:"slug_full"`
	Title       string             `json:"title" bson:"title"`
	Path        string             `json:"path" bson:"path"`
	Name        string             `json:"name" bson:"name"`
	ContentType string             `json:"contentType" bson:"content_type"`
	Publish     bool               `json:"publish" bson:"publish"`
	SortOrder   int                `json:"sortOrder" bson:"sort_order"`
	Setting     interface{}        `json:"setting" bson:"setting"`

	Component     Component       `json:"-" bson:"component"`
	Layout        Component       `json:"-" bson:"layout"`
	XXX           any             `json:"-" bson:"xxx"`
	Content       any             `json:"content" bson:"content"`
	ComponentData []ComponentData `json:"-" bson:"component_data"`
	CreatedAt     time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt" bson:"updatedAt"`
}

// type PageContent struct {
// }

type PageInputData struct {
	Slug        string             `json:"slug" bson:"slug" form:"slug"`
	SlugFull    string             `json:"slugFull" bson:"slug_full" form:"slugFull"`
	Title       string             `json:"title" bson:"title" form:"title"`
	Path        string             `json:"path" bson:"path" form:"path"`
	Name        string             `json:"name" bson:"name" form:"name"`
	ContentType string             `json:"contentType" bson:"content_type" form:"contentType"`
	Publish     bool               `json:"publish" bson:"publish" form:"publish"`
	SortOrder   int                `json:"sortOrder" bson:"sort_order" form:"sort_order"`
	Setting     interface{}        `json:"setting" bson:"setting" form:"setting"`
	ComponentID primitive.ObjectID `json:"componentId" bson:"componentId" form:"componentId" primitive:"true"`
	LayoutID    primitive.ObjectID `json:"layoutId" bson:"layoutId" form:"layoutId" primitive:"true"`
}

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
