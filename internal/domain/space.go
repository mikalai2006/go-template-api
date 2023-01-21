package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Space struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID     `json:"userId" bson:"user_id"`
	Title       string                 `json:"title" bson:"title"`
	Description string                 `json:"description" bson:"description"`
	Setting     map[string]interface{} `json:"setting" bson:"setting"`
	CreatedAt   time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updated_at"`
}

type SpaceInput struct {
	UserID      string                 `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	Title       string                 `json:"title" bson:"title" form:"title"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Setting     map[string]interface{} `json:"setting" bson:"setting" form:"setting"`
	CreatedAt   time.Time              `json:"createdAt" bson:"created_at" form:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updated_at" form:"updatedAt"`
}
