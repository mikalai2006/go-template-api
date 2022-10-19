package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	// swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login        string             `json:"login" binding:"required"`
	Email        string             `json:"email"`
	Password     string             `json:"password" binding:"required"`
	Strategy     string             `json:"strategy"`
	VkID         string             `json:"vk_id"`
	GoogleID     string             `json:"google_id" bson:"google_id"`
	GithubID     string             `json:"github_id"`
	AppleID      string             `json:"apple_id"`
	Verification Verification       `json:"verification" bson:"verification"`
	Session      Session            `json:"session" bson:"session"`
	Roles        []string           `json:"roles" bson:"roles"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}
type SignInInput struct {
	Login    string `json:"login" bson:"login"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Strategy string `json:"strategy" bson:"strategy"`
	VkID     string `json:"-"`
	GoogleID string `json:"-"`
}

type DataForClaims struct {
	Roles  []string `json:"roles"`
	UserID string   `json:"user_id"`
}
