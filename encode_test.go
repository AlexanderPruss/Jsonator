package Jsonator

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//Not a classic go test, tests cases instead of individual functions.

func marshalAndUnmarshal(t *testing.T, v interface{}, tagId string) *gabs.Container {
	var bytes []byte
	var err error
	if tagId == "" {
		bytes, err = DefaultMarshal(v)
	} else {
		bytes, err = Marshal(v, tagId)
	}
	require.Nil(t, err, "Failed to marshal: %v", err)

	container, err := gabs.ParseJSON(bytes)
	require.Nil(t, err, "Failed to unmarshal: %v", err)

	return container
}

func TestDefaultMarshal(t *testing.T) {
	type testStruct struct {
		Value      string `jsonator:"nonexistentMap.value"`
		OtherValue int    `jsonator:"renamedValue"`
	}

	foo := testStruct{"value", 10}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.EqualValues(t, foo.OtherValue, container.Search("renamedValue").Data())
}

func TestMarshal_EncodingJsonTagsIgnored(t *testing.T) {
	type testStruct struct {
		Value      string `json:",-"`
		OtherValue int    `json:"renamedValue"`
	}
	foo := testStruct{"value", 10}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("Value").Data())
	assert.EqualValues(t, foo.OtherValue, container.Search("OtherValue").Data())
}

func TestMarshal_multipleJsonatorMappings(t *testing.T) {
	type testStruct struct {
		Value string `jsonator_abc:"abc" jsonator_foo:"foo"`
	}
	foo := testStruct{"value"}
	containerAbc := marshalAndUnmarshal(t, foo, "abc")
	containerFoo := marshalAndUnmarshal(t, foo, "foo")

	assert.Equal(t, foo.Value, containerAbc.Search("abc").Data())
	assert.Equal(t, foo.Value, containerFoo.Search("foo").Data())
}

func TestMarshal_ignoreField(t *testing.T) {
	type testStruct struct {
		Value        string `jsonator:"nonexistentMap.value"`
		IgnoredValue int    `jsonator:",-"`
	}

	foo := testStruct{"value", 10}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.Nil(t, container.Search("IgnoredValue").Data())
}

func TestMarshal_omitEmptyField(t *testing.T) {
	type testStruct struct {
		Value      string `jsonator:"nonexistentMap.value"`
		EmptyValue int    `jsonator:",omitempty"`
	}

	foo := testStruct{}
	foo.Value = "Value"
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.Nil(t, container.Search("EmptyValue").Data())
}

func TestMarshal_recursiveStruct(t *testing.T) {
	type innerStruct struct {
		InnerStructValue string `jsonator:"newInnerValue"`
	}

	type testStruct struct {
		Value       string      `jsonator:"nonexistentMap.value"`
		InnerStruct innerStruct `jsonator:"newStructName"`
	}

	foo := testStruct{"value", innerStruct{"innerValue"}}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.Equal(t, foo.InnerStruct.InnerStructValue, container.Search("newStructName", "newInnerValue").Data())
}

func TestMarshal_recursiveStructPointer(t *testing.T) {
	type innerStruct struct {
		InnerStructValue string `jsonator:"newInnerValue"`
	}

	type testStruct struct {
		Value       string       `jsonator:"nonexistentMap.value"`
		InnerStruct *innerStruct `jsonator:"newStructName"`
	}

	foo := testStruct{"value", &innerStruct{"innerValue"}}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.Equal(t, foo.InnerStruct.InnerStructValue, container.Search("newStructName", "newInnerValue").Data())
}

type stringerStruct struct {
	InnerStructValue string `jsonator:"newInnerValue"`
}

func (s stringerStruct) String() string {
	return "stringerStruct is a stringer."
}

func TestMarshal_recursiveInterfaceStruct(t *testing.T) {

	type testStruct struct {
		Value       string       `jsonator:"nonexistentMap.value"`
		InnerStruct fmt.Stringer `jsonator:"newStructName"`
	}

	foo := testStruct{"value", stringerStruct{"innerValue"}}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.Equal(t, "innerValue", container.Search("newStructName", "newInnerValue").Data())
}

func TestMarshal_mapIntoPreviousStruct(t *testing.T) {
	type innerStruct struct {
		InnerStructValue string `jsonator:"[previous].newInnerValue"`
	}

	type testStruct struct {
		Value       string `jsonator:"nonexistentMap.value"`
		InnerStruct *innerStruct
	}

	foo := testStruct{"value", &innerStruct{"innerValue"}}
	container := marshalAndUnmarshal(t, foo, "")

	assert.Equal(t, foo.Value, container.Search("nonexistentMap", "value").Data())
	assert.Equal(t, foo.InnerStruct.InnerStructValue, container.Search("newInnerValue").Data())
}
