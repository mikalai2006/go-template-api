package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db *mongo.Database
}

func NewAuthMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{db: db}
}

func (r *AuthMongo) CreateAuth(user *domain.Auth) (string, error) {
	var id string

	collection := r.db.Collection(TblAuth)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	auth, err := collection.InsertOne(ctx, user)
	if err != nil {
		return id, err
	}
	id = auth.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func chooseProvider(auth *domain.SignInInput) bson.D {
	if auth.Strategy == "local" {
		var filter bson.A
		if auth.Email != "" {
			// add email filter
			filter = append(filter, bson.M{"login": auth.Email}, bson.M{"email": auth.Email})
		}
		if auth.Login != "" {
			// add email filter
			filter = append(filter, bson.M{"login": auth.Login}, bson.M{"email": auth.Login})
		}
		return bson.D{
			{Key: "$or", Value: filter},
			{Key: "password", Value: auth.Password},
		}
	}

	if auth.VkID != "" {
		return bson.D{{Key: "vkid", Value: auth.VkID}}
	} else if auth.GoogleID != "" {
		return bson.D{{Key: "googleid", Value: auth.GoogleID}}
	}
	return bson.D{{Key: "vkid", Value: "none"}}
}

func (r *AuthMongo) CheckExistAuth(auth *domain.SignInInput) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	filter := chooseProvider(auth)
	err := r.db.Collection(TblAuth).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = nil
		}
	}
	return user, err
}

func (r *AuthMongo) GetAuth(auth *domain.Auth) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	query := bson.M{"login": auth.Login, "password": auth.Password}
	// fmt.Println("")
	// fmt.Printf("GetAuth: query=%s", query)
	err := r.db.Collection(TblAuth).FindOne(ctx, query).Decode(&user)

	return user, err
}

func (r *AuthMongo) GetByCredentials(auth *domain.SignInInput) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	filter := chooseProvider(auth)

	err := r.db.Collection(TblAuth).FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (r *AuthMongo) SetSession(authID primitive.ObjectID, session domain.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	_, err := r.db.Collection(TblAuth).UpdateOne(
		ctx,
		bson.M{"_id": authID},
		bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}},
	)

	return err
}

func (r *AuthMongo) VerificationCode(userID, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	res, err := r.db.Collection(TblAuth).UpdateOne(ctx,
		bson.M{"verification.code": code, "_id": id},
		bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed verify")
	}

	return nil
}

func (r *AuthMongo) RefreshToken(refreshToken string) (domain.Auth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Auth

	if err := r.db.Collection(TblAuth).FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
