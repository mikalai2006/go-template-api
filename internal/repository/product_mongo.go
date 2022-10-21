package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewProductMongo(db *mongo.Database, i18n config.I18nConfig) *ProductMongo {
	return &ProductMongo{db: db, i18n: i18n}
}

func (r *ProductMongo) CreateProduct(userID string, data *domain.ProductInput) (domain.Product, error) {
	var result domain.Product

	collection := r.db.Collection(TblProduct)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}
	count, err := r.db.Collection(TblProduct).EstimatedDocumentCount(ctx)
	if err != nil {
		return result, err
	}

	shopID, err := primitive.ObjectIDFromHex(data.ShopID)
	if err != nil {
		return result, err
	}
	categoryID, err := primitive.ObjectIDFromHex(data.CategoryID)
	if err != nil {
		return result, err
	}
	// make pretty urls from title.
	prettyurl := utils.EncodeRus(data.Title) // fmt.Sprintf("%d-%s", count, utils.EncodeRus(data.Title[r.i18n.Default]))

	newPage := domain.Product{
		Description: data.Description,
		Title:       data.Title,
		SeoID:       count,

		UserID:     userIDPrimitive,
		ShopID:     shopID,
		CategoryID: categoryID,
		Seo:        prettyurl,
		Locale:     data.Locale,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPage)
	if err != nil {
		return result, err
	}

	err = r.db.Collection(TblProduct).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ProductMongo) GetProduct(id string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Product

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Product{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(TblProduct).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Product{}, err
	}

	return result, nil
}

func (r *ProductMongo) FindProduct(params domain.RequestParams) (domain.Response[domain.Product], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Product
	var response domain.Response[domain.Product]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Product]{}, err
	}
	fmt.Println(pipe)
	cursor, err := r.db.Collection(TblProduct).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Product, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	count, err := r.db.Collection(TblProduct).CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Product]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ProductMongo) UpdateProduct(id string, data interface{}) (domain.Product, error) {
	var result domain.Product
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblProduct)

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

func (r *ProductMongo) DeleteProduct(id string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Product{}
	collection := r.db.Collection(TblProduct)

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
