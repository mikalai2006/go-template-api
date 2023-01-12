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
)

type UserMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewUserMongo(db *mongo.Database, i18n config.I18nConfig) *UserMongo {
	return &UserMongo{db: db, i18n: i18n}
}

func (r *UserMongo) Iam(userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.User

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.M{"user_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}

func (r *UserMongo) GetUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.User

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}

func (r *UserMongo) FindUser(params domain.RequestParams) (domain.Response[domain.User], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.User
	var response domain.Response[domain.User]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.User]{}, err
	}
	// fmt.Println("gogo", pipe)
	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.User, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(tblUsers).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.User]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *UserMongo) CreateUser(userID string, user *domain.User) (*domain.User, error) {
	var result *domain.User

	collection := r.db.Collection(tblUsers)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newUser := domain.User{
		Avatar:    user.Avatar,
		Name:      user.Name,
		UserID:    userIDPrimitive,
		Login:     user.Login,
		Lang:      user.Lang,
		Currency:  user.Currency,
		Online:    user.Online,
		Verify:    user.Verify,
		LastTime:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblUsers).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserMongo) DeleteUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.User{}
	collection := r.db.Collection(tblUsers)

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

func (r *UserMongo) UpdateUser(id string, user *domain.User) (domain.User, error) {
	var result domain.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	data, err := utils.GetBodyToData(user)
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
