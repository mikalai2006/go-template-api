package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PageMongo struct {
	db *mongo.Database
}

func NewPageMongo(db *mongo.Database) *PageMongo {
	return &PageMongo{db:db}
}


func (r *PageMongo) GetPage(id string) (domain.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Page

	userIdPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Page{}, err
	}

	filter := bson.M{"_id": userIdPrimitive}

	err = r.db.Collection(tblPage).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Page{}, err
	}

	return result, nil
}

func (r *PageMongo) FindPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Page
	var response domain.Response[domain.Page]
	pipe, err := CreatePipeline(params)
	if err != nil {
		return domain.Response[domain.Page]{}, err
	}

	cursor, err := r.db.Collection(tblPage).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return response, err
	}

	var resultSlice []domain.Page = make([]domain.Page, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count,err := r.db.Collection(tblPage).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Page]{
		Total: count,
		Skip: int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data: resultSlice,
	}
	return response, nil
}

func (r *PageMongo) CreatePage(userId string, page domain.Page) (*domain.Page, error) {
	var result *domain.Page

	collection := r.db.Collection(tblPage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIdPrimitive, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	newPage := domain.Page{
		Name: page.Name,
		Title: page.Title,
		UserId: userIdPrimitive,
		Slug: page.Slug,
		SlugFull: fmt.Sprintf("/%s", page.Slug),
		ComponentId: page.ComponentId,
		Publish: page.Publish,
		SortOrder: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPage)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblPage).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PageMongo) DeletePage(id string) (domain.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Page{}
	collection := r.db.Collection(tblPage)

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

func (r *PageMongo) UpdatePage(id string, page domain.Page) (domain.Page, error) {
	var result domain.Page
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblPage)

	data, err := utils.GetBodyToData(page)
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