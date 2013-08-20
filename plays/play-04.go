package main

import (
	. "allochi/tcolor"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"reflect"
)

var tc = TColor

func init() {
	Contacts = gohst.PostJsonDataStore{"allochi_contactizer", "allochi", ""}
}

func main() {
	SamplePostJSon()
}

var Contacts gohst.PostJsonDataStore

func SamplePostJSon() {

	gohst.DataStore = Contacts
	var contacts []Contact
	gohst.GET(&contacts, nil)
	for _, contact := range contacts {
		fmt.Println(contact.Name())
	}

}

func SampleArray() {
	allochi := Person{"Ali", "Anwar", 40, "Switzerland"}
	vanessa := Person{"Vanessa", "-", 26, "Switzerland"}

	// people = append(people, vanessa)

	gohst.DataStore = People
	fmt.Printf("Main (people b4PUT): %s\n", People)
	response := gohst.PUT(allochi)
	fmt.Printf("%s\n", response.Message)
	if response.Message == "OK" {
		fmt.Printf("$v inserted", allochi)
	}
	if gohst.PUT(vanessa).Message == "OK" {
		fmt.Printf("$v inserted", vanessa)
	}
	fmt.Printf("Main (people A8PUT): %s\n", People)

	gohst.DataStore = others
	fmt.Printf("Main (others b4PUT): %s\n", others)
	gohst.PUT(allochi)
	fmt.Printf("Main (others A8PUT): %s\n", others)
}

// --------------------------------------------------------------------------------
// Array Data Stores
// --------------------------------------------------------------------------------
type PeopleDataStore []Person

var People PeopleDataStore

func (ds PeopleDataStore) PUT(object interface{}) (response gohst.Response) {
	_object := reflect.ValueOf(object)
	_people := reflect.Indirect(reflect.ValueOf(&People))
	_people.Set(reflect.Append(_people, _object))

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds PeopleDataStore) GET(interface{}, interface{}) (response gohst.Response) {
	return
}

type OthersDataStore []Person

var others OthersDataStore

func (ds OthersDataStore) PUT(object interface{}) (response gohst.Response) {
	_object := reflect.ValueOf(object)
	_others := reflect.Indirect(reflect.ValueOf(&others))
	_others.Set(reflect.Append(_others, _object))

	response.Message = "Ok"
	response.Error = nil
	response.Size = 0
	return
}

func (ds OthersDataStore) GET(interface{}, interface{}) (response gohst.Response) {
	return
}
