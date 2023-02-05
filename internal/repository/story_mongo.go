package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoryMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewStoryMongo(db *mongo.Database, i18n config.I18nConfig) *StoryMongo {
	return &StoryMongo{db: db, i18n: i18n}
}

func (r *StoryMongo) PublishStory(id string, data domain.StoryInputData) (domain.Story, error) {
	var result domain.Story
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblStory)

	userID, err := primitive.ObjectIDFromHex(data.UserID) // primitive.ObjectID
	if err != nil {
		return result, err
	}

	// disable old data.
	spaceID, err := primitive.ObjectIDFromHex(data.SpaceID) // primitive.ObjectID
	if err != nil {
		return result, err
	}
	// if LID, ok := data["space_id"]; ok {
	// 	spaceID = LID.(primitive.ObjectID)
	// 	// fmt.Println("layout_id=", layoutID)
	// }
	pageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	// fmt.Println("Story_id=", StoryID)
	filter := bson.M{
		"$or": bson.A{
			bson.M{
				"page_id":  pageID,
				"space_id": spaceID,
			},
			// bson.M{
			// 	"page_id":  pageID,
			// 	"space_id": spaceID,
			// },
		},
	}

	// unpublish old story data.
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "publish", Value: false}}}}
	_, err = collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return result, err
	}

	// publish new story data.
	// if val, ok := data["content"]; ok {
	// 	data := val.([]domain.ComponentData)
	// 	var ui []interface{}
	// 	for _, t := range data {
	// 		ui = append(ui, t)
	// 	}
	newItemData := domain.Story{
		UserID:    userID,
		SpaceID:   spaceID,
		PageID:    pageID,
		LayoutID:  data.LayoutID,
		Name:      data.Name,
		Title:     data.Title,
		Slug:      data.Slug,
		SlugFull:  data.SlugFull,
		Content:   data.Content,
		Publish:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	newItem, err := collection.InsertOne(ctx, newItemData)
	if err != nil {
		return result, err
	}
	insertedID := newItem.InsertedID // .(primitive.ObjectID).Hex()
	// }

	err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *StoryMongo) GetStory(params domain.RequestParams) (domain.Story, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	result := domain.Story{}
	results := []domain.Story{}

	params.Filter.(bson.M)["publish"] = true
	pipe, err := CreatePipeline(params, &r.i18n)
	// userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return result, err
	// }

	// filter := bson.M{"_id": userIDPrimitive}

	cursor, err := r.db.Collection(tblStory).Aggregate(ctx, pipe) // .FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return result, er
	}

	if len(results) > 0 {
		result = results[0]
	} else {
		return result, errors.New("not found item")
	}
	// if cursor.Next(ctx) {
	// 	if er := cursor.Decode(&result); er != nil {
	// 		return result, er
	// 	}
	// }

	return result, nil
}
