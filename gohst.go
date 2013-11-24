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
	GetById(interface{}, []int64) error
	GetRawById(interface{}, []int64) (string, error)
	GetRaw(interface{}, Requester) (string, error)
	Delete(interface{}, Requester) error
	DeleteById(interface{}, []int64) error
	Index(interface{}, string) error
	Prepare(string, interface{}, Requester) error
	ExecutePrepared(string, interface{}, ...interface{}) error
	Execute(interface{}, string) error
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

	return ds.container.Get(object, request)
}

func (ds *DataStore) GetById(object interface{}, ids []int64) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	return ds.container.GetById(object, ids)
}

func (ds *DataStore) GetAll(object interface{}) error {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	return ds.container.GetById(object, []int64{})
}

// Works just like Get() but returns a JSON array in a string instead of objects array.
func (ds *DataStore) GetRaw(object interface{}, request Requester) (string, error) {

	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return "", fmt.Errorf("gohst.GetRaw() only accepts an object or a pointer to an object of type struct")
	}

	return ds.container.GetRaw(object, request)
}

// Works just like GetRaw() and returns a JSON array in a string based on a list of ids
func (ds *DataStore) GetRawById(object interface{}, ids []int64) (result string, err error) {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return "", fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	return ds.container.GetRawById(object, ids)
}

func (ds *DataStore) GetAllRaw(object interface{}) (string, error) {

	_objectKind := KindOf(object)
	if _objectKind != Pointer2SliceOfStruct {
		return "", fmt.Errorf("gohst.Get() accepts a pointer to slice of a struct type as an object")
	}

	return ds.container.GetRawById(object, []int64{})
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

// Delete all objects based on their IDs
func (ds *DataStore) DeleteById(object interface{}, ids []int64) error {

	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.DeleteById() accepts only an object or a pointer to an object of type struct")
	}

	return ds.container.DeleteById(object, ids)
}

// Delete all objects based on a query
func (ds *DataStore) Delete(object interface{}, request Requester) error {

	_kind := KindOf(object)
	if _kind != Struct && _kind != Pointer2Struct {
		return fmt.Errorf("gohst.Delete() accepts only an object or a pointer to an object of type struct")
	}

	return ds.container.Delete(object, request)
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
