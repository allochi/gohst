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
	fmt.Printf("Main (people b4PUT): %s\n", People)
	response := PUT(allochi)
	fmt.Println(response)
	if PUT(allochi).Type == "OK" {
		fmt.Printf("$v inserted", allochi)
	}
	if PUT(vanessa).Type == "OK" {
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
	Type  string
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
// Array Data Stores
// --------------------------------------------------------------------------------
type PeopleDataStore []Person

var People PeopleDataStore

func (ds PeopleDataStore) PUT(object interface{}) Response {
	_object := reflect.ValueOf(object)
	_people := reflect.Indirect(reflect.ValueOf(&People))
	_people.Set(reflect.Append(_people, _object))

	response := Response{}
	response.Type = "asjdhalhdlas"
	response.Error = nil
	response.Size = 100
	fmt.Println(response)
	return response
}

func (ds PeopleDataStore) GET(*interface{}, interface{}) Response {
	return Response{}
}

type OthersDataStore []Person

var others OthersDataStore

func (ds OthersDataStore) PUT(object interface{}) Response {
	_object := reflect.ValueOf(object)
	_others := reflect.Indirect(reflect.ValueOf(&others))
	_others.Set(reflect.Append(_others, _object))
	return Response{}
}

func (ds OthersDataStore) GET(*interface{}, interface{}) Response {
	return Response{}
}
