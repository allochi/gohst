package gohst

import "errors"

// var DataStore DataStoreDeligate
var datastores map[string]DataStore

func init() {
	datastores = make(map[string]DataStore)
}

type DataStore interface {
	PUT(interface{}) Response
	GET(interface{}, interface{}) Response
}

type Response struct {
	Message string
	Error   error
	Size    int
}

type Request struct {
	Ids    []int64
	Source string
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

func PUT(name string, object interface{}) Response {
	return datastores[name].PUT(object)
}

func GET(name string, object interface{}, query interface{}) Response {
	return datastores[name].GET(object, query)
}
