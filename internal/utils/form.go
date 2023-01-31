package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BindAndValid binds and validates data.
func BindAndValid[V any](c *gin.Context, form V) (interface{}, error) {
	var body map[string]interface{}
	if c.Request.Body == nil {
		return nil, fmt.Errorf("not found data for patch")
	}
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
				} else {
					result[tagValue] = val
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

// BindAndValid component data for page.
func BindPageWithContent[V any](c *gin.Context, form V) (map[string]interface{}, error) {
	var body map[string]interface{}
	if er := json.NewDecoder(c.Request.Body).Decode(&body); er != nil {
		return nil, er
	}
	pageID := c.Param("id")
	var typeContent string
	if typeC, ok := body["type"]; ok {
		typeContent = typeC.(string)
	}
	var layoutID string
	if id, ok := body["layoutId"]; ok {
		layoutID = id.(string)
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
				if tagValue == "content" {
					result[tagValue] = buildFlatDataFormTree(val.(map[string]interface{}), layoutID, pageID, typeContent)
					// value := myDataReflect.Field(i)
					// fmt.Println("   === default: tag=", tagValue, value)
					// fmt.Println("   === default: value=", value)
					// fmt.Println("   === default: tag primitiveValue=", primitiveValue)
					// fmt.Println("   === default: kind= ", myDataReflect.Field(i).Kind())
				}
			}
		}
	}

	result["updated_at"] = time.Now()
	return result, nil
}

type StackNode struct {
	Node   interface{}
	Parent string
	Global bool
}

