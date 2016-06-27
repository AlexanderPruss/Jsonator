// Copyright 2016 Alexander Pruss. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//Jsonator defines JSON mapping tags that expand upon Go's encoding/json package.
//In addition to having expanded functionality, fields can define multiple different
//jsonator tags, allowing a single model object to map to and from multiple json formats.
//
//The first part of each Jsonator tag is the path to which the field should be mapped. As in
//the encoding/json package, this field can be explicitly kept empty. The path to which
//a field is mapped can be a simple renaming of the field, or can be a period-separated list.
//Following the path, the "-" and "omitempty" tags from the encoding/json package can be added.
//
//Example:	type Foo struct{
//			RenamedValue string 	`jsonator:"newValueName"` 	//This value will be renamed to Foo.newValueName
//			OmittedValue int 	`jsonator:",-"`			//This value will not be mapped
//			ValueSavedIntoMap int	`jsonator:"newMap.value"`	//This map will be saved in a new JSON map under Foo.newMap.value
//			NotTagged string					//Not tagged, so mapped to Foo.NotTagged
// 			}
//
//Struct and interface fields of structs will be recursively iterated into. By using
//(possibly successive) [previous] path tags, fields of substructs can be mapped into
// their parent struct.
//
//Example:	type Foo struct{
//			MapsIntoBar string `jsonator:"bar.cameFromFoo"` 	//This value will be mapped into Foo.bar.cameFromFoo
//			Bar Bar
//		}
//		type Bar struct{
//			MapsIntoFoo string `jsonator:"[previous].cameFromBar"` 	//This value will be mapped into Foo.cameFromBar
//		}
//
//Multiple Json mappings can be defined by either mixing Jsonator tags with encoding/json tags
//or by appending a tag ID to each Jsonator tag. A jsonator tag with a tag id looks like
//jsonator_tagId:"tags"
//
//Example: 	type Foo struct{
//			TwoDifferentMappings string `jsonator_rename:"newName" jsonator_map:"newMap.TwoDifferentMappings"`
//			//When using tagId "rename", the string will be renamed to newName
//			//When using tagId "map", the string will be saved in a new map called newMap.
//		}
package Jsonator

import (
	"fmt"
)

//Proof of concept of marshalling. Unmarshalling coming soon.
func main() {

	fmt.Println("Marshaling 'foo' with my own marshaller, hot damn")
	foo := &Foo{"foo", 1, &Bar{"shouldStillBeMarshalled", 3, "came from bar"}}
	jsonBytes, err := Marshal(foo, "map")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))
	fmt.Println("Marshaling 'foo' with the abc jsonator key")
	jsonBytesAbc, err := Marshal(foo, "abc")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytesAbc))

	//holding on to this for reference
	//here's a stupid version of what we're going to be doing
	/*foo := &Foo{"foo",1}

	fmt.Println("\nUnmarshaling custom json 'foo' with gabs")
	parsedFoo, err := gabs.ParseJSON(mapFoo.Bytes())
	if err!=nil{
		fmt.Println(err)
		return
	}
	foo = &Foo{}
	foo.Value = parsedFoo.Search(mapTags...).Data().(string) //TODO: Not quite right yet
	fmt.Println(foo)*/
}

type Foo struct {
	Value       string `jsonator_map:"nonexistentMap.value" jsonator_abc:"renamedValue"`
	NotTheValue int
	Blah        fmt.Stringer
}

func (f *Foo) String() string {
	return fmt.Sprintf("Value: %v\nNotTheValue: %v", f.Value, f.NotTheValue)
}

type Bar struct {
	BarThing                    string `json:"-"`
	BizarroInt                  int    `jsonator_map:"aMapInsideOfBar.bizarroIntLandedHere"`
	GoesIntoNonexistentMapOfFoo string `jsonator_map:"[previous].nonexistentMap.thisCameFromBar"`
}

func (b Bar) String() string {
	return fmt.Sprintf("Bar Thing: %v\n", b.BarThing)
}
