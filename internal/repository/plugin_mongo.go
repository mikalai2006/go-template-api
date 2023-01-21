package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PluginMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewPluginMongo(db *mongo.Database, i18n config.I18nConfig) *PluginMongo {
	return &PluginMongo{db: db, i18n: i18n}
}

func (r *PluginMongo) CreatePlugin(userID string, plugin *domain.PluginInput) (*domain.Plugin, error) {
	var result *domain.Plugin

	collection := r.db.Collection(tblPlugin)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	SpaceIDPrimitive, err := primitive.ObjectIDFromHex(plugin.SpaceID)
	if err != nil {
		return nil, err
	}
	newItem := domain.Plugin{
		UserID:      userIDPrimitive,
		SpaceID:     SpaceIDPrimitive,
		Name:        plugin.Name,
		Description: plugin.Description,
		Body:        plugin.Body,
		Code:        plugin.Code,
		Options:     plugin.Options,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PluginMongo) GetPlugin(id string) (domain.Plugin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Plugin

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Plugin{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblPlugin).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Plugin{}, err
	}

	return result, nil
}

func (r *PluginMongo) FindPlugin(params domain.RequestParams) (domain.Response[domain.Plugin], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Plugin
	var response domain.Response[domain.Plugin]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Plugin]{}, err
	}

	cursor, err := r.db.Collection(tblPlugin).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Plugin, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	response = domain.Response[domain.Plugin]{
		Total: 0,
		Skip:  int(params.Skip),
		Limit: int(params.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *PluginMongo) UpdatePlugin(id string, data interface{}) (domain.Plugin, error) {
	var result domain.Plugin
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblPlugin)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *PluginMongo) DeletePlugin(id string) (domain.Plugin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Plugin{}
	collection := r.db.Collection(tblPlugin)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}