func buildFlatDataFormTree(
	tree interface{},
	layoutID string,
	pageID string,
	typeContent string,
	// relations map[string][]*domain.Field,
	// level int,
	// i18n config.I18nConfig,
) []domain.ComponentData {
	// level++
	var result []domain.ComponentData
	var stack []StackNode
	stack = append(stack, StackNode{
		Node:   tree,
		Parent: "root",
		// Global: false,
	})
	PID, _ := primitive.ObjectIDFromHex(pageID)
	// LID, _ := primitive.ObjectIDFromHex(layoutID)

	for len(stack) > 0 {
		n := len(stack) - 1 // Top element
		currentNode := stack[n]
		stack = stack[:n] // Pop

		val := reflect.ValueOf(currentNode.Node)

		// if valData.Kind() == reflect.Struct {
		// 	fmt.Println("val type: ", valData.Interface().(map[string]interface{}))
		// 	valInt := reflect.ValueOf(currentNode.Node).Elem()
		// 	val = reflect.ValueOf(valInt)
		// 	// copy := reflect.New(val.Type()).Elem()
		// 	// for i := 0; i < val.NumField(); i += 1 {
		// 	// 	copy.Field(i).Set(val.Field(i))
		// 	// 	fmt.Println("val.Field(i)=", val.Field(i))
		// 	// }
		// 	// val = copy
		// }
		// fmt.Println("currentNode.Node->", currentNode, " kind=", val.Kind())

		if val.Kind() == reflect.Map {

			// create file css for page.
			if typeContent == "page" {
				cssValue := val.MapIndex(reflect.ValueOf("___cssPage"))
				if cssValue.IsValid() {
					filePath := fmt.Sprintf("./public/css/p_%s.css", pageID)
					f, err := os.Create(filePath)
					if err != nil {
						fmt.Println(err)
					}
					defer f.Close()

					f.Write([]byte(cssValue.Elem().String()))
				}
			}

			// create file css for layout.
			if typeContent == "layout" {
				cssLayoutValue := val.MapIndex(reflect.ValueOf("___cssLayout"))
				if cssLayoutValue.IsValid() {
					filePath := fmt.Sprintf("./public/css/l_%s.css", pageID)
					f, err := os.Create(filePath)
					if err != nil {
						fmt.Println(err)
					}
					defer f.Close()

					f.Write([]byte(cssLayoutValue.Elem().String()))
				}
			}

			globalValue := val.MapIndex(reflect.ValueOf("global"))
			var PPID primitive.ObjectID
			if globalValue.IsValid() {
				PPID = primitive.NilObjectID
			} else {
				PPID = PID
			}

			keyUID := val.MapIndex(reflect.ValueOf("_uid"))
			var parent string
			if keyUID.Elem().String() == layoutID {
				parent = "layout"
			} else if keyUID.Elem().String() == pageID {
				parent = "root"
			} else {
				parent = currentNode.Parent
			}

			res := domain.ComponentData{
				Parent: parent,
				UID:    keyUID.Elem().String(),
				PageID: PPID,
				// LayoutID: LID,
				// Component: "Component",
				Publish:   true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			resData := map[string]interface{}{}

			// fmt.Println("keyUID->", keyUID)

			for _, e := range val.MapKeys() {
				v := val.MapIndex(e)
				key := e.Interface().(string)
				if key == "___cssLayout" || key == "___cssPage" {
					// res.Component = v.Elem().String()
				} else if key == "component" {
					res.Component = v.Elem().String()
				} else if key == "_uid" {
					res.UID = v.Elem().String()
				} else {
					switch t := v.Interface().(type) {
					case int:
						// fmt.Println(e, "=", t)
						resData[key] = t
					case string:
						// fmt.Println(e, "=", t)
						resData[key] = t
					case bool:
						// fmt.Println(e, "=", t)
						resData[key] = t
					default:
						// fmt.Println("default:", e, "=", t, reflect.TypeOf(t).Kind())
						if reflect.TypeOf(t).Kind() == reflect.Slice {
							// reflect.ValueOf(t)

							s := reflect.ValueOf(t)
							// keysFirstElement := s.Index(0)
							// isNestedBloks := false
							// fmt.Println("keysFirstElement=", keysFirstElement, keysFirstElement.Type(), reflect.TypeOf(keysFirstElement).Kind())
							// if reflect.TypeOf(keysFirstElement).Kind() == reflect.Map {
							// 	for _, keyProperty := range keysFirstElement.MapKeys() {
							// 		keyProPertyAsString := keyProperty.Interface().(string)
							// 		fmt.Println("keysFirstElement=", keysFirstElement.MapIndex(keyProperty).String(), " | key", key, " == ", keyProPertyAsString)
							// 		if keyProPertyAsString == "_uid" {
							// 			isNestedBloks = true
							// 		}
							// 	}
							// }
							// fmt.Println("isNestedBloks=", isNestedBloks, " key", key)
							// if isNestedBloks {
							// if slice - is nested bloks.
							uids := []string{}

							for i := 0; i < s.Len(); i++ {
								// fmt.Println("s.Index(i).Elem()=", s.Index(i))
								// global := false
								// fmt.Println("global=", key, globalValue)
								// if key == "layout" {
								// 	global = true
								// }

								// add UID node to slice.
								valChild := s.Index(i).Elem()
								isBlok := false
								// fmt.Println("type =>", valChild.Kind(), reflect.TypeOf(valChild).Kind())
								if valChild.Kind() == reflect.Map {
									for _, ee := range valChild.MapKeys() {
										vChild := valChild.MapIndex(ee)
										keyChild := ee.Interface().(string)
										if keyChild == "_uid" {
											valUID := vChild.Elem().String()
											if valUID == pageID {
												valUID = "page"
											}
											uids = append(uids, valUID)
											isBlok = true
										}
									}

								}

								if isBlok == true {

									stack = append(stack, StackNode{
										Node:   s.Index(i).Elem().Interface(), // s.Index(i),
										Parent: keyUID.Elem().String(),
										// Global: global,
									})
								}
								// fmt.Println("========================")
								// fmt.Println(s.Index(i).Elem().Interface())
								// fmt.Println("========================")
							}
							// fmt.Println("=========", key, "===============", uids)
							if len(uids) > 0 {
								// if exist key by uids, insert slice uid.
								resData[key] = map[string]interface{}{
									"uids": uids,
								}
							} else {
								// if slice - is custom array.
								resData[key] = t
							}

							// } else {
							// 	// if slice - is custom array.
							// 	resData[key] = t
							// }
						} else {
							// fmt.Println(" custom object ===", key, t)
							// if value - is custom object.
							// write custom object as property.
							resData[key] = t
						}
					}
				}
			}
			res.Data = resData
			result = append(result, res)
		}
		// fmt.Println("stack->", len(stack), stack)
	}

	// for i, field := range fields {
	// 	node := map[string]interface{}{
	// 		"_uid":   field.UID,
	// 		"parent": field.Parent,
	// 		"name":   field.Name,
	// 		"level":  level,
	// 	}

	// 	for k, fieldData := range field.Data.Value {
	// 		switch c := fieldData.(type) {
	// 		case map[string]interface{}:
	// 			if _, ok := c[i18n.Default]; ok {
	// 				for i, data := range fieldData.(map[string]any) {
	// 					key := fmt.Sprintf("%s%s%s", k, i18n.Prefix, i)
	// 					node[key] = data
	// 					if i == i18n.Default {
	// 						node[k] = data
	// 					}
	// 				}
	// 				// delete(field.Data.Value, k)
	// 			}
	// 		default:
	// 			node[k] = fieldData
	// 		}
	// 	}

	// 	if childIDS, ok := relations[field.UID]; ok {
	// 		node["child"] = buildTree(childIDS, relations, level, i18n)
	// 	}
	// 	if node["parent"] == nil || node["level"].(int) != 1 {
	// 		tree[i] = node
	// 	}
	// }

	// fmt.Println("============result======================")
	// fmt.Println(result)
	// fmt.Println("==========================================")

	return result
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
