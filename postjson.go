package gohst

import (
	"bitbucket.org/pkg/inflect"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"reflect"
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
	DB                    *sql.DB
	CheckCollections      bool
	AutoCreateCollections bool
	CollectionNames       map[string]bool
	CollectionStmts       map[string]map[string]*sql.Stmt
}

type Record struct {
	Id        int64
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPostJson creates new store object only
func NewPostJson(DatabaseName, User, Password string) *PostJsonDataStore {
	store := new(PostJsonDataStore)
	store.DatabaseName = DatabaseName
	store.User = User
	store.Password = Password
	store.CollectionNames = make(map[string]bool)
	store.CollectionStmts = make(map[string]map[string]*sql.Stmt)
	return store
}

func (ds *PostJsonDataStore) Connect() (err error) {
	ds.DB, err = sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return
	}

	err = ds.DB.Ping()
	if err != nil {
		return
	}

	err = ds.loadFunctions()
	if err != nil {
		return
	}

	err = ds.loadCollectionNames()
	return
}

func (ds *PostJsonDataStore) Disconnect() (err error) {
	if ds.DB != nil {
		err = ds.DB.Close()
	}
	return
}

func (ds *PostJsonDataStore) loadCollectionNames() (err error) {

	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name LIKE 'json_%';"
	rows, err := ds.DB.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	if err == nil {
		for rows.Next() {
			var tableName string
			rows.Scan(&tableName)
			ds.CollectionNames[tableName] = true
			ds.prepareStatements(tableName)
		}
	}

	return
}

func (ds *PostJsonDataStore) loadFunctions() (err error) {

	var count int
	query := "SELECT count(*) FROM pg_proc WHERE proname IN ('_array','_date');"
	err = ds.DB.QueryRow(query).Scan(&count)

	if err != nil {
		return
	}

	functions := []string{
		`
			CREATE OR REPLACE FUNCTION _array(_j json, _key text) RETURNS text[] as $$
			SELECT concat('{',btrim(_j->>_key,'[]'),'}')::text[]
			$$ LANGUAGE SQL IMMUTABLE;
			`,
		`
			CREATE OR REPLACE FUNCTION _date(_j json, _key text) RETURNS TIMESTAMP as $$
			SELECT (_j->>_key)::TIMESTAMP
			$$ LANGUAGE SQL IMMUTABLE;
			`,
		`
			CREATE OR REPLACE FUNCTION _in(_list text) RETURNS setof int as $$
			SELECT unnest(string_to_array(btrim(_list,'()'), ',')::integer[]);
			$$ LANGUAGE SQL IMMUTABLE;
			`,
	}

	if count < len(functions) {
		for _, sql := range functions {
			_, err = ds.DB.Exec(sql)
			if err != nil {
				return
			}
		}
	}

	return
}

func (ds *PostJsonDataStore) collectionExists(name string) (exist bool, err error) {
	if len(ds.CollectionNames) == 0 {
		if ds.CheckCollections {
			err = ds.loadCollectionNames()
		}
	}

	exist = ds.CollectionNames[name]

	return
}

func (ds *PostJsonDataStore) createCollection(name string) error {

	stmt := fmt.Sprintf(`CREATE TABLE %s ("id" SERIAL PRIMARY KEY, "data" json, "created_at" timestamp(6) NULL, "updated_at" timestamp(6) NULL)`, name)

	_, err := ds.DB.Exec(stmt)
	if err != nil {
		return err
	}
	ds.CollectionNames[name] = true
	ds.prepareStatements(name)
	return nil
}

func pack(_value reflect.Value) (record Record, err error) {

	record.Id = _value.FieldByName("Id").Int()

	data, err := json.Marshal(_value.Interface())
	if err != nil {
		return
	}

	record.Data = bytes.Replace(data, sqoute, sqouteESC, -1)
	record.CreatedAt = _value.FieldByName("CreatedAt").Interface().(time.Time)
	record.UpdatedAt = _value.FieldByName("UpdatedAt").Interface().(time.Time)

	return

}

