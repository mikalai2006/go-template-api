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

type ComponentMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewComponentMongo(db *mongo.Database, i18n config.I18nConfig) *ComponentMongo {
	return &ComponentMongo{db: db, i18n: i18n}
}

func (r *ComponentMongo) GetComponent(id string) (domain.Component, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	result := domain.Component{}

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": userIDPrimitive}

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
	pipe, err := CreatePipeline(params, &r.i18n)

	// Populate Parent field
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "component_schemas",
		"as":   "schema",
		// "localField":   "_id",
		// "foreignField": "componentId",
		"let": bson.D{{Key: "componentId", Value: "$_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$componentId", "$$componentId"}}}}},

			bson.D{{Key: "$lookup", Value: bson.M{
				"from": "librarys",
				"as":   "library",
				// "let":  bson.D{{"libraryId", "$libraryId"}},
				// "pipeline": mongo.Pipeline{
				// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$libraryId"}}}}},
				// 	bson.D{{"$unwind", bson.D{{"path", "$library"}, {"preserveNullAndEmptyArrays", true}}}},
				// },
				"localField":   "libraryId",
				"foreignField": "_id",
			}}},

			bson.D{{Key: "$lookup", Value: bson.M{
				"from": "component_schemadatas",
				"as":   "schema_data",
				// "let":  bson.D{{"libraryId", "$libraryId"}},
				// "pipeline": mongo.Pipeline{
				// 	bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$libraryId"}}}}},
				// 	bson.D{{"$unwind", bson.D{{"path", "$library"}, {"preserveNullAndEmptyArrays", true}}}},
				// },
				"localField":   "_id",
				"foreignField": "schemaId",
			}}},
			// bson.D{{"$project", bson.M{
			// 	"library": bson.M{"$arrayElemAt": []interface{}{"$library", 0}},
			// }}},
			// bson.D{{"$project", bson.M{
			// 	"schema": bson.M{"$arrayElemAt": []interface{}{"$schema", 0}},
			// }}},
		},
	}}}, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "component_groups",
		"as":           "groups",
		"localField":   "group",
		"foreignField": "_id",
	}}})

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(tblComponent).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// fmt.Println("results=", results[0].Schema)

	resultSlice := make([]domain.Component, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblComponent).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Component]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ComponentMongo) CreateComponent(userID string, component *domain.ComponentCreate) (*domain.Component, error) {
	var result *domain.Component

	collection := r.db.Collection(tblComponent)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newItem := domain.Component{
		Name:      component.Name,
		Title:     component.Title,
		UserID:    userIDPrimitive,
		Group:     component.Group,
		IsPage:    component.IsPage,
		IsGlobal:  component.IsGlobal,
		IsLayout:  component.IsLayout,
		SortOrder: 0,
		Publish:   component.Publish,
		Tpl:       "tpl",
		Setting:   component.Setting,
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
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ComponentMongo) UpdateComponent(id string, data interface{}) (domain.Component, error) {
	var result domain.Component
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblComponent)

	// data, err := utils.GetBodyToData(component)
	// if err != nil {
	// 	return result, err
	// }

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

func (r *ComponentMongo) FindByPopulate(params domain.RequestParams) (domain.Response[domain.Component], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.Component{}
	response := domain.Response[domain.Component]{}

	// agSearch := bson.M{"$match": params.Filter}
	// var agSort bson.D
	// if params.Options.Sort != nil {
	// 	// opts.SetSort(params.Options.Sort)
	// 	agSort = bson.D{{Key: "$sort", Value: params.Options.Sort}}
	// }
	// var agSkip bson.D
	// if params.Options.Skip != 0 {
	// 	agSkip = bson.D{{Key: "$skip", Value: params.Options.Skip}}
	// }
	// var agLimit bson.D
	// if params.Options.Limit != 0 {
	// 	agLimit = bson.D{{Key: "$limit", Value: params.Options.Limit}}
	// }

	// pipe = append(pipe, bson.D{
	// 	{Key: "$group", Value: bson.M{
	// 		"_id":    "$title",
	// 		"count": bson.M{"$sum": 1},
	// }}})

	// Populate Parent field
	agPopulate := bson.M{
		"$lookup": bson.M{
			"from":         "component_schemas",
			"as":           "schema",
			"foreignField": "componentId",
			"localField":   "_id",
		},
	}

	// agUnWind := bson.M{"$unwind": bson.D{{"path", "$schem"}, {"preserveNullAndEmptyArrays", true}}}
	// agProject := bson.D{{"$project", bson.M{
	// 	"schema": bson.M{"$arrayElemAt": []interface{}{"$schema", 0}},
	// }}}

	// Take first element from the populated array (there is only one)
	// aggProject = bson.M{"$project": bson.M{
	//   "parent": bson.M{"$arrayElemAt": []interface{}{"$parent", 0}},
	// }}

	cursor, err := r.db.Collection(tblComponent).Aggregate(ctx, []bson.M{
		// agSearch,
		// agSort, agSkip, agLimit,
		agPopulate,
		// agUnWind,
		// agProject,
	}) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// fmt.Println("results=", results)

	resultSlice := make([]domain.Component, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblComponent).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Component]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ComponentMongo) FindLibrarys(params domain.RequestParams) (domain.Response[domain.Library], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	results := []domain.Library{}
	response := domain.Response[domain.Library]{}
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return response, err
	}

	pipe = append(pipe, bson.D{{
		Key: "$lookup", Value: bson.M{
			"from": "fields",
			"let":  bson.M{"libraryId": "$_id"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{
					"$expr":   bson.M{"$eq": [2]string{"$libraryId", "$$libraryId"}},
					"publish": true,
				}}},
			},
			"as": "data",
			// "localField":   "_id",
			// "foreignField": "libraryId",
		},
	}})
	// agUnWind := bson.M{"$unwind": bson.D{{"path", "$data"}, {"preserveNullAndEmptyArrays", true}}}

	cursor, err := r.db.Collection(tblLibrary).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	// fmt.Println("results=", results)

	resultSlice := make([]domain.Library, len(results))
	copy(resultSlice, results)

	for keyLibrary := range resultSlice {
		// for keyData, _ := range library.Data {
		// resultSlice[keyLibrary].Gogo = append(resultSlice[keyLibrary].Gogo, buildTree(library.Data, ))
		relations := relationMap(resultSlice[keyLibrary].Data)
		// resultSlice[keyLibrary].Gogo = relations
		treeLibrary := buildTree(resultSlice[keyLibrary].Data, relations, 0, r.i18n)

		var treeRoot interface{}
		for _, tree := range treeLibrary {
			if tree != nil {
				treeRoot = tree
			}
		}

		resultSlice[keyLibrary].Tree = treeRoot
		// if err != nil {
		// 	return response, err
		// }
		// resultSlice[keyLibrary].Tree = res
		// do something here
		// fmt.Println(val)
	}

	count, err := r.db.Collection(tblLibrary).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Library]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}

	return response, nil
}

