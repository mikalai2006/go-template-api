package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BindAndValid binds and validates data.
func BindAndValid[V any](c *gin.Context, form V) (interface{}, error) {
	var body map[string]interface{}
	if er := json.NewDecoder(c.Request.Body).Decode(&body); er != nil {
		return nil, er
	}
	result := make(map[string]interface{}, len(body))
	var tagValue, primitiveValue, tagJSONValue string
	myDataReflect := reflect.Indirect(reflect.ValueOf(form))

	for i := 0; i < myDataReflect.NumField(); i++ {
		typeField := myDataReflect.Type().Field(i)
		tag := typeField.Tag
		tagValue = tag.Get("bson")
		tagJSONValue = tag.Get("json")
		primitiveValue = tag.Get("primitive")
		if val, ok := body[tagJSONValue]; ok {
			// fmt.Println(tagValue, tagJSONValue, reflect.TypeOf(val))
			switch myDataReflect.Field(i).Kind() {
			case reflect.String:
				result[tagValue] = val.(string)

			case reflect.Bool:
				result[tagValue] = val.(bool)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// s := val.(string)
				// i, err := strconv.ParseInt(s, 10, 64)
				// if err == nil {
				// 	result[tagValue] = i
				// 	continue
				// }
				// f, err := strconv.ParseFloat(s, 64)
				// if err == nil {
				// 	result[tagValue] = f
				// 	continue
				// }
				result[tagValue] = val
			default:
				if primitiveValue == "true" {
					if reflect.ValueOf(val).Kind() == reflect.Slice {
						l := len(val.([]interface{}))
						idsPrimititiveSlice := make([]primitive.ObjectID, l)
						allValue := val.([]interface{})
						for i := range val.([]interface{}) {
							id, err := primitive.ObjectIDFromHex(allValue[i].(string))
							if err != nil {
								// todo error
								return result, err
							}
							idsPrimititiveSlice[i] = id
						}
						result[tagValue] = idsPrimititiveSlice
						// fmt.Println("default: ", tagValue, reflect.ValueOf(val).Kind())
					} else {

						id, err := primitive.ObjectIDFromHex(val.(string))
						if err != nil {
							// todo error
							return result, err
						}
						// fmt.Println("default: ", tagValue, reflect.ValueOf(val).Kind())
						result[tagValue] = id
					}
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

func BindJSON[V any](data map[string]interface{}) (V, error) {
	var result V
	parsedData := make(map[string]map[string]string)

	// fmt.Println(reflect.TypeOf(result))
	// fmt.Println(reflect.ValueOf(result))
	// fmt.Println(reflect.ValueOf(result).Kind())

	// fmt.Println(reflect.TypeOf(data))
	// fmt.Println(reflect.ValueOf(data))
	// fmt.Println(reflect.ValueOf(data).Kind())

	for i := range data {
		key := strings.Split(i, "__i18n__")
		if len(key) == 2 {
			if parsedData[key[1]] == nil {
				parsedData[key[1]] = map[string]string{}
			}
			if reflect.ValueOf(data[i]).Kind() != reflect.String {
				return result, fmt.Errorf("field %s must be string", i)
			}
			parsedData[key[1]][key[0]] = data[i].(string)
		}
	}
	// for k, v := range parsedData {
	data["locale"] = parsedData
	// }

	// fmt.Println("==========data======================")
	// fmt.Println(data)
	// fmt.Println("==========/data======================")

	// elementsStructure := reflect.ValueOf(result)
	var tagValue, primitiveValue string
	structValue := reflect.ValueOf(&result).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		typeField := structValue.Type().Field(i)
		tag := typeField.Tag
		tagValue = tag.Get("json")
		primitiveValue = tag.Get("primitive")
		// structValue2 := structValue.FieldByName(tagValue)
		// fmt.Println(tagValue, "= structValue2= ", structValue2)
		if !structValue.Field(i).CanSet() {
			return result, fmt.Errorf("no canset field %s", tagValue)
			// fmt.Println("===== nocanset", tagValue)
			// fmt.Println(i, tagValue, structValue.Field(i), structValue.Field(i).Type())
			// continue
		}
		val, ok := data[tagValue]
		if ok {

			if reflect.TypeOf(val) != structValue.Field(i).Type() {
				return result, fmt.Errorf("field %s must be type %s", tagValue, structValue.Field(i).Type().String())
				// fmt.Println("===== no compare types", tagValue)
				// fmt.Println(i, structValue.Field(i), structValue.Field(i).Type())
				// fmt.Println(i, val, reflect.TypeOf(val))
				// continue
			}

			// fmt.Println("====ok============================")
			// fmt.Println(i, tagValue, structValue.Field(i), structValue.Field(i).Type())
			valStructure := reflect.ValueOf(val)
			// if tagValue == "locale" {
			if primitiveValue == "true" {
				_, err := primitive.ObjectIDFromHex(val.(string))
				if err != nil {
					return result, err
				}
				// structValue.Field(i).Set(id)
				structValue.Field(i).Set(valStructure)
			} else {
				structValue.Field(i).Set(valStructure)
			}
			// }
		}
	}

	// fmt.Println("===result====")
	// fmt.Printf("%v", result)
	// fmt.Println("=============")

	return result, nil
}
