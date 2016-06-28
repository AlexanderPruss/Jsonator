package Jsonator

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetTagKey(t *testing.T) {
	assert.Equal(t, TAG_KEY, getTagKey(""), "An empty tag key should return the default 'jsonator' tag.")
	assert.Equal(t, TAG_KEY+KEY_SEPARATOR+"foo", getTagKey("foo"), "An non-empty tag key should be preceded by 'jsonator_'.")
}

func TestParseTag(t *testing.T) {
	type testStruct struct {
		NoTags           string
		JsonTags         string `json:"jsonTag"`
		SingleTag        string `jsonator:"newStructName"`
		OmitEmptyTag     string `jsonator:",omitempty"`
		OmitEmptyAndPath string `jsonator:"omit.path,omitempty"`
		IgnoreTag        string `jsonator:",-"`
		MultipleKeyTags  string `jsonator_abc:"abc" jsonator_foo:"foo"`
	}

	foo := testStruct{}
	currentType := reflect.ValueOf(foo).Type()
	assert.Equal(t, tag{}, parseTag(currentType.Field(0), ""))
	assert.Equal(t, tag{}, parseTag(currentType.Field(1), ""))
	assert.Equal(t, tag{path: "newStructName", omitEmpty: false, skip: false}, parseTag(currentType.Field(2), ""))
	assert.Equal(t, tag{path: "", omitEmpty: true, skip: false}, parseTag(currentType.Field(3), ""))
	assert.Equal(t, tag{path: "omit.path", omitEmpty: true, skip: false}, parseTag(currentType.Field(4), ""))
	assert.Equal(t, tag{path: "", omitEmpty: false, skip: true}, parseTag(currentType.Field(5), ""))
	assert.Equal(t, tag{path: "abc", omitEmpty: false, skip: false}, parseTag(currentType.Field(6), "abc"))
	assert.Equal(t, tag{path: "foo", omitEmpty: false, skip: false}, parseTag(currentType.Field(6), "foo"))
}

func TestGetFieldPath(t *testing.T) {
	defaultPath := "defaultPath"
	firstPathElement := "firstPathElement"
	secondPathElement := "secondPathElement"
	emptyPath := make([]string, 0)
	nonemptyPath := make([]string, 0)
	nonemptyPath = append(nonemptyPath, firstPathElement)
	nonemptyPath = append(nonemptyPath, secondPathElement)

	//tags without a user-defined path
	tagWithoutPath := tag{}

	fieldPath := tagWithoutPath.getFieldPath(nonemptyPath, defaultPath)
	assert.Equal(t, []string{firstPathElement, secondPathElement, defaultPath}, fieldPath)

	fieldPath = tagWithoutPath.getFieldPath(emptyPath, defaultPath)
	assert.Equal(t, []string{defaultPath}, fieldPath)

	//tags with a user-defined path
	firstTagPath := "firstTagPath"
	secondTagPath := "secondTagPath"
	tagWithPath := tag{firstTagPath + PATH_SEPARATOR + secondTagPath, false, false}

	fieldPath = tagWithPath.getFieldPath(nonemptyPath, defaultPath)
	assert.Equal(t, []string{firstPathElement, secondPathElement, firstTagPath, secondTagPath}, fieldPath)

	fieldPath = tagWithPath.getFieldPath(emptyPath, defaultPath)
	assert.Equal(t, []string{firstTagPath, secondTagPath}, fieldPath)

	//tags with [previous] tags.
	tagWithPreviousPath := tag{PREVIOUS + PATH_SEPARATOR + firstTagPath + PATH_SEPARATOR + secondTagPath, false, false}

	fieldPath = tagWithPreviousPath.getFieldPath(nonemptyPath, defaultPath)
	assert.Equal(t, []string{firstPathElement, firstTagPath, secondTagPath}, fieldPath)

	fieldPath = tagWithPreviousPath.getFieldPath(emptyPath, defaultPath)
	assert.Equal(t, []string{firstTagPath, secondTagPath}, fieldPath)

	tagWithDoublePreviousPath := tag{PREVIOUS + PATH_SEPARATOR + PREVIOUS + PATH_SEPARATOR + firstTagPath + PATH_SEPARATOR + secondTagPath, false, false}

	fieldPath = tagWithDoublePreviousPath.getFieldPath(nonemptyPath, defaultPath)
	assert.Equal(t, []string{firstTagPath, secondTagPath}, fieldPath)

	fieldPath = tagWithDoublePreviousPath.getFieldPath(emptyPath, defaultPath)
	assert.Equal(t, []string{firstTagPath, secondTagPath}, fieldPath)
}
