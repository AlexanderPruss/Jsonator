package main

import (
	"reflect"
	"github.com/Jeffail/gabs"
	_ "fmt"
)

func Marshal(v interface{}, tagId string) ([]byte, error) {

	container, err := marshal(v, gabs.New(), tagId)
	if err != nil{
		return nil, err
	}
	return container.Bytes(), err
}

func marshal(v interface{}, container *gabs.Container, tagId string, currentPath ...string) (*gabs.Container, error) {
	val := reflect.ValueOf(v)
	if val.Kind()==reflect.Ptr{
		val = val.Elem()
	}
	currentType := val.Type()

	for fieldIndex := 0; fieldIndex < currentType.NumField(); fieldIndex++{
		typeField := currentType.Field(fieldIndex)
		valueField := val.Field(fieldIndex)

		tag := parseTag(typeField, tagId)
		if(tag.skip || (tag.omitEmpty && isEmptyValue(valueField))){
			continue
		}
		fieldPath := tag.getFieldPath(currentPath, typeField.Name)

		//see if the field is a struct, struct pointer, or interface. If so, recurse into it
		//TODO: Prevent infinite loops. Perhaps have an encoder object that saves encoding state?
		isStruct := false;
		var structValAsInterface interface{}
		if valueField.Kind()==reflect.Interface && valueField.Elem().Kind() == reflect.Ptr && valueField.Elem().Elem().Kind() == reflect.Struct{
			isStruct = true
			structValAsInterface = valueField.Elem().Elem().Interface()
		} else if valueField.Kind()==reflect.Ptr && valueField.Elem().Kind()==reflect.Struct{
			isStruct = true
			structValAsInterface = valueField.Elem().Interface()
		} else if valueField.Kind()==reflect.Struct{
			isStruct = true
			structValAsInterface = valueField.Interface()
		}
		if isStruct{
			_, err := marshal(structValAsInterface, container, tagId, fieldPath...)
			if(err!=nil){
				return nil, err
			}
		} else{
			_, err := container.Set(valueField.Interface(), fieldPath...)
			if(err!=nil){
				return nil, err
			}
		}
	}
	return container, nil
}

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