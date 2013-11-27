package gohst

import (
	"fmt"
	"reflect"
	"time"
)

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
	Put(interface{}, Trx) error
	Get(interface{}, Requester, Trx) error
	GetById(interface{}, []int64, Trx) error
	GetRawById(interface{}, []int64) (string, error)
	GetRaw(interface{}, Requester) (string, error)
	Delete(interface{}, Requester) error
	DeleteById(interface{}, []int64) error
	Index(interface{}, string) error
	Prepare(string, interface{}, Requester) error
	ExecutePrepared(string, interface{}, ...interface{}) error
	Execute(interface{}, string) error
	ExecuteRaw(string) (string, error)
	Drop(interface{}, bool) error
	Begin(string) (Trx, error)
	Commit(Trx) error
	Rollback(Trx) error
}

// When creating a data store, gohst use Register() to keep a reference by name of that store
// this way, the store can be retried in any part of the code, it's very important
// to use Register() and GetDataStore() and use the data store they return, as all the validation
// is done on this layer. Data store names can't be reassigned once their are used.
// More than a store can be created and registered to the same database.
func Register(name string, container DataStoreContainer) error {
	if name == "" {
		return fmt.Errorf("gohst.Register() can't have empty data store name")
	}

	if container == nil {
		return fmt.Errorf("gohst.Register() can't have empty data store")
	}

	if datastores[name] != nil {
		return fmt.Errorf("gohst.Register() can't reassign data store name")
	}

	datastores[name] = container
	return nil
}

// When creating a data store, gohst use Register() to keep a reference by name of that store
// this way, the store can be retried in any part of the code, it's very important
// to use Register() and GetDataStore() and use the data store they return, as all the validation
// is done on this layer.
func GetDataStore(name string) (ds DataStore, err error) {
	if name == "" {
		err = fmt.Errorf("gohst.GetDataStore() can't have empty data store name")
	}

	dsc := datastores[name]

	if dsc == nil {
		err = fmt.Errorf("gohst.GetDataStore() couldn't find the data store")
	} else {
		ds = DataStore{dsc}
	}

	return ds, err
}

// Retrieves an array of objects using their IDs. A pointer to a slice should be passed
// with another array of IDs, the function uses the try of the slice and fill the slice with
// retrieved objects, if the slice is not empty it will be appended. If the IDs slice is empty
// all object in the table will be retrieved. This function doesn't check for duplicates.
func (ds *DataStore) Get(object interface{}, request interface{}) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	// GetById
	if reflect.TypeOf(request).String() == "[]int64" {
		return ds.container.GetById(object, request.([]int64), Trx{})
	}

	// Get(request)
	requester, ok := request.(Requester)
	if ok {
		return ds.container.Get(object, requester, Trx{})
	}

	return fmt.Errorf("gohst.Get() has no proper request parameters to process.")
}

// Get with transaction
func (ds *DataStore) Get__(object interface{}, request interface{}, trx Trx) error {

	// return fmt.Errorf("Get__ not implemented yet!")

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	// GetById
	if reflect.TypeOf(request).String() == "[]int64" {
		return ds.container.GetById(object, request.([]int64), trx)
	}

	// Get(request)
	requester, ok := request.(Requester)
	if ok {
		return ds.container.Get(object, requester, trx)
	}

	return fmt.Errorf("gohst.Get__() has no proper request parameters to process.")
}

// Get all objects from datastore of the type of passed slice
func (ds *DataStore) GetAll(object interface{}) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.GetAll() accepts a pointer to slice of a struct type as an object")
	}

	return ds.container.GetById(object, []int64{}, Trx{})

}

// Works just like Get() but returns a JSON array in a string instead of objects array.
func (ds *DataStore) GetRaw(object interface{}, params ...interface{}) (string, error) {

	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return "", fmt.Errorf("gohst.GetRaw() only accepts an object or a pointer to an object of type struct")
	}

	// Check the type of params
	if len(params) > 0 {
		options := params[0]

		// GetRawById
		if reflect.TypeOf(options).String() == "[]int64" {
			return ds.container.GetRawById(object, options.([]int64))
		}

		// GetRaw(request)
		request, ok := options.(Requester)
		if ok {
			return ds.container.GetRaw(object, request)
		}
	} else {
		// GetRawAll
		return ds.container.GetRawById(object, []int64{})
	}

	return "", fmt.Errorf("gohst.GetRaw() has no proper request parameters to process.")
}

// Execute a procedure in the database and return an array of objects, the array is of the same type
// of the passed array. Only fields that matches names between the query and the object type will be filled.
// Execute can be used to have a custom struct type filled by custom SQL statement
func (ds *DataStore) Execute(object interface{}, procedure string) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	if procedure == "" {
		return fmt.Errorf("gohst.Execute() requires procedure name as a string")
	}

	return ds.container.Execute(object, procedure)
}

