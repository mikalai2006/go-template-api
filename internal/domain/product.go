package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty" primitive:"true"`
	UserID     primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	ShopID     primitive.ObjectID `json:"shopId" bson:"shop_id" primitive:"true"`
	CategoryID primitive.ObjectID `json:"categoryId" bson:"category_id" primitive:"true"`
	SeoID      int64              `json:"seoId" bson:"seo_id"`

	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`

	Seo       string    `json:"seo" bson:"seo"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ProductInput struct {
	Locale map[string]map[string]string `json:"locale" bson:"locale" form:"locale"`

	Title       string `json:"title" bson:"title" form:"title"`
	Description string `json:"description" bson:"description" form:"description"`
	ShopID      string `json:"shopId" bson:"shop_id" primitive:"true"`
	CategoryID  string `json:"categoryId" bson:"category_id" primitive:"true"`
	Seo         string `json:"seo" bson:"seo" form:"seo"`
}
