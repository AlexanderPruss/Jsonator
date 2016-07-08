package Jsonator

import (
	"github.com/Jeffail/gabs"
	"reflect"
	"errors"
)

//Unmarshals a JSON to the input pointer v using the "jsonator" tag.
func DefaultUnmarshal(data []byte, v interface{}) error {
	return Unmarshal(data, v, "")
}

//Unmarshals a JSON to the input pointer v using the "jsonator" tag, but with pretty indentation.
func Unmarshal(data []byte, v interface{}, tagId string) error {
	val := reflect.ValueOf(v)
	if(val.Kind()!=reflect.Ptr){
		return errors.New("Unmarshalling requires the target interface value to be a pointer.")
	}

	container, err := gabs.ParseJSON(data)
	if err != nil {
		return err
	}
	return unmarshal(val.Elem(), container, tagId)
}

//Iterates through the fields of a struct, marshalling them in turn. Recurses into any struct fields found.
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


//Sets the value of the JSON Value to the reflection value. Potentially requires recursing into subvalues.
func setValue(val reflect.Value, jsonValue interface{}, container *gabs.Container, tagId string, currentPath ...string) {
	//TODO: Ugly interface. Maybe hide container and tagId into an object and turns this into a method?
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

//Helper method for setting collection values, as collection reflection is somewhat complicated in GO, even before custom tags get involved.
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
		if mapKeyType != stringType { //we can only map maps with string keys //TODO: Until we get more clever with tags!
			break
		}

		mapType := reflect.MapOf(mapKeyType, mapValueType)
		mapValue := reflect.MakeMap(mapType)

		for key, value := range jsonValue.(map[string]interface{}) {
			mapKeyValue := reflect.ValueOf(key)
			mapValuePointer := reflect.New(mapValueType)
			setValue(mapValuePointer.Elem(), value, container, tagId, currentPath...)
			mapValue.SetMapIndex(mapKeyValue, mapValuePointer.Elem())
		}
		val.Set(mapValue)

	case reflect.Slice:
		sliceVal := reflect.MakeSlice(val.Type(), 0, 0)
		for _, value := range jsonValue.([]interface{}) {
			sliceVal = reflect.Append(sliceVal, reflect.ValueOf(value))
		}
		val.Set(sliceVal)
	}
}
