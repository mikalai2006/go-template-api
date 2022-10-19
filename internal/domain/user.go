package domain

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty" primitive:"true"`

	Type     string `json:"type" db:"type" bson:"type"`
	Name     string `json:"name,omitempty" db:"name" bson:"name"`
	Login    string `json:"login" db:"login" bson:"login"`
	Currency string `json:"currency" bson:"currency" db:"currency"`
	Lang     string `json:"lang" db:"lang" bson:"lang"`
	Avatar   string `json:"avatar" db:"avatar" bson:"avatar"`

	Roles     []string  `json:"roles" bson:"-" db:"-" form:"-"`
	Online    bool      `json:"online" db:"online" bson:"online"`
	Verify    bool      `json:"verify" db:"verify" bson:"verify"`
	LastTime  time.Time `json:"last_time" db:"last_time" bson:"last_time"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bson:"updated_at"`
}

type UserInput struct {
	Type     string `json:"type" db:"type" bson:"type" form:"type"`
	Name     string `json:"name,omitempty" db:"name" bson:"name" form:"name"`
	Login    string `json:"login" db:"login" bson:"login" form:"login"`
	Currency string `json:"currency" bson:"currency" db:"currency" form:"currency"`
	Lang     string `json:"lang" db:"lang" bson:"lang" form:"lang"`
	Avatar   string `json:"avatar" db:"avatar" bson:"avatar" form:"avatar"`
}

func (user *User) BodyToData() (interface{}, error) {
	result := bson.M{}
	var tagValue string
	elementsFilter := reflect.ValueOf(User{})

	for i := 0; i < elementsFilter.NumField(); i++ {
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
	// if user.Online != nil {
	// 	result["online"] = user.Online
	// }
	// if user.Verify != nil {
	// 	result["verify"] = user.Verify
	// }

	result["updated_at"] = time.Now()
	return result, nil
}
