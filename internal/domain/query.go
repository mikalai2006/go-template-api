package domain

type RequestParams struct {
	Options
	Filter interface{} `json:"filter" bson:"filter"`
	Group  interface{} `json:"group" bson:"$group"`
}

type Options struct {
	Limit int64       `json:"$limit" bson:"limit" form:"$limit"`
	Skip  int64       `json:"$skip" bson:"skip" form:"$skip"`
	Sort  interface{} `json:"$sort" bson:"sort" form:"$sort"`
}

type RefreshInput struct {
	Token string `json:"token" bson:"token" form:"token"`
}