// Put insert or update an object to the database
// If the object's Id == 0 then INSERT is used to insert a new object
// If a pointer to object is passed and Id == 0 then the Id will be updated from the database
// If Id != 0 then UPDATE is used and the object is updated
func (ds *PostJsonDataStore) Put(object interface{}) error {

	_name, _, _value := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	// Check & Create Collections
	// If you are sure that the table exist then CheckCollections can be set to false for performance

	if ds.CheckCollections {
		exist, err := ds.collectionExists(tableName)
		if err != nil {
			return fmt.Errorf("Couldn't check for data store collection \"%s\" exists - %s", err)
		}
		if !exist && err == nil {
			if ds.AutoCreateCollections {
				ds.createCollection(tableName)
			} else {
				return fmt.Errorf("Data store \"%s\" collection doesn't exist", tableName)
			}
		}
	}

	_kind := KindOf(object)

	// Single object
	switch _kind {
	case Pointer2Struct:
		_value = reflect.ValueOf(object).Elem()
		fallthrough
	case Struct:
		err := ds.saveOrUpdate(_value, _kind, tableName)
		return err
	}

	// Slice of objects
	for i := 0; i < _value.Len(); i++ {
		err := ds.saveOrUpdate(_value.Index(i), _kind, tableName)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (ds *PostJsonDataStore) saveOrUpdate(_value reflect.Value, _kind Kind, tableName string) error {
	record, err := pack(_value)
	if err != nil {
		return err
	}

	if record.Id == 0 {
		if _kind == Pointer2SliceOfStruct || _kind == Pointer2Struct {
			err := ds.CollectionStmts[tableName]["INSID"].QueryRow(record.Data).Scan(&record.Id)
			if err != nil {
				return err
			}
			_value.FieldByName("Id").SetInt(record.Id)
		} else {
			_, err := ds.CollectionStmts[tableName]["INS"].Exec(record.Data)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := ds.CollectionStmts[tableName]["UPD"].Exec(record.Data, record.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns an an array of objects based on list of ids
// This should be faster than normal Get, since it search by IDs and uses prepared statement
func (ds *PostJsonDataStore) GetById(object interface{}, ids []int64) (err error) {

	_name, _type, _value := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	var rows *sql.Rows
	switch {
	case len(ids) == 0:
		rows, err = ds.CollectionStmts[tableName]["SELALL"].Query()
	case len(ids) == 1:
		rows, err = ds.CollectionStmts[tableName]["SELID"].Query(ids[0])
		// rows, err = ds.DB.Query(fmt.Sprintf("SELECT * FROM %s where id = $1;", tableName), ids[0])
	case len(ids) > 1:
		rows, err = ds.CollectionStmts[tableName]["SELIN"].Query(IN(ids))
	}

	if err != nil {
		return
	}

	defer rows.Close()

	unpackRows(rows, _type, _value)

	return

}

// Returns an an array of objects based on list of ids
// This should be faster than normal Get, since it search by IDs and uses prepared statement
func (ds *PostJsonDataStore) GetRawById(object interface{}, ids []int64) (result string, err error) {

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	switch {
	case len(ids) == 0:
		err = ds.CollectionStmts[tableName]["SELALLRAW"].QueryRow().Scan(&result)
	case len(ids) == 1:
		err = ds.CollectionStmts[tableName]["SELIDRAW"].QueryRow(ids[0]).Scan(&result)
	case len(ids) > 1:
		err = ds.CollectionStmts[tableName]["SELINRAW"].QueryRow(IN(ids)).Scan(&result)
	}

	return

}

// Returns an an array of objects based on a request
func (ds *PostJsonDataStore) Get(object interface{}, request Requester) (err error) {

	_name, _type, _value := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	var rows *sql.Rows
	sql := fmt.Sprintf("SELECT * FROM %s %s", tableName, request.Bake(object))
	rows, err = ds.DB.Query(sql)

	if err != nil {
		return
	}

	defer rows.Close()

	unpackRows(rows, _type, _value)

	return

}

// Returns the result as a json in a string instead of unpacking it into an array of objects
func (ds *PostJsonDataStore) GetRaw(object interface{}, request Requester) (result string, err error) {

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	sql := fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s %s) row_data;", tableName, request.Bake(object))
	err = ds.DB.QueryRow(sql).Scan(&result)

	return

}

// Delete the objects from the table based on the request
func (ds *PostJsonDataStore) Delete(object interface{}, request Requester) (err error) {

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	sql := fmt.Sprintf("DELETE FROM %s %s", tableName, request.Bake(object))
	_, err = ds.DB.Exec(sql)

	return

}

// Delete the objects from the table based on the request
// This should be faster than normal Delete, since it search by IDs and uses prepared statement
func (ds *PostJsonDataStore) DeleteById(object interface{}, ids []int64) (err error) {

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	switch {
	case len(ids) == 0:
		err = fmt.Errorf("Delete can't have empty id list")
	case len(ids) == 1:
		_, err = ds.CollectionStmts[tableName]["DELID"].Exec(ids[0])
	case len(ids) > 1:
		_, err = ds.CollectionStmts[tableName]["DELIN"].Exec(IN(ids))
	}

	return

}

func (ds *PostJsonDataStore) prepareStatements(tableName string) (err error) {

	sqls := make(map[string]string)
	prepared := make(map[string]*sql.Stmt)

	// INSERT
	sqls["INS"] = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES ($1,NOW(),NOW());", tableName)
	sqls["INSID"] = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES ($1,NOW(),NOW()) RETURNING id;", tableName)

	// UPDATE
	sqls["UPD"] = fmt.Sprintf("UPDATE %s SET data=$1, updated_at=NOW() WHERE id = $2;", tableName)

	// SELECT for GetById
	sqls["SELALL"] = fmt.Sprintf("SELECT * FROM %s;", tableName)
	sqls["SELID"] = fmt.Sprintf("SELECT * FROM %s where id = $1;", tableName)
	sqls["SELIN"] = fmt.Sprintf("SELECT * FROM %s where id IN (SELECT _in($1));", tableName)

	// SELECT RAW JSON for GetRawById
	sqls["SELRAW"] = fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s) row_data;", tableName)
	sqls["SELIDRAW"] = fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s WHERE id = $1) row_data;", tableName)
	sqls["SELINRAW"] = fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s WHERE id IN (SELECT _in($1))) row_data;", tableName)

	// DELETE for DeleteById
	sqls["DELID"] = fmt.Sprintf("DELETE FROM %s where id = $1;", tableName)
	sqls["DELIN"] = fmt.Sprintf("DELETE FROM %s where id in (select _in($1));", tableName)

	for key, value := range sqls {
		prepared[key], err = ds.DB.Prepare(value)
		if err != nil {
			return
		}
	}

	ds.CollectionStmts[tableName] = prepared

	return

}

// Using index helped in some cases to reduce query time to 4.7%
func (ds *PostJsonDataStore) Index(object interface{}, field string) error {

	if field == "Id" || field == "CreatedAt" || field == "UpdatedAt" {
		return fmt.Errorf("gohst.Index: Can't index any of the fields Id, CreatedAt and UpdatedAt")
	}

	_type := TypeOf(object, field)
	if _type == "" {
		// Can't index struct ot pointers, field has to have basic type
		return fmt.Errorf("gohst.Index: Unable to index %s either it doesn't exist or it's type is not basic", field)
	}

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	_json := JsonName(object, field)
	if _json != "" {
		field = _json
	}

	sql := ""
	if strings.Contains(_type, "[]") {
		sql = fmt.Sprintf("CREATE INDEX _array_%s_%s_idx ON %s USING GIN (%s(data,'%s'));", tableName, field, tableName, "_array", field)
	} else {
		switch _type {
		case "string":
			sql = fmt.Sprintf("CREATE INDEX %s_%s_idx ON %s(((data->>'%s')::%s));", tableName, field, tableName, field, SQLTypes[_type])
		case "time.Time":
			sql = fmt.Sprintf("CREATE INDEX %s_%s_idx ON %s(_date(data,'%s'));", tableName, field, tableName, field)
		}
	}
	// simple field (int, float, string, date)
	// array (string, float, array, date)

	// fmt.Println(sql)
	ds.DB.Exec(sql)
	return nil
}

func (ds *PostJsonDataStore) Execute(object interface{}, procedure string) (err error) {

	_, _type, _value := TypeName(object)

	// TODO: SQL injection!

	sql := fmt.Sprintf("SELECT * FROM %s;", procedure)
	rows, err := ds.DB.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	unpackRows(rows, _type, _value)

	return
}

func (ds *PostJsonDataStore) ExecuteRaw(procedure string) (result string, err error) {
	// TODO: SQL injection!
	sql := fmt.Sprintf("select array_to_json(array_agg(row_to_json(row_data))) from (SELECT * FROM %s) row_data;", procedure)
	err = ds.DB.QueryRow(sql).Scan(&result)
	return
}

func (ds *PostJsonDataStore) Prepare(name string, object interface{}, request Requester) error {

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	if ds.CollectionStmts[tableName][name] != nil {
		return fmt.Errorf("gohst.Prepare [%s][%s] statement already exist", tableName, name)
	}

	sql := fmt.Sprintf("SELECT * FROM %s %s", tableName, request.Bake(object))
	prepared, err := ds.DB.Prepare(sql)
	if err != nil {
		fmt.Printf("gohst.PostJson.Prepare: %s\nQuery: %s", err, sql)
	}
	ds.CollectionStmts[tableName][name] = prepared

	return nil

}

func (ds *PostJsonDataStore) ExecutePrepared(name string, object interface{}, values ...interface{}) (err error) {

	_name, _type, _value := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	prepared := ds.CollectionStmts[tableName][name]
	if prepared == nil {
		return fmt.Errorf("gohst.ExecutePrepared [%s][%s] statement doesn't exist", tableName, name)
	}

	var rows *sql.Rows
	rows, err = prepared.Query(values...)

	if err != nil {
		return
	}
	defer rows.Close()

	unpackRows(rows, _type, _value)

	return
}

func unpackRows(rows *sql.Rows, _type reflect.Type, _value reflect.Value) {

	for rows.Next() {
		var record Record
		_object := reflect.New(_type)

		rows.Scan(&record.Id, &record.Data, &record.CreatedAt, &record.UpdatedAt)

		_object.Elem().FieldByName("Id").SetInt(record.Id)
		_object.Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(record.CreatedAt))
		_object.Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(record.UpdatedAt))
		json.Unmarshal(record.Data, _object.Interface())
		_value.Set(reflect.Append(_value, _object.Elem()))
	}

}

func (ds *PostJsonDataStore) Drop(object interface{}, confirmed bool) (err error) {

	_name, _, _ := TypeName(object)
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_name))

	sql := fmt.Sprintf("DROP TABLE %s", tableName)
	if !confirmed {
		sql = fmt.Sprintf("ALTER TABLE %s RENAME TO %s_%d;", tableName, tableName, time.Now().UnixNano())
	}

	_, err = ds.DB.Exec(sql)
	if err == nil {
		// Remove from collections
		ds.CollectionNames[tableName] = false
		ds.CollectionStmts[tableName] = nil
	}
	return
}
