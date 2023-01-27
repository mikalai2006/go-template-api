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

type ComponentGroupMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewComponentGroupMongo(db *mongo.Database, i18n config.I18nConfig) *ComponentGroupMongo {
	return &ComponentGroupMongo{db: db, i18n: i18n}
}

func (r *ComponentGroupMongo) FindComponentGroup(params domain.RequestParams) (domain.Response[domain.ComponentGroup], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.ComponentGroup{}
	response := domain.Response[domain.ComponentGroup]{}
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(tblComponentGroup).Aggregate(ctx, pipe) // .Find(ctx, params)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.ComponentGroup, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(tblComponentGroup).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.ComponentGroup]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ComponentGroupMongo) CreateComponentGroup(userID string, componentGroup *domain.ComponentGroup) (*domain.ComponentGroup, error) {
	var result *domain.ComponentGroup

	collection := r.db.Collection(tblComponentGroup)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	componentGroup.UserID = userIDPrimitive
	componentGroup.CreatedAt = time.Now()
	componentGroup.UpdatedAt = time.Now()

	res, err := collection.InsertOne(ctx, componentGroup)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblComponentGroup).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ComponentGroupMongo) UpdateComponentGroup(id string, data interface{}) (domain.ComponentGroup, error) {
	var result domain.ComponentGroup
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblComponentGroup)

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

func (r *ComponentGroupMongo) DeleteComponentGroup(id string) (domain.ComponentGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.ComponentGroup{}
	collection := r.db.Collection(tblComponentGroup)

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
