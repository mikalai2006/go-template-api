package app

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BindAndValid binds and validates data.
func BindAndValid[V any](c *gin.Context, form V) (interface{}, error) {
	result := bson.M{}

	var body map[string]interface{}
	if er := json.NewDecoder(c.Request.Body).Decode(&body); er != nil {
		return result, er
	}

	var tagValue, primitiveValue string
	myDataReflect := reflect.Indirect(reflect.ValueOf(form))

	for i := 0; i < myDataReflect.NumField(); i++ {
		typeField := myDataReflect.Type().Field(i)
		tag := typeField.Tag
		tagValue = tag.Get("bson")
		primitiveValue = tag.Get("primitive")

		if val, ok := body[tagValue]; ok {
			switch myDataReflect.Field(i).Kind() {
			case reflect.String:
				result[tagValue] = val.(string)

			case reflect.Bool:
				result[tagValue] = val.(bool)

			case reflect.Int:
				result[tagValue] = val
			default:
				if primitiveValue == "true" {
					id, err := primitive.ObjectIDFromHex(val.(string))
					if err != nil {
						// todo error
						return result, err
					}
					result[tagValue] = id
				}
				// value := myDataReflect.Field(i)
				// fmt.Println("   === default: tag=", tagValue, value)
				// fmt.Println("   === default: value=", value)
				// fmt.Println("   === default: tag primitiveValue=", primitiveValue)
				// fmt.Println("   === default: kind= ", myDataReflect.Field(i).Kind())
			}
		}
	}

	// fmt.Println("============result======================")
	// fmt.Println(result)
	// fmt.Println("==========================================")

	result["updated_at"] = time.Now()
	return result, nil
}
