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
	Put(interface{}) error
	Get(interface{}, Requester) error
	Delete(interface{}, interface{}) error
	Index(interface{}, string, string) error
	Execute(interface{}, string) error
	GetRaw(interface{}, interface{}, string) (string, error)
	ExecuteRaw(string) (string, error)
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

// Put() is not as simple as it looks, it does simply put an object in the data store, but, it create
// new one in the store if the Id=0, and if it's not, then it updates the object. Also, if the data store
// has CheckCollections = true it will check first if the collection exists, otherwise returns an error.
// And if AutoCreateCollections = true, then it will create one if it doesn't exist.
func (ds *DataStore) Put(object interface{}) error {
	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Put() only accepts an object or a pointer to an object of type struct")
	}
	return ds.container.Put(object)
}

// Retrieves an array of objects using their IDs. A pointer to a slice should be passed
// with another array of IDs, the function uses the try of the slice and fill the slice with
// retrieved objects, if the slice is not empty it will be appended. If the IDs slice is empty
// all object in the table will be retrieved. This function doesn't check for duplicates.
func (ds *DataStore) Get(object interface{}, request Requester) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	// _idsKind := KindOf(ids)
	// if _idsKind != SliceOfPrimitive {
	// 	return fmt.Errorf("gohst.Get() accepts slice of a primitive type as ids e.g int64 or string")
	// }

	return ds.container.Get(object, request)
}

// Works just like Get() but returns a JSON array in a string instead of objects array.
func (ds *DataStore) GetRaw(object interface{}, ids interface{}, sort string) (string, error) {

	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return "", fmt.Errorf("gohst.GetRaw() only accepts an object or a pointer to an object of type struct")
	}

	_idsKind := KindOf(ids)
	if _idsKind != SliceOfPrimitive {
		return "", fmt.Errorf("gohst.GetRaw() accepts slice of a primitive type as ids e.g int64 or string")
	}

	return ds.container.GetRaw(object, ids, sort)
}

// Execute a procedure in the database and return an array of objects, the array is of the same type
// of the passed array. Only fields that matches names between the query and the object type will be filled.
// Execute can be used to have a mutant struct type filled by custom SQL statement
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

// Deletes all the objects using their IDs, an array of IDs must be passed
// and the array should not be empty otherwise nothing happens
func (ds *DataStore) Delete(object interface{}, ids interface{}) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct && _objectKind != SliceOfStruct {
		return fmt.Errorf("gohst.Delete() accepts a slice or pointer to slice of a struct type as an object")
	}

	_idsKind := KindOf(ids)
	if _idsKind != SliceOfPrimitive {
		return fmt.Errorf("gohst.Delete() accepts slice of a primitive type as ids e.g int64 or string")
	}

	return ds.container.Delete(object, ids)
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
func (ds *DataStore) Index(object interface{}, field string, indexSqlType string) error {
	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Index() only accepts an object or a pointer to an object of type struct")
	}

	if field == "" {
		return fmt.Errorf("gohst.Index() field name can't be empty")
	}

	return ds.container.Index(object, field, indexSqlType)
}
