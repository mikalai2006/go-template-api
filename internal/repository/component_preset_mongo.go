package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComponentPresetMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewComponentPresetMongo(db *mongo.Database, i18n config.I18nConfig) *ComponentPresetMongo {
	return &ComponentPresetMongo{db: db, i18n: i18n}
}

func (r *ComponentPresetMongo) FindComponentPreset(params domain.RequestParams) (domain.Response[domain.ComponentPreset], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.ComponentPreset{}
	response := domain.Response[domain.ComponentPreset]{}
	pipe, err := CreatePipeline(params, &r.i18n)

	// if err != nil {
	// 	return response, err
	// }
	fmt.Println("pipe preset", pipe)
	// cursor, err := r.db.Collection(tblComponentPreset).Find(ctx, bson.M{})
	cursor, err := r.db.Collection(tblComponentPreset).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.ComponentPreset, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(tblComponentPreset).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.ComponentPreset]{
		Total: int(count),
		Skip:  0,          // int(params.Options.Skip),
		Limit: int(count), // int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ComponentPresetMongo) CreateComponentPreset(userID string, preset *domain.ComponentPresetInput) (*domain.ComponentPreset, error) {
	var result *domain.ComponentPreset

	collection := r.db.Collection(tblComponentPreset)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	data := domain.ComponentPreset{
		ComponentID: preset.ComponentID,
		Title:       preset.Title,
		Description: preset.Description,
		Data:        preset.Data,
		UserID:      userIDPrimitive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblComponentPreset).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ComponentPresetMongo) UpdateComponentPreset(id string, data interface{}) (domain.ComponentPreset, error) {
	var result domain.ComponentPreset
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblComponentPreset)

	// data, err := utils.GetBodyToData(component)
	// if err != nil {
	// 	return result, err
	// }

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// data["user_id"] = idPrimitive

	filter := bson.M{"_id": idPrimitive}

	// fmt.Println("data=", data)
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

func (r *ComponentPresetMongo) DeleteComponentPreset(id string) (domain.ComponentPreset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.ComponentPreset{}
	collection := r.db.Collection(tblComponentPreset)

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
