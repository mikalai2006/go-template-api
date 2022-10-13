package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ComponentMongo struct {
	db *mongo.Database
}

func NewComponentMongo(db *mongo.Database) *ComponentMongo {
	return &ComponentMongo{db:db}
}


func (r *ComponentMongo) GetComponent(id string) (domain.Component, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	result := domain.Component{}

	userIdPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": userIdPrimitive}

	err = r.db.Collection(tblComponent).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ComponentMongo) FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.Component{}
	response := domain.Response[domain.Component]{}
	pipe, err := CreatePipeline(params)
	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(tblComponent).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return response, err
	}

	var resultSlice []domain.Component = make([]domain.Component, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count,err := r.db.Collection(tblComponent).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Component]{
		Total: count,
		Skip: int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data: resultSlice,
	}
	return response, nil
}

func (r *ComponentMongo) CreateComponent(userId string, component domain.Component) (*domain.Component, error) {
	var result *domain.Component

	collection := r.db.Collection(tblComponent)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIdPrimitive, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	newItem := domain.Component{
		Name: component.Name,
		Title: component.Title,
		UserId: userIdPrimitive,
		Publish: component.Publish,
		Group: component.Group,
		Tpl: "tpl",
		SortOrder: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblComponent).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ComponentMongo) DeleteComponent(id string) (domain.Component, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Component{}
	collection := r.db.Collection(tblComponent)

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
	if err != nil  {
		return result, err
	}

	return result, nil
}

func (r *ComponentMongo) UpdateComponent(id string, component domain.Component) (domain.Component, error) {
	var result domain.Component
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblComponent)

	data, err := utils.GetBodyToData(component)
	if err != nil {
		return result, err
	}

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

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