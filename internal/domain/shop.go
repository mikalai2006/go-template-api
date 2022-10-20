package domain

import "time"

type Shop struct {
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Seo         string    `json:"seo" bson:"seo"`
	UserID      string    `json:"user_id" bson:"user_id"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type ShopInput struct {
	RequestParams
	Title       string `json:"title" bson:"title" form:"title"`
	Description string `json:"description" bson:"description" form:"description"`
	Seo         string `json:"seo" bson:"seo" form:"seo"`
}
