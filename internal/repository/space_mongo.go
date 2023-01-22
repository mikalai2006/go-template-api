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

type SpaceMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewSpaceMongo(db *mongo.Database, i18n config.I18nConfig) *SpaceMongo {
	return &SpaceMongo{db: db, i18n: i18n}
}

func (r *SpaceMongo) CreateSpace(userID string, space *domain.SpaceInput) (*domain.Space, error) {
	var result *domain.Space

	collection := r.db.Collection(tblSpace)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newItem := domain.Space{
		UserID:      userIDPrimitive,
		Title:       space.Title,
		Description: space.Description,
		Setting:     space.Setting,
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

func (r *SpaceMongo) GetSpace(id string) (domain.Space, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Space

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Space{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblSpace).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Space{}, err
	}

	return result, nil
}

func (r *SpaceMongo) FindSpace(params domain.RequestParams) (domain.Response[domain.Space], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// opts := getPaginationOpts(&params.PaginationQuery)
	// f := createFilter(&params.PageFilter)
	// fmt.Println("f===", f)

	var results []domain.Space
	var response domain.Response[domain.Space]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Space]{}, err
	}
	// var pipe mongo.Pipeline
	// pipe = append(pipe, bson.D{{"$match", params.PageFilter}})
	// fmt.Printf("params - %v", pipe)

	cursor, err := r.db.Collection(tblSpace).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Space, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// var options options.CountOptions
	// // options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	// count, err := r.db.Collection(tblPage).CountDocuments(ctx, bson.M{}, &options)
	// if err != nil {
	// 	return response, err
	// }

	response = domain.Response[domain.Space]{
		Total: 0,
		Skip:  int(params.Skip),
		Limit: int(params.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *SpaceMongo) UpdateSpace(id string, data interface{}) (domain.Space, error) {
	var result domain.Space
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblSpace)

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

func (r *SpaceMongo) DeleteSpace(id string) (domain.Space, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Space{}
	collection := r.db.Collection(tblSpace)

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
