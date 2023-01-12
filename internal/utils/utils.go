package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Utils interface {
	BodyToData() (interface{}, error)
	// ParamsToFilter() (interface{}, error)
}

// Create interface from body request to update item mongodb.
func GetBodyToData(u Utils) (interface{}, error) {
	data, err := u.BodyToData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Parse request params and return struct domain.RequestParams.
func GetParamsFromRequest[V any](c *gin.Context, filterStruct V, i18n *config.I18nConfig) (domain.RequestParams, error) {
	params := domain.RequestParams{
		Filter: filterStruct,
	}
	// set current locale.
	lang := c.Query("lang")
	if lang != "" {
		params.Lang = lang
	} else {
		params.Lang = i18n.Default
	}

	// fmt.Println("lang", c.QueryMap("name"))
	var filter V
	if err := c.ShouldBind(&filter); err != nil {
		// disable error for convert string to primitive.ObjectID.
		return domain.RequestParams{}, err
	}
	// fmt.Println("filter: ", filter)

	filterValues := c.Request.URL.Query()
	// fmt.Println("filterValues: ", filterValues)
	dataFilter := bson.M{}
	var tagValue, primitiveValue, tagJSONValue, tagMapQuery string
	elementsFilter := reflect.ValueOf(filter)
	for i := 0; i < elementsFilter.NumField(); i++ {
		tag := elementsFilter.Type().Field(i).Tag
		tagValue = tag.Get("bson")
		primitiveValue = tag.Get("primitive")
		tagJSONValue = tag.Get("json")
		tagMapQuery = fmt.Sprintf("%s[]", tag.Get("json"))

		valueParam := filterValues[tagJSONValue]
		// if len(valueParam) == 0 {
		// 	valueParam = filterValues[tagMapQuery]
		// }

		// fmt.Println(tagValue, tagJSONValue, tagMapQuery, valueParam)
		if len(valueParam) != 0 {
			// fmt.Println(tagValue, elementsFilter.Field(i).Kind())
			switch elementsFilter.Field(i).Kind() {
			case reflect.String:
				value := elementsFilter.Field(i).String()
				if primitiveValue == "true" {
					id, _ := primitive.ObjectIDFromHex(valueParam[0])
					// fmt.Println("===== string add ", tagValue, filterValues[tagJSONValue])
					dataFilter[tagValue] = id
				} else {
					dataFilter[tagValue] = value
				}

			case reflect.Bool:
				value := elementsFilter.Field(i).Bool()
				dataFilter[tagValue] = value

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value := elementsFilter.Field(i).Int()
				dataFilter[tagValue] = value

			default:

			}
		} else {
			valueParam = filterValues[tagMapQuery]
			if len(valueParam) != 0 {
				// fmt.Println("custom:", tagValue, valueParam)
				// fmt.Println("default ", tagValue, elementsFilter.Field(i).Type(), primitiveValue)
				if primitiveValue == "true" {
					var sliceIn bson.D
					a := []primitive.ObjectID{}
					for i := range valueParam {
						id, _ := primitive.ObjectIDFromHex(valueParam[i])
						// fmt.Println("===== add ", tagValue, id)
						a = append(a, id)
					}
					// fmt.Println("===== a ", a)
					sliceIn = bson.D{{"$in", a}}
					dataFilter[tagValue] = sliceIn
					// id, _ := primitive.ObjectIDFromHex(filterValues[tagValue][0])
					// // if err != nil {
					// // 	// todo error
					// // }
					// dataFilter[tagValue] = id
				} else {
					// fmt.Println(tagValue, tagJSONValue, tagMapQuery, valueParam)
					var sliceIn bson.D
					sliceIn = bson.D{{"$in", valueParam}}
					dataFilter[tagValue] = sliceIn
				}
			}
		}
	}

	// bind query params.
	var opts domain.Options
	// limit, err := strconv.ParseInt(c.Query("$limit"), 10, strconv.IntSize)
	// if err != nil {
	// 	return domain.RequestParams{}, err
	// }
	// fmt.Println("query > ", limit)
	if err := c.Bind(&opts); err != nil {
		return domain.RequestParams{}, err
	}
	// fmt.Println("Bind options", opts)

	sort := c.QueryMap("$sort")
	var testBson bson.D
	if len(sort) > 0 {
		for k, v := range sort {
			value, err := strconv.ParseInt(v, 10, strconv.IntSize)
			if err != nil {
				return domain.RequestParams{}, err
			}
			testBson = append(testBson, bson.E{Key: k, Value: value})
		}
		opts.Sort = testBson
	}
	// TODO opts.Limit.
	if opts.Limit == 0 || opts.Limit > 100 {
		opts.Limit = 10
	}
	params.Filter = dataFilter
	params.Options = opts
	// fmt.Println("query map: ", c.Request.URL.Query())
	// fmt.Println("params: ", params)
	return params, nil
}
