package gohst

import (
	"allochi/inflect"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	// "log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	sqoute    = []byte("'")
	sqouteESC = []byte("''")
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

type Result struct {
	id    int64
	count int64
}

func (r Result) LastInsertId() (int64, error) {
	return r.id, nil
}

func (r Result) RowsAffected() (int64, error) {
	return r.count, nil
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

	_, err := ds.sqlExecute(stmt)
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
	record.Data = bytes.Replace(data, sqoute, sqouteESC, -1)
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

	var sqlStmt string
	if record.Id == 0 {
		sqlStmt = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES (E'%s',NOW(),NOW()) RETURNING id", tableName, record.Data)
	} else {
		sqlStmt = fmt.Sprintf("UPDATE %s SET data=E'%s', updated_at=NOW() WHERE id = %d", tableName, record.Data, record.Id)
	}

	result, err := ds.sqlExecute(sqlStmt)
	id, err := result.LastInsertId()
	_elem.FieldByName("Id").SetInt(id)
	spew.Dump(object)
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

func (ds PostJsonDataStore) sqlExecute(sqlStmt string) (sql.Result, error) {

	db, err := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	result := Result{}
	err = db.QueryRow(sqlStmt).Scan(&result.id)
	// result, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ds PostJsonDataStore) sqlQuery(sqlStmt string) (*sql.Rows, error) {

	db, err := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
