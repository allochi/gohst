package main

import (
	. "allochi/tcolor"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"time"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
	Country   string
}

var tc = TColor

func init() {
	Contacts = PostJsonDataStore{"allochi_contactizer", "allochi", ""}
}

func main() {
	SamplePostJSon()
}

func SamplePostJSon() {

	DataStore = Contacts
	var contacts []Contact
	response := GET(&contacts, nil)
	fmt.Println(response)
	spew.Dump(contacts)

}

func SampleArray() {
	allochi := Person{"Ali", "Anwar", 40, "Switzerland"}
	vanessa := Person{"Vanessa", "-", 26, "Switzerland"}

	// people = append(people, vanessa)

	DataStore = People
	fmt.Printf("Main (people b4PUT): %s\n", People)
	response := PUT(allochi)
	fmt.Printf("%s\n", response.Message)
	if response.Message == "OK" {
		fmt.Printf("$v inserted", allochi)
	}
	if PUT(vanessa).Message == "OK" {
		fmt.Printf("$v inserted", vanessa)
	}
	fmt.Printf("Main (people A8PUT): %s\n", People)

	DataStore = others
	fmt.Printf("Main (others b4PUT): %s\n", others)
	PUT(allochi)
	fmt.Printf("Main (others A8PUT): %s\n", others)
}

// --------------------------------------------------------------------------------
// Gohst
// --------------------------------------------------------------------------------
type DataStoreDeligate interface {
	PUT(interface{}) Response
	GET(interface{}, interface{}) Response
}

var DataStore DataStoreDeligate

type Response struct {
	Message string
	Error   error
	Size    int
}

type Request struct {
	Ids    []int64
	Source string
}

func PUT(object interface{}) Response {
	log.Printf("DataStore: %v", DataStore)
	return DataStore.PUT(object)
}

func GET(object interface{}, query interface{}) Response {
	log.Printf("DataStore: %v", DataStore)
	return DataStore.GET(object, query)
}

// --------------------------------------------------------------------------------
// Array Data Stores
// --------------------------------------------------------------------------------
type PeopleDataStore []Person

var People PeopleDataStore

func (ds PeopleDataStore) PUT(object interface{}) (response Response) {
	_object := reflect.ValueOf(object)
	_people := reflect.Indirect(reflect.ValueOf(&People))
	_people.Set(reflect.Append(_people, _object))

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds PeopleDataStore) GET(interface{}, interface{}) (response Response) {
	return Response{}
}

type OthersDataStore []Person

var others OthersDataStore

func (ds OthersDataStore) PUT(object interface{}) (response Response) {
	_object := reflect.ValueOf(object)
	_others := reflect.Indirect(reflect.ValueOf(&others))
	_others.Set(reflect.Append(_others, _object))

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds OthersDataStore) GET(interface{}, interface{}) (response Response) {
	return Response{}
}

// --------------------------------------------------------------------------------
// PostJSON Data Stores
// --------------------------------------------------------------------------------
type PostJsonDataStore struct {
	DatabaseName string
	User         string
	Password     string
}

var Contacts PostJsonDataStore

// func (ds PostJsonDataStore) Init() {
// 	log.Println("Initializing Database Connection...")

// 	db, _ = sql.Open("postgres", "user="+ds.User+" dbname="+dbName+" sslmode=disable")
// 	defer ds.Database.Close()

// 	if err := ds.Database.Ping(); err != nil {
// 		log.Fatalf("Couldn't connect to the database: %s", err)
// 	}
// }

func (ds PostJsonDataStore) PUT(object interface{}) (response Response) {

	// TODO: Need insert logic here

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds PostJsonDataStore) GET(object interface{}, ids interface{}) (response Response) {

	_slice := reflect.Indirect(reflect.ValueOf(object))
	_type := reflect.TypeOf(object).Elem().Elem()

	db, _ := sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Couldn't connect to the database: %s", err)
	}

	rows, err := db.Query(`select "id","data","created_at","updated_at" from json_contacts order by "id" limit $1;`, 10)
	// rows, err := db.Query(`select "id","data","created_at","updated_at" from json_contacts order by "id";`)
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

type Record struct {
	Id        int64
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Contact struct {
	Id        int64     `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	JobTitle  int64     `json:"job_title_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
