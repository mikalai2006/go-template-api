package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Language struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Publish   bool               `json:"publish" bson:"publish"`
	Flag      string             `json:"flag" bson:"flag"`
	Name      string             `json:"name" bson:"name"`
	Code      string             `json:"code" bson:"code"`
	Locale    string             `json:"locale" bson:"locale"`
	SortOrder int64              `json:"sort_order" bson:"sort_order"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type LanguageInput struct {
	Publish   bool   `json:"publish" bson:"publish" form:"publish"`
	Flag      string `json:"flag" bson:"flag" form:"flag"`
	Name      string `json:"name" bson:"name" form:"name"`
	Code      string `json:"code" bson:"code" form:"code"`
	Locale    string `json:"locale" bson:"locale" form:"locale"`
	SortOrder int64  `json:"sort_order" bson:"sort_order" form:"sort_order"`
}

type Category struct {
	ID          int64             `json:"id" bson:"_id,omitempty" form:"-"`
	ParentID    int64             `json:"parent_id" bson:"parent_id"`
	Title       map[string]string `json:"title" bson:"title" form:"title"`
	Description map[string]string `json:"description" bson:"description" form:"description"`
	Seo         string            `json:"seo" bson:"seo" form:"seo"`
	SortOrder   int64             `json:"sort_order" bson:"sort_order" form:"sort_order"`
	MPath       string            `json:"mpath" bson:"mpath" form:"mpath"`
	Level       string            `json:"level" bson:"level" form:"level"`
	Status      bool              `json:"status" bson:"status" form:"status"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type CategoryInput struct {
	Title       map[string]string `json:"title" bson:"title" form:"title"`
	Description map[string]string `json:"description" bson:"description" form:"description"`
	Seo         string            `json:"seo" bson:"seo" form:"seo"`
	SortOrder   int64             `json:"sort_order" bson:"sort_order" form:"sort_order"`
	MPath       string            `json:"mpath" bson:"mpath" form:"mpath"`
	Level       string            `json:"level" bson:"level" form:"level"`
	Status      bool              `json:"status" bson:"status" form:"status"`
}
