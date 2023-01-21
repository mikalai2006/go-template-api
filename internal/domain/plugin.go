package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plugin struct {
	ID          primitive.ObjectID       `json:"id" bson:"_id,omitempty"`
	SpaceID     primitive.ObjectID       `json:"spaceId" bson:"space_id"  primitive:"true"`
	UserID      primitive.ObjectID       `json:"userId" bson:"user_id"`
	Name        string                   `json:"name" bson:"name"`
	Description string                   `json:"description" bson:"description"`
	Body        string                   `json:"body" bson:"body"`
	Code        string                   `json:"code" bson:"code"`
	Options     []map[string]interface{} `json:"options" bson:"options"`
	CreatedAt   time.Time                `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                `json:"updatedAt" bson:"updated_at"`
}

type PluginInput struct {
	UserID      string                   `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	SpaceID     string                   `json:"spaceId" bson:"space_id" form:"spaceId" primitive:"true"`
	Name        string                   `json:"name" bson:"name" form:"name"`
	Description string                   `json:"description" bson:"description" form:"description"`
	Body        string                   `json:"body" bson:"body" form:"body"`
	Code        string                   `json:"code" bson:"code" form:"code"`
	Options     []map[string]interface{} `json:"options" bson:"options" form:"options"`
}
