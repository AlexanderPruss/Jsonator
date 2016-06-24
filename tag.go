package main

import (
	"reflect"
	"strings"
)

const TAG_KEY = "jsonator"
const KEY_SEPARATOR = "_"
const TAG_SEPARATOR = ","
const PATH_SEPARATOR = "."

const OMIT_EMPTY = "omitempty"
const SKIP = "-"
const PREVIOUS = "[previous]"


type tag struct{
	path string
	omitEmpty bool
	skip bool
}

func parseTag(field reflect.StructField, tagId string) tag{
	t := tag{}
	tags := field.Tag.Get(getTagKey(tagId))
	if tags!=""{
		splitTags := strings.Split(tags, TAG_SEPARATOR)
		t.path = splitTags[0]
		for i := 1; i<len(splitTags); i++{
			currentTag := splitTags[i]
			if currentTag==OMIT_EMPTY{
				t.omitEmpty=true
			} else if currentTag == SKIP{
				t.skip = true
			}
		}
	}
	return t
}

func getTagKey(tagId string) string{
	return TAG_KEY + KEY_SEPARATOR +tagId
}

func (t tag) getFieldPath(currentPath []string, defaultPath string) []string{
	fieldPath := make([]string, 0)
	fieldPath = append(fieldPath, currentPath...)
	if t.path == ""{
		return append(fieldPath, defaultPath)
	}

	splitPath := strings.Split(t.path, PATH_SEPARATOR)
	escapes := 0
	for escapes = 0; escapes < len(splitPath) && splitPath[escapes]==PREVIOUS; escapes++{
	}

	if escapes == 0{
		return append(fieldPath, splitPath...)
	} else{
		//remove escapes from path
		splitPath = splitPath[escapes:]
		if escapes > len(fieldPath){
			return splitPath
		} else{
			fieldPath = fieldPath[:len(fieldPath)-escapes]
			return append(fieldPath, splitPath...)
		}
	}
}