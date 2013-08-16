package main

import (
	. "allochi/tcolor"
	"fmt"
	"reflect"
)

type Person struct {
	FirstName string
	LastName  string
	Age       int
	Country   string
}

var tc = TColor

func main() {

	allochi := Person{"Ali", "Anwar", 40, "Switzerland"}
	vanessa := Person{"Vanessa", "-", 26, "Switzerland"}

	// people = append(people, vanessa)

	DataStore = People
	fmt.Printf("Main (people b4PUT): %s\n", people)
	PUT(allochi)
	PUT(vanessa)
	fmt.Printf("Main (people A8PUT): %s\n", people)

	// DataStore = Others
	// fmt.Printf("Main (others b4PUT): %s\n", others)
	// PUT(allochi)
	// fmt.Printf("Main (others A8PUT): %s\n", others)

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
	Error error
	Size  int
}

type Request struct {
	Ids    []int64
	Source string
}

func PUT(object interface{}) Response {
	DataStore.PUT(object)
	return Response{}
}

func GET(object *interface{}, query interface{}) Response {
	DataStore.GET(object, query)
	return Response{}
}

// --------------------------------------------------------------------------------
// Data Stores
// --------------------------------------------------------------------------------
var people []Person
var People PeopleDataStore

type PeopleDataStore DataStoreDeligate

func Put(object interface{}) {
	_object := reflect.ValueOf(object)
	_people := reflect.Indirect(reflect.ValueOf(&people))
	_people.Set(reflect.Append(_people, _object))
}

func GET(*interface{}, interface{}) Response {
	return Response{}
}

// var others []Person
// var Others DataStoreDeligate

// func OthersPut(object interface{}) {
// 	_object := reflect.ValueOf(object)
// 	_others := reflect.Indirect(reflect.ValueOf(&others))
// 	_others.Set(reflect.Append(_others, _object))
// }
