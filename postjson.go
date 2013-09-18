package gohst

import (
	"allochi/inflect"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
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

// Result is an implementation of sql.Result to use with methods like sql.Execute
// Its only intention is for LastInsertedId for PostJson
type Result struct {
	id    int64
	count int64
}

// LastInsertId returns the last inserted record id
func (r Result) LastInsertId() (int64, error) {
	return r.id, nil
}

// RowsAffected returns how many records were affected by a query
func (r Result) RowsAffected() (int64, error) {
	return r.count, nil
}

// NewPostJson creates new store object only
func NewPostJson(DatabaseName, User, Password string) (store PostJsonDataStore) {
	store.DatabaseName = DatabaseName
	store.User = User
	store.Password = Password
	store.CollectionNames = make(map[string]bool)
	store.CollectionStmts = make(map[string]map[string]*sql.Stmt)
	return
}

func (ds *PostJsonDataStore) Connect() (err error) {
	ds.DB, err = sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	if err != nil {
		return
	}
	err = ds.DB.Ping()
	return
}

func (ds *PostJsonDataStore) Disconnect() (err error) {
	if ds.DB != nil {
		err = ds.DB.Close()
	}
	return
}

func (ds *PostJsonDataStore) loadCollectionNames() (err error) {

	query := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';"
	rows, err := ds.sqlQuery(query)
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

func (ds PostJsonDataStore) collectionExists(name string) (exist bool, err error) {
	if len(ds.CollectionNames) == 0 {
		if ds.CheckCollections {
			err = ds.loadCollectionNames()
		}
	}

	exist = ds.CollectionNames[name]

	return
}

func (ds PostJsonDataStore) createCollection(name string) error {

	stmt := fmt.Sprintf(`CREATE TABLE %s ("id" SERIAL PRIMARY KEY, "data" json, "created_at" timestamp(6) NULL, "updated_at" timestamp(6) NULL)`, name)

	_, err := ds.sqlExecute(stmt)
	if err != nil {
		return err
	}
	ds.CollectionNames[name] = true
	ds.prepareStatements(name)
	return nil
}

func (ds PostJsonDataStore) PUT(object interface{}) error {

	// Check what kind of object passed
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

	record := Record{}
	// ID
	record.Id = _elem.FieldByName("Id").Int()

	// Data
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	record.Data = bytes.Replace(data, sqoute, sqouteESC, -1)

	// Time Stamps
	record.CreatedAt = _elem.FieldByName("CreatedAt").Interface().(time.Time)
	record.UpdatedAt = _elem.FieldByName("UpdatedAt").Interface().(time.Time)

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
		if _kind == Pointer2Struct {
			sqlStmt = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES (E'%s',NOW(),NOW()) RETURNING id", tableName, record.Data)
		} else {
			sqlStmt = fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES (E'%s',NOW(),NOW())", tableName, record.Data)
		}
	} else {
		sqlStmt = fmt.Sprintf("UPDATE %s SET data=E'%s', updated_at=NOW() WHERE id = %d", tableName, record.Data, record.Id)
	}

	result, err := ds.sqlExecute(sqlStmt)
	if _kind == Pointer2Struct {
		var id int64
		id, err = result.LastInsertId()
		_elem.FieldByName("Id").SetInt(id)
	}

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

	if ds.DB == nil {
		return nil, fmt.Errorf("Data store not connected")
	}

	result := Result{}
	err := ds.DB.QueryRow(sqlStmt).Scan(&result.id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ds PostJsonDataStore) sqlQuery(sqlStmt string) (*sql.Rows, error) {

	if ds.DB == nil {
		return nil, fmt.Errorf("Data store not connected")
	}

	rows, err := ds.DB.Query(sqlStmt)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (ds *PostJsonDataStore) prepareStatements(tableName string) (err error) {

	stmts := make(map[string]*sql.Stmt)

	insertSQL := fmt.Sprintf("INSERT INTO %s (data, created_at, updated_at) VALUES (E$1,NOW(),NOW())", tableName)
	updateSQL := fmt.Sprintf("UPDATE %s SET data=E$1, updated_at=NOW() WHERE id = $2", tableName)
	selectSQL := fmt.Sprintf("select * from %s", tableName)

	stmts["INS"], err = ds.DB.Prepare(insertSQL)
	stmts["INSID"], err = ds.DB.Prepare(insertSQL + " RETURNING id")
	stmts["UPD"], err = ds.DB.Prepare(updateSQL)
	stmts["SEL"], err = ds.DB.Prepare(selectSQL)
	stmts["SELID"], err = ds.DB.Prepare(selectSQL + " where id = $1")
	stmts["SELIN"], err = ds.DB.Prepare(selectSQL + " where id in (select unnest(string_to_array($1, ',')::integer[]))")

	return

}
