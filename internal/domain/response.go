package domain

import "time"

type Response[D any] struct {
	Total int `json:"total" bson:"total"`
	Limit int `json:"limit" bson:"limit"`
	Skip  int `json:"skip" bson:"skip"`
	Data  []D `json:"data" bson:"data"`
}

type ResponseTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type GeneralFieldDB struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
