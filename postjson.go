package gohst

import (
	"allochi/inflect"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
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

func (ds *PostJsonDataStore) Put(object interface{}) error {

	// Check type of object & get the name of collection ------------------------------
	var _elem reflect.Value
	_kind := KindOf(object)
	switch _kind {
	case Struct:
		_elem = reflect.ValueOf(object)
	case Pointer2Struct:
		_elem = reflect.Indirect(reflect.ValueOf(object))
	}

	record := Record{}
	_type := _elem.Type()
	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	// ID -----------------------------------------------------------------------------
	record.Id = _elem.FieldByName("Id").Int()

	// Data ---------------------------------------------------------------------------
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	record.Data = bytes.Replace(data, sqoute, sqouteESC, -1)

	// Time Stamps --------------------------------------------------------------------

	record.CreatedAt = _elem.FieldByName("CreatedAt").Interface().(time.Time)
	record.UpdatedAt = _elem.FieldByName("UpdatedAt").Interface().(time.Time)

	// Check & Create Collections -----------------------------------------------------

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

	// Write to database --------------------------------------------------------------

	if record.Id == 0 {
		if _kind == Pointer2Struct {
			ds.CollectionStmts[tableName]["INSID"].QueryRow(record.Data).Scan(&record.Id)
			_elem.FieldByName("Id").SetInt(record.Id)
		} else {
			ds.CollectionStmts[tableName]["INS"].Exec(record.Data)
		}
	} else {
		ds.CollectionStmts[tableName]["UPD"].Exec(record.Data, record.Id)
	}

	return err
}

func (ds *PostJsonDataStore) Get(object interface{}, ids interface{}, sort string) (err error) {

	_slice := reflect.Indirect(reflect.ValueOf(object))
	_type := reflect.TypeOf(object).Elem().Elem()

	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	var rows *sql.Rows
	_ids := ids.([]int64)
	if len(_ids) > 1 {
		_idsStr := make([]string, len(_ids))
		for i, id := range _ids {
			_idsStr[i] = strconv.FormatInt(id, 10)
		}
		if sort == "" {
			rows, err = ds.CollectionStmts[tableName]["SELIN"].Query(strings.Join(_idsStr, ","))
		} else {
			sql := fmt.Sprintf("SELECT * FROM %s where id IN (SELECT unnest(string_to_array($1, ',')::integer[])) ORDER BY %s;", tableName, sort)
			rows, err = ds.DB.Query(sql, strings.Join(_idsStr, ","))
		}
	} else if len(_ids) == 1 {
		rows, err = ds.CollectionStmts[tableName]["SELID"].Query(_ids[0])
	} else {
		if sort == "" {
			rows, err = ds.CollectionStmts[tableName]["SEL"].Query()
		} else {
			sql := fmt.Sprintf("SELECT * FROM %s ORDER BY %s;", tableName, sort)
			rows, err = ds.DB.Query(sql)
		}
	}

	if err != nil {
		return
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

	return
}

func (ds *PostJsonDataStore) GetRaw(object interface{}, ids interface{}, sort string) (result string, err error) {

	var _elem reflect.Value
	_kind := KindOf(object)
	switch _kind {
	case Struct:
		_elem = reflect.ValueOf(object)
	case Pointer2Struct:
		_elem = reflect.Indirect(reflect.ValueOf(object))
	}

	_type := _elem.Type()
	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	_ids := ids.([]int64)
	if len(_ids) > 1 {
		_idsStr := make([]string, len(_ids))
		for i, id := range _ids {
			_idsStr[i] = strconv.FormatInt(id, 10)
		}
		if sort == "" {
			err = ds.CollectionStmts[tableName]["SELINRAW"].QueryRow(strings.Join(_idsStr, ",")).Scan(&result)
		} else {
			sql := fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s WHERE id IN (SELECT unnest(string_to_array($1, ',')::integer[])) ORDER BY %s) row_data;", tableName, sort)
			err = ds.DB.QueryRow(sql, strings.Join(_idsStr, ",")).Scan(&result)
		}
	} else if len(_ids) == 1 {
		err = ds.CollectionStmts[tableName]["SELIDRAW"].QueryRow(_ids[0]).Scan(&result)
	} else {
		if sort == "" {
			err = ds.CollectionStmts[tableName]["SELRAW"].QueryRow().Scan(&result)
		} else {
			sql := fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s ORDER BY %s) row_data;", tableName, sort)
			err = ds.DB.QueryRow(sql).Scan(&result)
		}
	}

	return
}

