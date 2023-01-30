package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Partner struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`

	SeoID       int64                        `json:"seoId" bson:"seo_id"`
	Seo         string                       `json:"seo" bson:"seo"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type PartnerInput struct {
	UserID string `json:"userId" bson:"user_id" form:"userId" primitive:"true"`

	SeoID       int64                        `json:"seoId" bson:"seo_id" form:"seoId"`
	Seo         string                       `json:"seo" bson:"seo" form:"seo"`
	Title       string                       `json:"title" bson:"title" form:"title"`
	Description string                       `json:"description" bson:"description" form:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale" form:"locale"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at" form:"updatedAt"`
}

type PartnerPopulate struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty" primitive:"true"`
	UserID      primitive.ObjectID           `json:"userId" bson:"user_id" primitive:"true"`
	SeoID       int64                        `json:"seoId" bson:"seo_id"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Seo         string                       `json:"seo" bson:"seo"`

	Images []Image `json:"images" bson:"images"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
