package repository

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewPageMongo(db *mongo.Database, i18n config.I18nConfig) *PageMongo {
	return &PageMongo{db: db, i18n: i18n}
}

func (r *PageMongo) CreatePage(userID string, page *domain.Page) (*domain.Page, error) {
	var result *domain.Page

	collection := r.db.Collection(tblPage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newPage := domain.Page{
		Name:        page.Name,
		Title:       page.Title,
		UserID:      userIDPrimitive,
		Slug:        page.Slug,
		SlugFull:    fmt.Sprintf("/%s", page.Slug),
		ComponentID: page.ComponentID,
		Publish:     page.Publish,
		LayoutID:    page.LayoutID,
		SortOrder:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

func (r *PageMongo) GetFullPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Page
	var response domain.Response[domain.Page]
	pipe, err := CreatePipeline(params, &r.i18n)

	// Populate Parent field.
	pipe = append(pipe,
		bson.D{{
			Key: "$lookup",
			Value: bson.M{
				"from": tblComponent,
				"as":   "component",
				// "localField":   "_id",
				// "foreignField": "componentId",
				"let": bson.D{{Key: "componentId", Value: "$componentId"}},
				"pipeline": mongo.Pipeline{
					bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$componentId"}}}}},
				},
			},
		}},
		bson.D{{Key: "$unwind", Value: "$component"}},
		// bson.D{{"$project", bson.M{
		// 	"component": bson.M{"$arrayElemAt": []interface{}{"$component", 0}},
		// }}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": tblComponentData,
			"as":   "component_data",
			// "localField":   "_id",
			// "foreignField": "pageId",
			"let": bson.D{{Key: "pageId", Value: "$_id"}, {Key: "layoutId", Value: "$layoutId"}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$and", Value: bson.A{
						bson.D{{Key: "$or",
							Value: bson.A{
								bson.M{"$expr": bson.M{"$eq": [2]string{"$pageId", "$$pageId"}}},
								bson.D{{Key: "$and", Value: bson.A{
									bson.M{"$expr": bson.M{"$eq": [2]string{"$layoutId", "$$layoutId"}}},
									// bson.M{"$expr": bson.M{"$eq": [2]string{"$pageId", string(primitive.NilObjectID[0])}}},
									bson.M{"pageId": nil},
								}},
								},
							},
						}},
						bson.M{"publish": true},
					},
					},
				}}},
			},
		}},
		},
		bson.D{
			{
				Key: "$lookup",
				Value: bson.M{
					"from": tblComponent,
					"as":   "layout",
					"let":  bson.D{{Key: "layoutId", Value: "$layoutId"}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$layoutId"}}}}},
					},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: "$layout"}},
	)

	if err != nil {
		return domain.Response[domain.Page]{}, err
	}
	cursor, err := r.db.Collection(tblPage).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Page, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	options.SetSkip(params.Skip)
	count, err := r.db.Collection(tblPage).CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	for keyPage := range resultSlice {
		// mapP := relationMapX(resultSlice[keyPage].ComponentData, r.i18n)
		// resultSlice[keyPage].XXX = mapP
		datas := filterComponentData(resultSlice[keyPage].ComponentData, mytest)
		mapP := createContent(resultSlice[keyPage].ComponentData, datas, r.i18n)
		if len(mapP) > 0 {
			resultSlice[keyPage].Content = mapP[0]
		}
		// fmt.Println("pages", pages)
	}

	response = domain.Response[domain.Page]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}
func mytest(s domain.ComponentData) bool {
	return s.Parent == "page"
}

func filterComponentData(
	ss []domain.ComponentData,
	test func(domain.ComponentData) bool) (ret []domain.ComponentData) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func filterComponentData2(
	ss []domain.ComponentData,
	ids interface{},
	test func(domain.ComponentData, interface{}) bool) (ret []domain.ComponentData) {
	for _, s := range ss {
		if test(s, ids) {
			ret = append(ret, s)
		}
	}
	return
}
func mytestx(s domain.ComponentData, ids interface{}) bool {
	kindIds := reflect.TypeOf(ids).Kind()
	if kindIds != reflect.Slice {
		return false
	}
	sIds := reflect.ValueOf(ids)
	for i := 0; i < sIds.Len(); i++ {
		if s.UID == sIds.Index(i).Elem().String() {
			return true
		}
	}
	return false
}

func createContent(
	allData []domain.ComponentData,
	datasets []domain.ComponentData,
	i18n config.I18nConfig,
) []interface{} {
	s := reflect.ValueOf(datasets)
	if s.Kind() != reflect.Slice {
		fmt.Println("plot() given a non-slice type")
	}
	lenSlice := s.Len()
	relations := make([]interface{}, 0, lenSlice)

	for keyData := range datasets {
		newBlok := map[string]interface{}{
			"_uid":      datasets[keyData].UID,
			"parent":    datasets[keyData].Parent,
			"component": datasets[keyData].Component,
		}

		for keyProperty, property := range datasets[keyData].Data {
			// resultItem := map[string]interface{}{}
			if reflect.ValueOf(property).Kind() != reflect.Map {
				// if property not uids.
				newBlok[keyProperty] = property
				continue
			}
			// if slice uids in data.
			if val, ok := property.(map[string]interface{})["uids"]; ok {
				newBlok[keyProperty] = property
				// recurse.
				nested := filterComponentData2(allData, val, mytestx)
				newBlok[keyProperty] = createContent(allData, nested, i18n)
			}
			// if map i18n.
			if _, ok := property.(map[string]interface{})[i18n.Default]; ok {
				for i, data := range property.(map[string]interface{}) {
					key := fmt.Sprintf("%s%s%s", keyProperty, i18n.Prefix, i)
					newBlok[key] = data
					if i == i18n.Default {
						newBlok[keyProperty] = data
					}
				}
			}
		}

		relations = append(relations, newBlok)
	}
	return relations
}

