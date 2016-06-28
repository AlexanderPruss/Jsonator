package Jsonator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

//Not a classic go test, tests cases instead of individual functions.

func testUnmarshal(t *testing.T, data []byte, v interface{}, tagId string) {
	var err error
	if tagId == "" {
		err = DefaultUnmarshal(data, v)
	} else {
		err = Unmarshal(data, v, tagId)
	}
	require.Nil(t, err, "Failed to unmarshal: %v", err)
}

type innerStruct struct {
	InnerStructVal string
}

func (s innerStruct) String() string {
	return "inner value: " + s.InnerStructVal
}

func TestDefaultUnmarshal(t *testing.T) {

	//TODO: get maps to work, split this test up into smaller, more readables ones ala encode_test.go

	type testStruct struct {
		Value            *string `jsonator:"nonexistentMap.value"`
		OtherValue       int     `jsonator:"renamedValue"`
		ASlice           []string
		AnArray          [2]string
		AMap             map[string]interface{}
		InStruct         innerStruct
		AnInterface      interface{}
		WontMapInterface fmt.Stringer
	}

	json := `
	{
		"nonexistentMap":{"value":"value"},
		"renamedValue":10,
		"ASlice":[
		"foo",
		"bar"],
		"AnArray":[
		"array1",
		"array2"],
		"AMap":{"mapKey":"mapValue"},
		"InStruct":{"InnerStructVal":"innerValue"},
		"AnInterface":{"InnerStructVal":"interfaceValue"},
		"WontMapInterface":{"InnerStructVal":"wontMap"}
	}`

	unmarshalTarget := &testStruct{}
	testUnmarshal(t, []byte(json), unmarshalTarget, "")

	stringValue := "value"

	expectedValue := testStruct{&stringValue, 10, []string{"foo", "bar"},
		[2]string{"array1", "array2"}, map[string]interface{}{"mapKey": "mapValue"}, innerStruct{"innerValue"},
		map[string]interface{}{"InnerStructVal": "interfaceValue"},
		nil}
	assert.Equal(t, expectedValue, *unmarshalTarget)
}

//todo: test passing in a non-pointer
