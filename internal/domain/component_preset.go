package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComponentPreset struct {
	ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty" primitive:"true"`
	UserID      primitive.ObjectID     `json:"userId" bson:"user_id" primitive:"true"`
	ComponentID primitive.ObjectID     `json:"componentId" bson:"component_id" primitive:"true"`
	Title       string                 `json:"title" bson:"title"`
	Description string                 `json:"description" bson:"description"`
	Data        map[string]interface{} `json:"data" bson:"data"`
	Image       map[string]interface{} `json:"image" bson:"image"`
	CreatedAt   time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updated_at"`
}

type ComponentPresetInput struct {
	ComponentID primitive.ObjectID     `json:"componentId" bson:"component_id" form:"componentId" primitive:"true"`
	Title       string                 `json:"title" bson:"title" form:"title"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Data        map[string]interface{} `json:"data" bson:"data" form:"data"`
	Image       map[string]interface{} `json:"image" bson:"image" form:"image"`
}
type ComponentPresetFind struct {
	ComponentID string                 `json:"componentId" bson:"component_id" form:"componentId" primitive:"true"`
	Title       string                 `json:"title" bson:"title" form:"title"`
	Description string                 `json:"description" bson:"description" form:"description"`
	Data        map[string]interface{} `json:"data" bson:"data" form:"data"`
}
