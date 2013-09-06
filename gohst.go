package gohst

import "errors"
import "fmt"

// var DataStore DataStoreDeligate
var datastores map[string]DataStore

func init() {
	datastores = make(map[string]DataStore)
}

type DataStore interface {
	PUT(interface{}) error
	GET(interface{}, interface{}) error
}

func Register(name string, datastore DataStore) error {
	if name == "" {
		return errors.New("Can't have empty data store name")
	}

	if datastore == nil {
		return errors.New("Can't have empty data store")
	}

	if datastores[name] != nil {
		return errors.New("Can't reassign data store name")
	}

	datastores[name] = datastore
	return nil
}

func PUT(name string, object interface{}) error {
	if !IsStructOrPtr2Struct(object) {
		return errors.New("PUT only accepts an object or a pointer to an object of type struct")
	}
	return datastores[name].PUT(object)
}

func GET(name string, object interface{}, ids interface{}) error {
	fmt.Println("Testing array")
	if !IsPtr2SliceOfStruct(object) {
		return errors.New("GET accepts a pointer to slice of a struct type as an object")
	}
	fmt.Println("Testing ids")
	if ids != nil && !IsSliceOrPtr2SliceOfPrimitive(ids) {
		return errors.New("GET accepts pointer to slice of a primitive type as ids e.g int64 or string")
	}
	fmt.Println("Calling Real GET")
	return datastores[name].GET(object, ids)
}
