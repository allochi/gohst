package gohst

import "errors"

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
	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return errors.New("PUT only accepts an object or a pointer to an object of type struct")
	}
	return datastores[name].PUT(object)
}

func GET(name string, object interface{}, ids interface{}) error {
	if name == "" {
		return errors.New("GET requires a data store name")
	}

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return errors.New("GET accepts a pointer to slice of a struct type as an object")
	}

	_idsKind := KindOf(ids)
	if _idsKind != SliceOfPrimitive {
		return errors.New("GET accepts slice of a primitive type as ids e.g int64 or string")
	}

	return datastores[name].GET(object, ids)
}
