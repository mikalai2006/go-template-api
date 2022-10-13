package domain

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	Id primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	ComponentId primitive.ObjectID `json:"component_id" bson:"component_id"`
	Slug string `json:"slug" bson:"slug" form:"slug"`
	SlugFull string `json:"slug_full" bson:"slug_full" form:"slug_full"`

	Title       string `json:"title" bson:"title" form:"title"`
	Path       string `json:"path" bson:"path" form:"path"`
	Name       string `json:"name" bson:"name" form:"name"`
	Publish *bool `json:"publish" bson:"publish" form:"publish"`

	SortOrder int `json:"sort_order" bson:"sort_order" form:"sort_order"`
	Setting         interface{} `json:"setting" bson:"setting" form:"setting"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

func (page Page) BodyToData() (interface{}, error) {

	result := bson.M{}
	var tagValue string
	elementsFilter := reflect.ValueOf(page)

	for i := 0; i < elementsFilter.NumField(); i += 1 {
		typeField := elementsFilter.Type().Field(i)
		tag := typeField.Tag

		tagValue = tag.Get("bson")

		if tagValue == "-" {
				continue
		}

		if elementsFilter.Field(i).Interface() == "" {
			continue
		}

		switch elementsFilter.Field(i).Kind() {
		case reflect.String:
				value := elementsFilter.Field(i).String()
				result[tagValue] = value

		case reflect.Bool:
				value := elementsFilter.Field(i).Bool()
				result[tagValue] = value

		case reflect.Int:
				value := elementsFilter.Field(i).Int()
				result[tagValue] = value

		}
	}
	if page.Publish != nil {
		result["online"] = *page.Publish
	}

	result["updated_at"] = time.Now()

	fmt.Println("user: new data =", result)

	return result, nil
}