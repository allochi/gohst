package gohst

import (
	"allochi/inflect"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type PostJsonDataStore struct {
	DatabaseName          string
	User                  string
	Password              string
	CheckCollections      bool
	AutoCreateCollections bool
	CollectionNames       []string
}

type Record struct {
	Id        int64
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPostJson(DatabaseName, User, Password string) (store PostJsonDataStore) {
	store.DatabaseName = DatabaseName
	store.User = User
	store.Password = Password
	return
}

func (ds PostJsonDataStore) loadCollectionNames() error {
	// TODO: Implement
	return errors.New("Not Implemented")
}

func (ds PostJsonDataStore) collectionExists(name string) bool {
	for _, _name := range ds.CollectionNames {
		if name == _name {
			return true
		}
	}
	return false
}

func (ds PostJsonDataStore) createCollection(name string) error {
	// TODO: Create the table
	var err error
	if err != nil {
		return err
	}
	ds.CollectionNames = append(ds.CollectionNames, name)
	return nil
}

func (ds PostJsonDataStore) PUT(object interface{}) error {

	record := Record{}

	_elem := reflect.ValueOf(object)

	// ID is a Sequence of integers, Not sure about other ID types
	record.Id = _elem.FieldByName("Id").Int()

	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	record.Data = data
	record.CreatedAt = _elem.FieldByName("CreatedAt").Interface().(time.Time)
	record.UpdatedAt = _elem.FieldByName("UpdatedAt").Interface().(time.Time)

	_type := reflect.TypeOf(object)
	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	if ds.CheckCollections && !ds.collectionExists(tableName) {
		if ds.AutoCreateCollections {
			ds.createCollection(tableName)
		} else {
			return fmt.Errorf("Data store \"%s\" collection doesn't exist", tableName)
		}
	}
	// Check if the table exists using table slice
	// If not in table slice then execute "IF EXIST" SQL and add it to the slice
	// If the SQL is false, check if "AutoCreateStore" is true and create the table

	// TODO: Transaction?
	var sqlStatement string
	if record.Id == 0 {
		sqlStatement = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES ('%s',NOW(),NOW())", tableName, record.Data)
	} else {
		sqlStatement = fmt.Sprintf("UPDATE %s SET data='%s', updated_at=NOW() WHERE id = %d", tableName, record.Data, record.Id)
	}

	db, err := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (ds PostJsonDataStore) GET(object interface{}, ids interface{}) error {

	_slice := reflect.Indirect(reflect.ValueOf(object))
	_type := reflect.TypeOf(object).Elem().Elem()

	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	db, err := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	var sqlStatement string
	if ids != nil {
		_ids := ids.([]int64)
		_idsStr := make([]string, len(_ids))
		for i, id := range _ids {
			_idsStr[i] = strconv.FormatInt(id, 10)
		}
		sqlStatement = fmt.Sprintf("select * from %s where id in (%s);", tableName, strings.Join(_idsStr, ","))
	} else {
		sqlStatement = fmt.Sprintf("select * from %s;", tableName)
	}
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var record Record
		_object := reflect.New(_type)

		rows.Scan(&record.Id, &record.Data, &record.CreatedAt, &record.UpdatedAt)

		_object.Elem().FieldByName("Id").SetInt(record.Id)
		_object.Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(record.CreatedAt))
		_object.Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(record.UpdatedAt))
		json.Unmarshal(record.Data, _object.Interface())
		_slice.Set(reflect.Append(_slice, _object.Elem()))
	}

	return nil
}
