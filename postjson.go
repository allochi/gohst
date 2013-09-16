package gohst

import (
	"allochi/inflect"
	"database/sql"
	"encoding/json"
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

func (ds *PostJsonDataStore) loadCollectionNames() (err error) {

	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';"
	rows, err := ds.sqlQuery(query)
	defer rows.Close()
	var names []string
	if err == nil {
		for rows.Next() {
			var table_name string
			rows.Scan(&table_name)
			names = append(names, table_name)
		}
		ds.CollectionNames = names
	}

	return
}

func (ds PostJsonDataStore) collectionExists(name string) (exist bool, err error) {
	if len(ds.CollectionNames) == 0 {
		if ds.CheckCollections {
			err = ds.loadCollectionNames()
		}
	}

	for _, _name := range ds.CollectionNames {
		if name == _name {
			return true, nil
		}
	}
	return
}

func (ds PostJsonDataStore) createCollection(name string) error {

	stmt := fmt.Sprintf(`CREATE TABLE %s ("id" SERIAL PRIMARY KEY, "data" json, "created_at" timestamp(6) NULL, "updated_at" timestamp(6) NULL)`, name)

	err := sqlExecute(stmt)
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

	if ds.CheckCollections {
		exist, err := ds.collectionExists(tableName)
		if err != nil {
			return fmt.Errorf("Couldn't check for data store collection \"%s\" - %s", err)
		}
		if !exist && err == nil {
			if ds.AutoCreateCollections {
				ds.createCollection(tableName)
			} else {
				return fmt.Errorf("Data store \"%s\" collection doesn't exist", tableName)
			}
		}
	}

	var sqlStmts string
	if record.Id == 0 {
		sqlStmts = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES ('%s',NOW(),NOW())", tableName, record.Data)
	} else {
		sqlStmts = fmt.Sprintf("UPDATE %s SET data='%s', updated_at=NOW() WHERE id = %d", tableName, record.Data, record.Id)
	}

	_, err = ds.sqlExecute(sqlStmts)

	return err
}

func (ds PostJsonDataStore) GET(object interface{}, ids interface{}) error {

	_slice := reflect.Indirect(reflect.ValueOf(object))
	_type := reflect.TypeOf(object).Elem().Elem()

	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	var sqlStmt string
	if ids != nil {
		_ids := ids.([]int64)
		_idsStr := make([]string, len(_ids))
		for i, id := range _ids {
			_idsStr[i] = strconv.FormatInt(id, 10)
		}
		sqlStmt = fmt.Sprintf("select * from %s where id in (%s);", tableName, strings.Join(_idsStr, ","))
	} else {
		sqlStmt = fmt.Sprintf("select * from %s;", tableName)
	}
	rows, err := ds.sqlQuery(sqlStmt)
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

func (ds PostJsonDataStore) sqlExecute(sqlStmts string) (result sql.Result, err error) {

	db, err := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return
	}
	defer db.Close()

	result, err = db.Exec(sqlStmts)
	if err != nil {
		return
	}

	return
}

func (ds PostJsonDataStore) sqlQuery(sqlStmt string) (rows *sql.Rows, err error) {

	db, err := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return
	}
	defer db.Close()

	rows, err = db.Query(sqlStmt)
	if err != nil {
		return
	}

	return
}