// func relationMapX(dataset []*domain.ComponentData, i18n config.I18nConfig) map[string]interface{} {
// 	relations := make(map[string]*domain.ComponentData)
// 	// range group bloks.
// 	for _, relation := range dataset {
// 		parent := relation.UID
// 		relations[parent] = relation // append(relations[parent], child)
// 	}

// 	/* structure each item.
// 	_uid
// 	parent
// 	component (short co)
// 	--- other few data.
// 	*/

// 	// go to poperties each group blok.
// 	result := make(map[string]interface{})
// 	for keyBlok, blok := range relations {
// 		resultItem := map[string]interface{}{}
// 		for keyProperty, property := range blok.Data {
// 			// fmt.Println("type=", keyProperty, "=")
// 			if reflect.ValueOf(property).Kind() == reflect.Map {
// 				if _, ok := property.(map[string]interface{})["uids"]; ok {
// 					resultItem[keyProperty] = property
// 				}
// 				if _, ok := property.(map[string]interface{})[i18n.Default]; ok {
// 					for i, data := range property.(map[string]interface{}) {
// 						key := fmt.Sprintf("%s%s%s", keyProperty, i18n.Prefix, i)
// 						resultItem[key] = data
// 						if i == i18n.Default {
// 							resultItem[keyProperty] = data
// 						}
// 					}
// 				}
// 			} else {
// 				resultItem[keyProperty] = property
// 			}
// 		}
// 		resultItem["_uid"] = blok.UID
// 		resultItem["parent"] = blok.Parent
// 		result[keyBlok] = resultItem
// 	}

// 	// create tree.
// 	tree := map[string]interface{}{}
// 	for keyItem, blok := range result {
// 		newBlok := map[string]interface{}{}
// 		for keyProperty, property := range blok.(map[string]interface{}) {
// 			// fmt.Println("property type: ", keyProperty, reflect.ValueOf(property).Kind())
// 			if reflect.ValueOf(property).Kind() == reflect.Map {
// 				// append blok with uids.
// 				if uids, ok := property.(map[string]interface{})["uids"]; ok {
// 					var sliceUids []interface{}
// 					for _, uid := range uids.(primitive.A) {
// 						if val, ok := result[uid.(string)]; ok {
// 							sliceUids = append(sliceUids, val)
// 						}
// 					}
// 					newBlok[keyProperty] = sliceUids
// 				} else {
// 					newBlok[keyProperty] = property
// 				}
// 			} else {
// 				newBlok[keyProperty] = property
// 				// fmt.Println("property=", property)
// 			}

// 		}

// 		if val, ok := blok.(map[string]interface{})["parent"]; ok {
// 			fmt.Println(keyItem, " ===== ", val.(string), val.(string) == "page")
// 			if val.(string) == "page" {
// 				tree[keyItem] = newBlok
// 			}
// 		}
// 	}

// 	return tree
// }

func (r *PageMongo) GetPage(id string) (domain.Page, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Page

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Page{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblPage).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Page{}, err
	}

	return result, nil
}

func (r *PageMongo) GetPageForRouters() (domain.Response[domain.PageRoutes], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.PageRoutes
	var response domain.Response[domain.PageRoutes]
	filter := bson.M{"publish": true}
	cursor, err := r.db.Collection(tblPage).Find(ctx, filter)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.PageRoutes, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(tblPage).CountDocuments(ctx, filter)
	if err != nil {
		return response, err
	}
	response = domain.Response[domain.PageRoutes]{
		Total: int(count),
		Skip:  0,
		Limit: int(count),
		Data:  resultSlice,
	}

	return response, nil
}

func (r *PageMongo) FindPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// opts := getPaginationOpts(&params.PaginationQuery)
	// f := createFilter(&params.PageFilter)
	// fmt.Println("f===", f)

	var results []domain.Page
	var response domain.Response[domain.Page]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Page]{}, err
	}
	// var pipe mongo.Pipeline
	// pipe = append(pipe, bson.D{{"$match", params.PageFilter}})
	fmt.Printf("params - %v", pipe)

	cursor, err := r.db.Collection(tblPage).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Page, len(results))
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

	response = domain.Response[domain.Page]{
		Total: 0,
		Skip:  int(params.Skip),
		Limit: int(params.Limit),
		Data:  resultSlice,
	}
	return response, nil
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
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *PageMongo) UpdatePage(id string, data interface{}) (domain.Page, error) {
	var result domain.Page
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblPage)

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
