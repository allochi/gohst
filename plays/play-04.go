package main

import (
	. "allochi/tcolor"
	"database/sql"
	"fmt"
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
	Contacts = PostJsonDataStore{nil, "allochi_contactizer", "allochi", ""}
}

func main() {
	SamplePostJSon()
}

func SamplePostJSon() {

	DataStore = Contacts
	response := GET(nil, nil)
	fmt.Println(response)

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
	GET(*interface{}, interface{}) Response
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

func GET(object *interface{}, query interface{}) Response {
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

func (ds PeopleDataStore) GET(*interface{}, interface{}) (response Response) {
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

func (ds OthersDataStore) GET(*interface{}, interface{}) (response Response) {
	return Response{}
}

// --------------------------------------------------------------------------------
// PostJSON Data Stores
// --------------------------------------------------------------------------------
type PostJsonDataStore struct {
	Database     *sql.DB
	DatabaseName string
	User         string
	Password     string
}

var Contacts PostJsonDataStore

func (ds PostJsonDataStore) Init() {
	log.Println("Initializing Database Connection...")

	ds.Database, _ = sql.Open("postgres", "user="+ds.User+" dbname="+ds.DatabaseName+" sslmode=disable")
	// defer ds.Database.Close()

	if err := ds.Database.Ping(); err != nil {
		log.Fatalf("Couldn't connect to the database: %s", err)
	}
}

func (ds PostJsonDataStore) PUT(object interface{}) (response Response) {

	// TODO: Need insert logic here

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds PostJsonDataStore) GET(object *interface{}, ids interface{}) (response Response) {

	ds.Database.Query("select * from json_contacts limit $1;", 10)
	// rows, err := ds.Database.Query("select * from json_contacts limit $1;", 10)
	// if err != nil {
	// 	log.Fatalf("Couldn't select table information: %s", err)
	// }

	// fmt.Printf("rows: %v\n", rows)

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return Response{}
}

type Contact struct {
	Id        int64
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
