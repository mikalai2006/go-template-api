package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartnerMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewPartnerMongo(db *mongo.Database, i18n config.I18nConfig) *PartnerMongo {
	return &PartnerMongo{db: db, i18n: i18n}
}

func (r *PartnerMongo) CreatePartner(userID string, data *domain.PartnerInput) (domain.Partner, error) {
	var result domain.Partner

	collection := r.db.Collection(tblPartner)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}
	count, err := r.db.Collection(tblPartner).EstimatedDocumentCount(ctx)
	if err != nil {
		return result, err
	}

	// make pretty urls from title.
	prettyurl := utils.EncodeRus(data.Title) // fmt.Sprintf("%d-%s", count, utils.EncodeRus(data.Title[r.i18n.Default]))

	newItem := domain.Partner{
		UserID: userIDPrimitive,
		SeoID:  count,

		Title:       data.Title,
		Description: data.Description,

		Seo:    prettyurl,
		Locale: data.Locale,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newItem)
	if err != nil {
		return result, err
	}

	err = r.db.Collection(tblPartner).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *PartnerMongo) GetPartner(id string) (domain.Partner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Partner

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Partner{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblPartner).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Partner{}, err
	}

	return result, nil
}

func (r *PartnerMongo) FindPartner(params domain.RequestParams) (domain.Response[domain.PartnerPopulate], error) {
	collection := r.db.Collection(tblPartner)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.PartnerPopulate
	var response domain.Response[domain.PartnerPopulate]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return response, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": tblImage,
			"as":   "images",
			// "localField":   "_id",
			// "foreignField": "componentId",
			"let": bson.D{{Key: "serviceId", Value: "$_id"}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
			},
		},
	}})

	// fmt.Println(pipe)
	cursor, err := collection.Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.PartnerPopulate, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	count, err := collection.CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.PartnerPopulate]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *PartnerMongo) UpdatePartner(id string, data *domain.PartnerInput) (domain.Partner, error) {
	var result domain.Partner
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblPartner)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	userIDPrimitive, err := primitive.ObjectIDFromHex(data.UserID)
	if err != nil {
		return result, err
	}
	filter := bson.M{"_id": idPrimitive}

	// make pretty urls from title.
	prettyurl := utils.EncodeRus(data.Title) // fmt.Sprintf("%d-%s", count, utils.EncodeRus(data.Title[r.i18n.Default]))

	newItem := domain.Partner{
		UserID: userIDPrimitive,
		SeoID:  data.SeoID,

		Seo:         prettyurl,
		Title:       data.Title,
		Description: data.Description,

		Locale: data.Locale,

		CreatedAt: data.CreatedAt,
		UpdatedAt: time.Now(),
	}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newItem})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *PartnerMongo) DeletePartner(id string) (domain.Partner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Partner{}
	collection := r.db.Collection(tblPartner)

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
