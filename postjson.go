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
	DatabaseName string
	User         string
	Password     string
}

type Record struct {
	Id        int64
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
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
	tableName := inflect.Pluralize(inflect.Underscore(_typeName))

	// Transaction?
	// is it insert or update?
	var sqlStatement string
	if record.Id == 0 {
		// TODO: If the store is not there
		sqlStatement = fmt.Sprintf("INSERT INTO json_%s (data, created_at, updated_at) VALUES ('%s',NOW(),NOW())", tableName, record.Data)
	} else {
		sqlStatement = fmt.Sprintf("UPDATE json_%s SET data='%s', updated_at=NOW() WHERE id = %d", tableName, record.Data, record.Id)
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
	tableName := inflect.Pluralize(inflect.Underscore(_typeName))

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
		sqlStatement = fmt.Sprintf("select * from json_%s where id in (%s);", tableName, strings.Join(_idsStr, ","))
	} else {
		sqlStatement = fmt.Sprintf("select * from json_%s;", tableName)
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
