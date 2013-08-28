package gohst

import (
	"allochi/inflect"
	"database/sql"
	"encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"log"
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

func (ds PostJsonDataStore) PUT(object interface{}) (response Response) {

	record := Record{}
	// isNew := true

	_elem := reflect.ValueOf(object)
	record.Id = _elem.FieldByName("Id").Int()
	data, err := json.Marshal(object)
	if err != nil {
		log.Printf("JSON Error: %s\n", err)
		return
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
		sqlStatement = fmt.Sprintf("INSERT INTO json_%s (data, created_at, updated_at) VALUES ('%s',NOW(),NOW())", tableName, record.Data)
	} else {
		sqlStatement = fmt.Sprintf("UPDATE json_%s SET data='%s', updated_at=NOW() WHERE id = %d", tableName, record.Data, record.Id)
	}

	log.Printf("%s \n", sqlStatement)

	db, _ := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	defer db.Close()

	_, err = db.Exec(sqlStatement)

	if err != nil {
		log.Fatalf("PUT - Database error: %s\n", err)
	}

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds PostJsonDataStore) GET(object interface{}, ids interface{}) (response Response) {

	_slice := reflect.Indirect(reflect.ValueOf(object))
	_type := reflect.TypeOf(object).Elem().Elem()

	_typeName := _type.Name()
	tableName := inflect.Pluralize(inflect.Underscore(_typeName))

	db, _ := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Couldn't connect to the database: %s", err)
	}

	// TODO: Generalize Query
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
	defer rows.Close()

	if err != nil {
		log.Fatalf("Couldn't select table information: %s", err)
	}

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

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return Response{}
}
