package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComponentGroup struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" primitive:"true"`
	SpaceID     primitive.ObjectID `json:"spaceId" bson:"space_id" primitive:"true"`
	UserID      primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type ComponentGroupInput struct {
	UserID      string `json:"userId" bson:"user_id" primitive:"true"`
	SpaceID     string `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}
