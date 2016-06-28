package Jsonator

import (
	"github.com/Jeffail/gabs"
	"reflect"
)

//Marshals an object to JSON using the default "jsonator" tag.
func DefaultMarshal(v interface{}) ([]byte, error) {

	return Marshal(v, "")
}

//Marshals an object to JSON using the "jsonator_tagId" tag.
func Marshal(v interface{}, tagId string) ([]byte, error) {

	container, err := marshal(v, gabs.New(), tagId)
	if err != nil {
		return nil, err
	}
	return container.Bytes(), err
}

//Iterates through the fields of a struct. Recurses into any struct fields found.
func marshal(v interface{}, container *gabs.Container, tagId string, currentPath ...string) (*gabs.Container, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	currentType := val.Type()

	for fieldIndex := 0; fieldIndex < currentType.NumField(); fieldIndex++ {
		typeField := currentType.Field(fieldIndex)
		valueField := val.Field(fieldIndex)

		tag := parseTag(typeField, tagId)
		if tag.skip || (tag.omitEmpty && isEmptyValue(valueField)) {
			continue
		}
		fieldPath := tag.getFieldPath(currentPath, typeField.Name)

		//see if the field is a struct, struct pointer, or interface. If so, recurse into it
		//TODO: Prevent infinite loops. Perhaps have an encoder object that saves encoding state?
		isStruct := false
		var structValAsInterface interface{}

		if valueField.Kind() == reflect.Interface {
			isStruct = true
			if valueField.Elem().Kind() == reflect.Struct {
				structValAsInterface = valueField.Elem().Interface()
			} else {
				structValAsInterface = valueField.Elem().Elem().Interface()
			}
		} else if valueField.Kind() == reflect.Ptr && valueField.Elem().Kind() == reflect.Struct {
			isStruct = true
			structValAsInterface = valueField.Elem().Interface()
		} else if valueField.Kind() == reflect.Struct {
			isStruct = true
			structValAsInterface = valueField.Interface()
		}
		if isStruct {
			_, err := marshal(structValAsInterface, container, tagId, fieldPath...)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := container.Set(valueField.Interface(), fieldPath...)
			if err != nil {
				return nil, err
			}
		}
	}
	return container, nil
}

//Checks for empty values. Weird, but that's how reflection works in Go I guess
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
