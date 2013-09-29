package gohst

import "fmt"

type DataStore struct {
	container DataStoreContainer
}

var datastores map[string]DataStoreContainer

func init() {
	datastores = make(map[string]DataStoreContainer)
}

type DataStoreContainer interface {
	Connect() error
	Disconnect() error
	PUT(interface{}) error
	GET(interface{}, interface{}) error
	DELETE(interface{}, interface{}) error
}

func (ds *DataStore) Register(name string, container DataStoreContainer) error {
	if name == "" {
		return fmt.Errorf("Can't have empty data store name")
	}

	if container == nil {
		return fmt.Errorf("Can't have empty data store")
	}

	if datastores[name] != nil {
		return fmt.Errorf("Can't reassign data store name")
	}

	ds.container = container
	datastores[name] = container
	return nil
}

func GetDataStore(name string) (ds DataStore, err error) {
	if name == "" {
		err = fmt.Errorf("Can't have empty data store name")
	}

	dsc := datastores[name]

	if dsc == nil {
		err = fmt.Errorf("Couldn't find the data store")
	} else {
		ds = DataStore{dsc}
	}

	return ds, err
}

func (ds *DataStore) PUT(object interface{}) error {
	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("PUT only accepts an object or a pointer to an object of type struct")
	}
	return ds.container.PUT(object)
}

func (ds *DataStore) GET(object interface{}, ids interface{}) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("GET accepts a pointer to slice of a struct type as an object")
	}

	_idsKind := KindOf(ids)
	if _idsKind != SliceOfPrimitive {
		return fmt.Errorf("GET accepts slice of a primitive type as ids e.g int64 or string")
	}

	return ds.container.GET(object, ids)
}

func (ds *DataStore) DELETE(object interface{}, ids interface{}) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct && _objectKind != SliceOfStruct {
		return fmt.Errorf("DELETE accepts a slice or pointer to slice of a struct type as an object")
	}

	_idsKind := KindOf(ids)
	if _idsKind != SliceOfPrimitive {
		return fmt.Errorf("DELETE accepts slice of a primitive type as ids e.g int64 or string")
	}

	return ds.container.DELETE(object, ids)
}

func (ds *DataStore) Connect() error {
	return ds.container.Connect()
}

func (ds *DataStore) Disconnect() error {
	return ds.container.Disconnect()
}
