package Jsonator

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"reflect"
)

//not yet implemented
func DefaultUnmarshal(data []byte, v interface{}) error {
	return Unmarshal(data, v, "")
}

//not yet implemented
func Unmarshal(data []byte, v interface{}, tagId string) error {
	//TODO: check that we're sending in a pointer? Or just let it crash when the value isn't settable?
	//TODO: What does encoding/json do?

	container, err := gabs.ParseJSON(data)
	if err != nil {
		return err
	}
	return unmarshal(reflect.ValueOf(v).Elem(), container, tagId)
}

//not yet implemented
func unmarshal(val reflect.Value, container *gabs.Container, tagId string, currentPath ...string) error {
	currentType := val.Type()
	for fieldIndex := 0; fieldIndex < currentType.NumField(); fieldIndex++ {
		typeField := currentType.Field(fieldIndex)
		valueField := val.Field(fieldIndex)

		tag := parseTag(typeField, tagId)
		if tag.skip || (tag.omitEmpty && isEmptyValue(valueField)) {
			continue
		}
		fieldPath := tag.getFieldPath(currentPath, typeField.Name)
		jsonValue := container.Search(fieldPath...).Data()

		if jsonValue == nil {
			continue
		}
		setValue(valueField, jsonValue, container, tagId, fieldPath...)
	}
	return nil
}

//TODO: Ugly interface. Maybe hide container and tagId into an object and turns this into a method?
func setValue(val reflect.Value, jsonValue interface{}, container *gabs.Container, tagId string, currentPath ...string) {
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		setCollectionValue(val, jsonValue, container, tagId, currentPath...)
	case reflect.String:
		val.SetString(jsonValue.(string))
	case reflect.Bool:
		val.SetBool(jsonValue.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val.SetInt(int64(jsonValue.(float64)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val.SetUint(uint64(jsonValue.(float64)))
	case reflect.Float32, reflect.Float64:
		val.SetFloat(jsonValue.(float64))
	case reflect.Ptr:
		newPtrVal := reflect.New(val.Type().Elem())
		setValue(newPtrVal.Elem(), jsonValue, container, tagId, currentPath...)
		val.Set(newPtrVal)
	case reflect.Interface:
		//Not knowing anything about what concrete type to map to, the best we can do is write the
		//value from jsonVal in directly. This will probably only work if val is an empty interface.
		jsonVal := reflect.ValueOf(jsonValue)
		if jsonVal.Type().AssignableTo(val.Type()) {
			val.Set(jsonVal)
		}
	case reflect.Struct:
		unmarshal(val, container, tagId, currentPath...)
	}
}

func setCollectionValue(val reflect.Value, jsonValue interface{}, container *gabs.Container, tagId string, currentPath ...string) {
	switch val.Kind() {
	case reflect.Array:
		arrayPointerVal := reflect.New(reflect.ArrayOf(val.Type().Len(), val.Type().Elem()))
		for i, value := range jsonValue.([]interface{}) {
			setValue(arrayPointerVal.Elem().Index(i), value, container, tagId, currentPath...)
		}
		val.Set(arrayPointerVal.Elem())
	case reflect.Map:
		mapValueType := val.Type().Elem()
		mapKeyType := val.Type().Key()
		stringType := reflect.ValueOf("").Type()
		if mapKeyType != stringType { //we can only map maps with string keys
			break
		}

		mapType := reflect.MapOf(mapKeyType, mapValueType)
		mapPointer := reflect.New(mapType)
		mapValue := mapPointer.Elem()
		fmt.Printf("map value type:%v\n", mapValue.Type)
		for key, value := range jsonValue.(map[string]interface{}) {
			mapKeyValue := reflect.ValueOf(key)
			mapValuePointer := reflect.New(mapValueType)
			setValue(mapValuePointer, value, container, tagId, currentPath...) //TODO: Unadressable values X_X
			val.SetMapIndex(mapKeyValue, mapValuePointer.Elem())
		}

	case reflect.Slice:
		sliceVal := reflect.MakeSlice(val.Type(), 0, 0)
		for _, value := range jsonValue.([]interface{}) {
			sliceVal = reflect.Append(sliceVal, reflect.ValueOf(value))
		}
		val.Set(sliceVal)
	}
}