// Works just like Execute() but returns a JSON array in a string instead of objects array.
// This is more fixable than Execute() since the SQL statement in function can hold any number
// of fields from multiple tables join, make sure to name them properly in case of aggregations
// and joins
func (ds *DataStore) ExecuteRaw(procedure string) (string, error) {

	if procedure == "" {
		return "", fmt.Errorf("gohst.ExecuteRaw() requires procedure name as a string")
	}

	return ds.container.ExecuteRaw(procedure)
}

// Put() is not as simple as it looks, it does simply put an object in the data store, but, it create
// new one in the store if the Id=0, and if it's not, then it updates the object. Also, if the data store
// has CheckCollections = true it will check first if the collection exists, otherwise returns an error.
// And if AutoCreateCollections = true, then it will create one if it doesn't exist.
func (ds *DataStore) Put(object interface{}) error {

	// object should be on struct, pointer to struct, slice of struct or pointer to slice of struct
	_kind := KindOf(object)
	if _kind != SliceOfStruct && _kind != Pointer2SliceOfStruct && _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Put() accepts struct, or a slice of struct as an object or pointer to these kinds.")
	}

	return ds.container.Put(object, Trx{})
}

func (ds *DataStore) Put__(object interface{}, trx Trx) error {

	// object should be on struct, pointer to struct, slice of struct or pointer to slice of struct
	_kind := KindOf(object)
	if _kind != SliceOfStruct && _kind != Pointer2SliceOfStruct && _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Put() accepts struct, or a slice of struct as an object or pointer to these kinds.")
	}

	return ds.container.Put(object, trx)
}

// Delete objects
func (ds *DataStore) Delete(object interface{}, params ...interface{}) error {

	// object should be on struct, pointer to struct, slice of struct or pointer to slice of struct
	_kind := KindOf(object)
	if _kind != SliceOfStruct && _kind != Pointer2SliceOfStruct && _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Delete() accepts struct, or a slice of struct as an object or pointer to these kinds.")
	}

	// Check the type of params
	if len(params) > 0 {
		options := params[0]

		// DeleteById
		if reflect.TypeOf(options).String() == "[]int64" {
			return ds.container.DeleteById(object, options.([]int64))
		}

		// Delete(request)
		request, ok := options.(Requester)
		if ok {
			return ds.container.Delete(object, request)
		}
	} else {
		_value := reflect.ValueOf(object)
		var ids []int64
		switch _kind {
		case Pointer2Struct:
			_value = _value.Elem()
			fallthrough
		case Struct:
			ids = append(ids, _value.FieldByName("Id").Interface().(int64))
		case Pointer2SliceOfStruct:
			_value = _value.Elem()
			fallthrough
		case SliceOfStruct:
			for i := 0; i < _value.Len(); i++ {
				ids = append(ids, _value.Index(i).FieldByName("Id").Interface().(int64))
			}
		}
		return ds.container.DeleteById(object, ids)
	}

	return fmt.Errorf("gohst.Delete() has no proper request parameters to process.")
}

// Connect to the database
func (ds *DataStore) Connect() error {

	err := ds.container.Connect()

	if err != nil {
		return fmt.Errorf("gohst.Connect() couldn't connect: %s", err)
	}

	return nil
}

// Disconnect from the database
func (ds *DataStore) Disconnect() error {

	return ds.container.Disconnect()

}

// Index takes an object e.g Contact{} and a name of field in the object e.g. "FirstName"
// and an SQL type e.g "VARCHAR" to create an index on that field named {field_name}_idx.
// In some test the speed of indexed field search was 4.7% of of same search without index.
func (ds *DataStore) Index(object interface{}, field string) error {
	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Index() only accepts an object or a pointer to an object of type struct")
	}

	if field == "" {
		return fmt.Errorf("gohst.Index() field name can't be empty")
	}

	return ds.container.Index(object, field)
}

func (ds *DataStore) Prepare(name string, object interface{}, request Requester) error {

	if name == "" {
		return fmt.Errorf("gohst.Prepare requires a name")
	}

	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Prepare() accepts only an object or a pointer to an object of type struct")
	}

	return ds.container.Prepare(name, object, request)
}

func (ds *DataStore) ExecutePrepared(name string, object interface{}, values ...interface{}) error {

	if name == "" {
		return fmt.Errorf("gohst.ExecutePrepared requires a name")
	}

	return ds.container.ExecutePrepared(name, object, values...)
}

func (ds *DataStore) Drop(object interface{}, confirmed bool) error {
	return ds.container.Drop(object, confirmed)
}

func (ds *DataStore) Begin(name string) (Trx, error) {
	if name == "" {
		name = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return ds.container.Begin(name)
}

func (ds *DataStore) Commit(trx Trx) error {
	return ds.container.Commit(trx)
}

func (ds *DataStore) Rollback(trx Trx) error {
	return ds.container.Rollback(trx)
}
