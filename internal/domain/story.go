package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Story struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID   primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	SpaceID  primitive.ObjectID `json:"spaceId" bson:"space_id" primitive:"true"`
	PageID   primitive.ObjectID `json:"pageId" bson:"page_id" primitive:"true"`
	LayoutID primitive.ObjectID `json:"layoutId" bson:"layout_id" primitive:"true"`

	Name     string                 `json:"name" bson:"name"`
	Title    string                 `json:"title" bson:"title"`
	Slug     string                 `json:"slug" bson:"slug"`
	SlugFull string                 `json:"slugFull" bson:"slug_full"`
	Content  map[string]interface{} `json:"content" bson:"content"`

	Publish   bool      `json:"publish" bson:"publish"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type StoryInputData struct {
	SpaceID  string `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	UserID   string `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	PageID   string `json:"pageId" bson:"page_id" form:"pageId" primitive:"true"`
	LayoutID string `json:"layoutId" bson:"layout_id" form:"layoutId" primitive:"true"`

	Slug     string                 `json:"slug" bson:"slug" form:"slug"`
	SlugFull string                 `json:"slugFull" bson:"slug_full" form:"slugFull"`
	Title    string                 `json:"title" bson:"title" form:"title"`
	Name     string                 `json:"name" bson:"name" form:"name"`
	Content  map[string]interface{} `json:"content" bson:"content" form:"content"`

	Publish   bool      `json:"publish" bson:"publish" form:"publish"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
