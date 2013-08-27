package gohst

import (
	"allochi/inflect"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"reflect"
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

	// TODO: Need insert logic here

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds PostJsonDataStore) GET(object interface{}, ids interface{}) (response Response) {

	recordIDs := ids.([]int64)

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
	sqlStatement := `select "id","data","created_at","updated_at" from json_` + tableName + ` where id in $1`
	rows, err := db.Query(sqlStatement, recordIDs)
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
