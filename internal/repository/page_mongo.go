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

func (r *PageMongo) CreatePage(userID string, page *domain.PageInputData) (*domain.Page, error) {
	var result *domain.Page

	collection := r.db.Collection(tblPage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newPage := domain.PageInputData{
		UserID:      userIDPrimitive,
		SpaceID:     page.SpaceID,
		Name:        page.Name,
		Title:       page.Title,
		Slug:        page.Slug,
		SlugFull:    fmt.Sprintf("/%s", page.Slug),
		ComponentID: page.ComponentID,
		Publish:     page.Publish,
		LayoutID:    page.LayoutID,
		SortOrder:   page.SortOrder,
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
				"let": bson.D{{Key: "componentId", Value: "$component_id"}},
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
			"let": bson.D{{Key: "pageId", Value: "$_id"}, {Key: "layoutId", Value: "$layout_id"}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$and", Value: bson.A{
						bson.D{{Key: "$or",
							Value: bson.A{
								bson.M{"$expr": bson.M{"$eq": [2]string{"$page_id", "$$pageId"}}},
								bson.D{{Key: "$and", Value: bson.A{
									bson.M{"$expr": bson.M{"$eq": [2]string{"$layout_id", "$$layoutId"}}},
									// bson.M{"$expr": bson.M{"$eq": [2]string{"$pageId", string(primitive.NilObjectID[0])}}},
									bson.M{"page_id": primitive.NilObjectID},
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
					"let":  bson.D{{Key: "layoutId", Value: "$layout_id"}},
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
		// datas = append(datas, domain.ComponentData{
		// 	Parent: "page",
		// 	Data: map[string]interface{}{
		// 		"layout": map[string]interface{}{
		// 			"uids": primitive.A{
		// 				resultSlice[keyPage].LayoutID.Hex(),
		// 			},
		// 			"global": true,
		// 		},
		// 	},
		// 	Publish:   true,
		// 	PageID:    resultSlice[keyPage].ID,
		// 	LayoutID:  resultSlice[keyPage].LayoutID,
		// 	UID:       resultSlice[keyPage].Component.ID.Hex(),
		// 	Component: resultSlice[keyPage].Component.Name,
		// })
		mapAllData := map[string]domain.ComponentData{}
		for i, _ := range resultSlice[keyPage].ComponentData {
			if resultSlice[keyPage].ComponentData[i].Parent == "global" {
				resultSlice[keyPage].ComponentData[i].Parent = resultSlice[keyPage].ID.Hex()
			}
			// if id component changed, edit component name to page.
			if resultSlice[keyPage].ComponentData[i].Parent == "page" {
				resultSlice[keyPage].ComponentData[i].Component = resultSlice[keyPage].Component.Name
			}
			mapAllData[resultSlice[keyPage].ComponentData[i].UID] = resultSlice[keyPage].ComponentData[i]
		}

		// add default page node.
		if _, ok := mapAllData[resultSlice[keyPage].ID.Hex()]; ok {
			// fmt.Println("yes page data=")
		} else {
			// fmt.Println("no page data!")
			mapAllData[resultSlice[keyPage].ID.Hex()] = domain.ComponentData{
				Parent: "page",
				Data:   map[string]interface{}{
					// "layout": map[string]interface{}{
					// 	"uids": primitive.A{
					// 		resultSlice[keyPage].LayoutID.Hex(),
					// 	},
					// 	"global": true,
					// },
				},
				Publish:   true,
				PageID:    resultSlice[keyPage].ID,
				LayoutID:  resultSlice[keyPage].LayoutID,
				UID:       resultSlice[keyPage].ID.Hex(),
				Component: resultSlice[keyPage].Component.Name,
			}
			resultSlice[keyPage].ComponentData = append(resultSlice[keyPage].ComponentData, mapAllData[resultSlice[keyPage].ID.Hex()])
			// mapAllData[resultSlice[keyPage].LayoutID.Hex()] = domain.ComponentData{
			// 	Parent: resultSlice[keyPage].Component.ID.Hex(),
			// 	Data: map[string]interface{}{
			// 		"global": true,
			// 	},
			// 	Publish:   true,
			// 	PageID:    primitive.NilObjectID,
			// 	LayoutID:  resultSlice[keyPage].LayoutID,
			// 	UID:       resultSlice[keyPage].LayoutID.Hex(),
			// 	Component: resultSlice[keyPage].Layout.Name,
			// }
		}

		// add default layout node.
		if _, ok := mapAllData[resultSlice[keyPage].LayoutID.Hex()]; ok {
		} else {
			mapAllData[resultSlice[keyPage].LayoutID.Hex()] = domain.ComponentData{
				Parent: "layout", //resultSlice[keyPage].Component.ID.Hex(),
				Data: map[string]interface{}{
					"global": true,
				},
				Publish:   true,
				PageID:    primitive.NilObjectID,
				LayoutID:  resultSlice[keyPage].LayoutID,
				UID:       resultSlice[keyPage].LayoutID.Hex(),
				Component: resultSlice[keyPage].Layout.Name,
			}
			resultSlice[keyPage].ComponentData = append(resultSlice[keyPage].ComponentData, mapAllData[resultSlice[keyPage].LayoutID.Hex()])
		}
		datas := filterComponentData(resultSlice[keyPage].ComponentData, mytest)
		mapP := createContent(mapAllData, datas, resultSlice[keyPage].LayoutID.Hex(), resultSlice[keyPage].ID.Hex(), r.i18n)
		// fmt.Println("datas=", datas)
		// fmt.Println("=========")
		// fmt.Println("mapAllData=", mapAllData)
		if len(mapP) > 0 {
			// fmt.Println("yes content")
			resultSlice[keyPage].Content = mapP[0]
			// fmt.Printf("%#v\n", mapP[0])
		} else {
			// fmt.Println("no content! default")
			resultSlice[keyPage].Content = bson.M{
				"component": resultSlice[keyPage].Component.Name,
				"_uid":      resultSlice[keyPage].ID,
				"layout": bson.A{
					bson.M{"component": resultSlice[keyPage].Layout.Name, "_uid": resultSlice[keyPage].LayoutID, "global": true},
				},
			}
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
	return s.Parent == "layout"
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
		// fmt.Println("sIds.Index(i)=", i, sIds.Index(i))
		if s.UID == sIds.Index(i).Elem().String() {
			return true
		}
	}
	return false
}

func createContent(
	allData map[string]domain.ComponentData,
	datasets []domain.ComponentData,
	layoutID string,
	pageID string,
	// global bool,
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
			// "global":    global,
		}

		for keyProperty, property := range datasets[keyData].Data {
			// resultItem := map[string]interface{}{}
			if reflect.ValueOf(property).Kind() != reflect.Map {
				// if property not uids.
				newBlok[keyProperty] = property
				continue
			}
			// fmt.Println(" map key=", keyProperty, property)
			// detect global.
			// if keyProperty == "layout" {
			// 	global = true
			// }

			newBlok[keyProperty] = property
			// if slice uids in data.
			if val, ok := property.(map[string]interface{})["uids"]; ok {
				newBlok[keyProperty] = property
				// recurse.
				// nested := filterComponentData2(allData, val, mytestx)
				nested := []domain.ComponentData{} // make(, len(val.(primitive.A)))
				for _, uID := range val.(primitive.A) {
					if uID == "page" {
						uID = pageID
						// fmt.Println("is page", allData[uID.(string)])
					}
					nested = append(nested, allData[uID.(string)])
				}
				newBlok[keyProperty] = createContent(allData, nested, pageID, layoutID, i18n)
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
				// continue
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

	datas := filterComponentData(result.ComponentData, mytest)
	mapAllData := map[string]domain.ComponentData{}
	for i, _ := range result.ComponentData {
		if result.ComponentData[i].Parent == "layout" {
			result.ComponentData[i].Parent = result.Layout.ID.Hex()
		}
		mapAllData[result.ComponentData[i].UID] = result.ComponentData[i]
	}
	mapP := createContent(mapAllData, datas, result.Layout.ID.Hex(), result.ID.Hex(), r.i18n)
	if len(mapP) > 0 {
		result.Content = mapP[0]
	} else {
		result.Content = bson.M{
			"component": result.Component.Name,
			"_uid":      result.ID,
			"layout": bson.A{
				bson.M{"component": result.Layout.Name, "_uid": result.LayoutID, "global": true},
			},
		}
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
	// fmt.Printf("params - %v", pipe)

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

func (r *PageMongo) UpdatePageWithContent(id string, data map[string]interface{}) (domain.Page, error) {
	var result domain.Page
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblComponentData)

	// disable old data.
	var layoutID primitive.ObjectID
	if LID, ok := data["layout_id"]; ok {
		layoutID = LID.(primitive.ObjectID)
		// fmt.Println("layout_id=", layoutID)
	}
	pageID, _ := primitive.ObjectIDFromHex(id)
	// fmt.Println("page_id=", pageID)
	filter := bson.M{
		"$or": bson.A{
			bson.M{
				"page_id":   pageID,
				"layout_id": layoutID,
			},
			bson.M{
				"page_id":   primitive.NilObjectID,
				"layout_id": layoutID,
			},
		},
	}
	// bson.M{"$or": bson.A{
	// 	bson.M{"page_id": pageID},
	// 	bson.M{"layout_id": layoutID},
	// }}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "publish", Value: false}}}}
	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		panic(err)
	}

	// set new data.
	if val, ok := data["content"]; ok {
		data := val.([]domain.ComponentData)
		var ui []interface{}
		for _, t := range data {
			ui = append(ui, t)
		}
		_, err := collection.InsertMany(ctx, ui)
		if err != nil {
			return result, err
		}

	}

	err = r.db.Collection(tblPage).FindOne(ctx, bson.M{"_id": pageID}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