func relationMap(dataset []*domain.Field) map[string][]*domain.Field {
	relations := make(map[string][]*domain.Field)
	for _, relation := range dataset {
		child, parent := relation, relation.Parent
		relations[parent] = append(relations[parent], child)
	}
	return relations
}

func buildTree(
	fields []*domain.Field,
	relations map[string][]*domain.Field,
	level int,
	i18n config.I18nConfig,
) []interface{} {
	level++
	tree := make([]interface{}, len(fields))
	for i, field := range fields {
		node := map[string]interface{}{
			"_uid":   field.UID,
			"parent": field.Parent,
			"name":   field.Name,
			"level":  level,
		}

		for k, fieldData := range field.Data.Value {
			switch c := fieldData.(type) {
			case map[string]interface{}:
				if _, ok := c[i18n.Default]; ok {
					for i, data := range fieldData.(map[string]any) {
						key := fmt.Sprintf("%s%s%s", k, i18n.Prefix, i)
						node[key] = data
						if i == i18n.Default {
							node[k] = data
						}
					}
					// delete(field.Data.Value, k)
				}
			default:
				node[k] = fieldData
			}
		}

		if childIDS, ok := relations[field.UID]; ok {
			node["child"] = buildTree(childIDS, relations, level, i18n)
		}
		if node["parent"] == nil || node["level"].(int) != 1 {
			tree[i] = node
		}
	}
	return tree
}
