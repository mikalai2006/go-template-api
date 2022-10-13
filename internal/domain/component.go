package domain

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Component struct {
	Id primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name" form:"name"`
	Title       string `json:"title" bson:"title" form:"title"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Publish *bool `json:"publish" bson:"publish" form:"publish"`
	Status *bool `json:"status" bson:"status" form:"status"`
	IsPage *bool `json:"is_page" bson:"is_page" form:"is_page"`
	IsGlobal *bool `json:"is_global" bson:"is_global" form:"is_global"`
	IsLayout *bool `json:"is_layout" bson:"is_layout" form:"is_layout"`
	Tpl string  `json:"tpl" bson:"tpl" form:"tpl"`
	SortOrder int `json:"sort_order" bson:"sort_order" form:"sort_order"`
	Group []primitive.ObjectID `json:"group" bson:"group"`
	Setting         interface{} `json:"setting" bson:"setting" form:"setting"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

func (component Component) BodyToData() (interface{}, error) {

	result := bson.M{}
	var tagValue string
	elementsFilter := reflect.ValueOf(component)

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
	if len(component.Group) == 0 || component.Group == nil {
		result["group"] = []primitive.ObjectID{}
		// for i := 0; i < len(component.Group) - 1; i += 1 {
		// 	 id, err := primitive.ObjectIDFromHex(component.Group[i])
		// 	 if err != nil {
		// 		return result, nil
		// 	 }
		// 	 result["group"][i] = id
		// }
	} else {
		result["group"] = component.Group
	}
	if component.IsPage != nil {
		result["is_page"] = *component.IsPage
	}
	if component.IsGlobal != nil {
		result["is_global"] = *component.IsGlobal
	}
	if component.IsLayout != nil {
		result["is_layout"] = *component.IsLayout
	}
	if component.Status != nil {
		result["status"] = *component.Status
	}
	if component.Publish != nil {
		result["publish"] = *component.Publish
	}

	result["updated_at"] = time.Now()

	fmt.Println("user: new data =", result)

	return result, nil
}