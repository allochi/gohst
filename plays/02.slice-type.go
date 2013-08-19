package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"reflect"
)

func main() {

	var persons []Person
	MakeSlice(&persons)
	spew.Dump(persons)

}

func MakeSlice(i interface{}) {

	_slice := reflect.Indirect(reflect.ValueOf(i))
	spew.Dump(_slice)
	_type := reflect.TypeOf(i).Elem().Elem()
	spew.Dump(_type)
	_object := reflect.New(_type)
	fmt.Printf("_object: %v \n", _object)
	// _slice.Set(reflect.Append(_slice, reflect.ValueOf(Person{"Allochi"})))
	// _slice.Set(reflect.Append(_slice, _object))

}

type Person struct {
	Name string
}

type Student struct {
}

type Employee struct {
}
