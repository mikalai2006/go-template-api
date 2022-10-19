package repository

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	tblShops = "shops"
	tblUsers = "users"
	TblAuth  = "auth"
	tblPage  = "pages"

	tblComponent           = "components"
	tblComponentData       = "component_datas"
	tblComponentGroup      = "component_groups"
	tblComponentSchema     = "component_schemas"
	tblComponentSchemaData = "component_schemadatas"

	tblLibrary = "librarys"
	tblFields  = "fields"

	TblLanguage = "langs"

	MongoQueryTimeout = 10 * time.Second
)

type ConfigMongoDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSL      bool
}

const timeDeadline = 30

func NewMongoDB(cfg *ConfigMongoDB) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeDeadline*time.Second)
	defer cancel()

	host := net.JoinHostPort(cfg.Host, cfg.Port)
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin&readPreference=primary&directConnection=true&ssl=%t",
		cfg.Username, cfg.Password, host, cfg.DBName, cfg.SSL)
	logger.Info(uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if er := client.Ping(ctx, readpref.Primary()); er != nil {
		return nil, er
	}

	return client, nil
}

func CreatePipeline(params domain.RequestParams) (mongo.Pipeline, error) {
	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{Key: "$match", Value: params.Filter}})
	// opts := options.Find()
	if params.Options.Sort != nil {
		// opts.SetSort(params.Options.Sort)
		pipe = append(pipe, bson.D{{Key: "$sort", Value: params.Options.Sort}})
	}
	if params.Options.Skip != 0 {
		// opts.SetSkip(params.Options.Skip)
		pipe = append(pipe, bson.D{{Key: "$skip", Value: params.Options.Skip}})
	}
	if params.Options.Limit != 0 {
		// opts.SetLimit(params.Options.Limit)
		pipe = append(pipe, bson.D{{Key: "$limit", Value: params.Options.Limit}})
	}

	// pipe = append(pipe, bson.D{
	// 	{Key: "$group", Value: bson.M{
	// 		"_id":    "$title",
	// 		"count": bson.M{"$sum": 1},
	// }}})

	// // pipe = append(pipe, bson.D{{"$unwind", bson.D{{"path", "$groups"}, {"preserveNullAndEmptyArrays", true}}}})
	// pipe = append(pipe, bson.D{{"$lookup", bson.M{
	// 	"from":         "component_schemas",
	// 	"as":           "schema",
	// 	"localField":   "_id",
	// 	"foreignField": "componentId",
	// }}})
	// pipe = append(pipe, bson.D{{"$unwind", bson.D{{"path", "$schema"}, {"preserveNullAndEmptyArrays", false}}}})
	// pipe = append(pipe, bson.D{{"$lookup", bson.M{
	// 	"from":         "librarys",
	// 	"as":           "schema.library",
	// 	"localField":   "schema._id",
	// 	"foreignField": "libraryId",
	// }}})
	// pipe = append(pipe, bson.D{{"$group", bson.M{
	// 	"_id":    "$_id",
	// 	"name":   bson.M{"$first": "$name"},
	// 	"schema": bson.M{"$push": "$schema"},
	// }}})
	// pipe = append(pipe, bson.D{{"$unwind", bson.D{{"path", "$schema"}, {"preserveNullAndEmptyArrays", true}}}})

	// pipe = append(pipe, bson.D{{"$project", bson.M{
	// 	"schema.schema_data": bson.M{"$arrayElemAt": []interface{}{"$schema.schema_data", 1}},
	// }}})

	// Take first element from the populated array (there is only one)
	// aggProject = bson.M{"$project": bson.M{
	//   "parent": bson.M{"$arrayElemAt": []interface{}{"$parent", 0}},
	// }}

	return pipe, nil
}

func CreateFilterAndOptions(params domain.RequestParams) (interface{}, *options.FindOptions, error) {
	opts := options.Find()
	if params.Options.Sort != nil {
		opts.SetSort(params.Options.Sort)
	}
	if params.Options.Skip != 0 {
		opts.SetSkip(params.Options.Skip)
	}
	if params.Options.Limit != 0 {
		opts.SetLimit(params.Options.Limit)
	}

	filter := params.Filter

	return filter, opts, nil
}