func (ds *PostJsonDataStore) Delete(object interface{}, ids interface{}) (err error) {

	_objectKind := KindOf(object)
	_type := reflect.TypeOf(object).Elem()
	if _objectKind == Pointer2SliceOfStruct {
		ds.Get(object, ids, "")
		_type = reflect.TypeOf(object).Elem().Elem()
	}

	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	_ids := ids.([]int64)
	if len(_ids) > 1 {
		_idsStr := make([]string, len(_ids))
		for i, id := range _ids {
			_idsStr[i] = strconv.FormatInt(id, 10)
		}
		_, err = ds.CollectionStmts[tableName]["DELIN"].Exec(strings.Join(_idsStr, ","))
	} else if len(_ids) == 1 {
		_, err = ds.CollectionStmts[tableName]["DELID"].Exec(_ids[0])
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

	// SELECT
	sqls["SEL"] = fmt.Sprintf("SELECT * FROM %s;", tableName)
	sqls["SELID"] = fmt.Sprintf("SELECT * FROM %s where id = $1;", tableName)
	sqls["SELIN"] = fmt.Sprintf("SELECT * FROM %s where id IN (SELECT unnest(string_to_array($1, ',')::integer[]));", tableName)

	// SELECT RAW JSON
	sqls["SELRAW"] = fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s) row_data;", tableName)
	sqls["SELIDRAW"] = fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s WHERE id = $1) row_data;", tableName)
	sqls["SELINRAW"] = fmt.Sprintf("SELECT array_to_json(array_agg(row_to_json(row_data))) FROM (SELECT * FROM %s WHERE id IN (SELECT unnest(string_to_array($1, ',')::integer[])) ) row_data;", tableName)

	// DELETE
	sqls["DELID"] = fmt.Sprintf("DELETE FROM %s where id = $1;", tableName)
	sqls["DELIN"] = fmt.Sprintf("DELETE FROM %s where id in (select unnest(string_to_array($1, ',')::integer[]));", tableName)

	for key, value := range sqls {
		prepared[key], err = ds.DB.Prepare(value)
		if err != nil {
			return
		}
	}

	ds.CollectionStmts[tableName] = prepared

	return

}

// Using index helped to reduce query time to 4.7%
func (ds *PostJsonDataStore) Index(object interface{}, field string, indexSqlType string) error {

	// Check type of object & get the name of collection ------------------------------
	var _elem reflect.Value
	_kind := KindOf(object)
	switch _kind {
	case Struct:
		_elem = reflect.ValueOf(object)
	case Pointer2Struct:
		_elem = reflect.Indirect(reflect.ValueOf(object))
	}

	_type := _elem.Type()
	_typeName := _type.Name()
	tableName := "json_" + inflect.Pluralize(inflect.Underscore(_typeName))

	sql := fmt.Sprintf("CREATE INDEX %s_idx ON %s(((data->>'%s')::%s));", field, tableName, field, indexSqlType)
	ds.DB.Exec(sql)
	return nil
}

func (ds *PostJsonDataStore) Execute(object interface{}, procedure string) (err error) {

	_slice := reflect.Indirect(reflect.ValueOf(object))
	_type := reflect.TypeOf(object).Elem().Elem()

	// TODO: SQL injection!

	sql := fmt.Sprintf("SELECT * FROM %s;", procedure)
	rows, err := ds.DB.Query(sql)
	if err != nil {
		return
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

	return
}

func (ds *PostJsonDataStore) ExecuteRaw(procedure string) (result string, err error) {
	// TODO: SQL injection!
	sql := fmt.Sprintf("select array_to_json(array_agg(row_to_json(row_data))) from (SELECT * FROM %s) row_data;", procedure)
	err = ds.DB.QueryRow(sql).Scan(&result)
	return
}
