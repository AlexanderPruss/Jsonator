package main

import (
	"fmt"
)

func main() {

	fmt.Println("Marshaling 'foo' with my own marshaller, hot damn")
	foo := &Foo{"foo",1, &Bar{"shouldStillBeMarshalled",3, "came from bar"}}
	jsonBytes, err := Marshal(foo,"map")
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBytes))
	fmt.Println("Marshaling 'foo' with the abc jsonator key")
	jsonBytesAbc, err := Marshal(foo,"abc")
	if err!=nil{
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

type Foo struct{
	Value string `jsonator_map:"nonexistentMap.value" jsonator_abc:"renamedValue"`
	NotTheValue int
	Blah fmt.Stringer
}

func (f *Foo) String() string{
	return fmt.Sprintf("Value: %v\nNotTheValue: %v", f.Value, f.NotTheValue)
}



type Bar struct{
	BarThing string `json:"-"`
	BizarroInt int `jsonator_map:"aMapInsideOfBar.bizarroIntLandedHere"`
	GoesIntoNonexistentMapOfFoo string `jsonator_map:"[previous].nonexistentMap.thisCameFromBar"`
}

func (b Bar) String() string{
	return fmt.Sprintf("Bar Thing: %v\n", b.BarThing)
}