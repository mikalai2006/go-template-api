package domain

type ResponseTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}