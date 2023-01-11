package domain

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Type     string `json:"type" db:"type" bson:"type"`

type User struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"user_id,omitempty" primitive:"true"`

	Name      string    `json:"name,omitempty" bson:"name"`
	Login     string    `json:"login" bson:"login"`
	Currency  string    `json:"currency" bson:"currency"`
	Lang      string    `json:"lang" bson:"lang"`
	Avatar    string    `json:"avatar" bson:"avatar"`
	Roles     []string  `json:"roles" bson:"-"`
	Online    bool      `json:"online" bson:"online"`
	Verify    bool      `json:"verify" bson:"verify"`
	LastTime  time.Time `json:"lastTime" bson:"last_time"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type UserInput struct {
	Name     string `json:"name,omitempty" bson:"name" form:"name"`
	Login    string `json:"login" bson:"login" form:"login"`
	Currency string `json:"currency" bson:"currency" form:"currency"`
	Lang     string `json:"lang" bson:"lang" form:"lang"`
	Avatar   string `json:"avatar" bson:"avatar" form:"avatar"`
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
