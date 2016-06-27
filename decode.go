package Jsonator

import "github.com/Jeffail/gabs"

//not yet implemented
func Unmarshal(data []byte, v interface{}, tagId string) error {
	//check that we're sending in a pointer? Or just let it crash when the value isn't settable?

	container, err := gabs.ParseJSON(data)
	if err != nil {
		return err
	}
	return unmarshal(v, container, tagId)
}

//not yet implemented
func unmarshal(v interface{}, container *gabs.Container, tagId string, currentPath ...string) error {

	return nil
}
