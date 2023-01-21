package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userId" bson:"user_id"`
	Service     string             `json:"service" bson:"service"`
	ServiceID   string             `json:"serviceId" bson:"service_id"`
	Path        string             `json:"path" bson:"path"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updated_at"`
}

type ImageInput struct {
	UserID      string `json:"userId" bson:"user_id"`
	Service     string `json:"service" bson:"service" form:"service" binding:"required"`
	ServiceID   string `json:"serviceId,omitempty" bson:"service_id" form:"serviceId"`
	Path        string `json:"path" bson:"path"`
	Description string `json:"description" bson:"description" form:"description"`
	Title       string `json:"title" bson:"title" form:"title"`
	// Images      *multipart.FileHeader `bson:"image" form:"image"`
}