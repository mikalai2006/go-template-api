package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
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
func GetParamsFromRequest[V any](c *gin.Context, filterStruct V) (domain.RequestParams, error) {
	params := domain.RequestParams{
		Filter: filterStruct,
	}
	var filter V
	if err := c.ShouldBind(&filter); err != nil {
		// disable error for convert string to primitive.ObjectID.
		return domain.RequestParams{}, nil
	}

	filterValues := c.Request.URL.Query()
	dataFilter := bson.M{}
	var tagValue, primitiveValue, tagJsonValue string
	elementsFilter := reflect.ValueOf(filter)
	for i := 0; i < elementsFilter.NumField(); i++ {
		typeField := elementsFilter.Type().Field(i)
		tag := typeField.Tag

		tagValue = tag.Get("bson")
		primitiveValue = tag.Get("primitive")
		tagJsonValue = tag.Get("json")
		primitiveValue = tag.Get("primitive")
		fmt.Println(tagValue, tagJsonValue)

		if len(filterValues[tagJsonValue]) != 0 {
			switch elementsFilter.Field(i).Kind() {
			case reflect.String:
				value := elementsFilter.Field(i).String()
				dataFilter[tagValue] = value

			case reflect.Bool:
				value := elementsFilter.Field(i).Bool()
				dataFilter[tagValue] = value

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value := elementsFilter.Field(i).Int()
				dataFilter[tagValue] = value

			default:
				fmt.Println("default ", tagValue, elementsFilter.Field(i).Type(), primitiveValue)
				if primitiveValue == "true" {
					// fmt.Println("===== add ", tagValue)
					id, _ := primitive.ObjectIDFromHex(filterValues[tagValue][0])
					// if err != nil {
					// 	// todo error
					// }
					dataFilter[tagValue] = id
				}
			}
		}
	}

	var opts domain.Options
	if err := c.Bind(&opts); err != nil {
		return domain.RequestParams{}, err
	}

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
	// err = bson.Unmarshal(sort, &sort)
	// fmt.Println("----------")
	// fmt.Printf("dataFilter=%s", dataFilter)
	// fmt.Println("----------")
	// fmt.Printf("len dataFilter=%s", len(dataFilter))
	// fmt.Println("----------")
	// fmt.Printf("filter=%s", filter)
	// fmt.Println("----------")
	// fmt.Printf("sort=%s", testBson)
	// fmt.Println("----------")
	// fmt.Printf("opts=%s", opts)
	// fmt.Println("----------")
	if opts.Limit == 0 || opts.Limit > 50 {
		opts.Limit = 10
	}
	params.Filter = dataFilter
	params.Options = opts

	fmt.Println(params)
	return params, nil
}
