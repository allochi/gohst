package main

import (
	. "allochi/gohst"
	"fmt"
	"reflect"
)

type Person struct {
	Name string
}

func main() {
	// test()

	person := Person{"Allochi"}
	fmt.Println(IsStructOrPtr2Struct(person))
	fmt.Println(IsStructOrPtr2Struct(&person))

	var something []uint64
	fmt.Println(IsPtr2SliceOfPrimitive(&something))

	var people []Person
	fmt.Println(IsPtr2SliceOfStruct(people))
	fmt.Println(IsPtr2SliceOfStruct(&people))

	fmt.Println(IsPtr2SliceOfStruct(Person{}))
	fmt.Println(IsPtr2SliceOfStruct(&Person{}))

	fmt.Println(IsStructOrPtr2Struct(people))
	fmt.Println(IsStructOrPtr2Struct(&people))

}

func check(i interface{}) {

	_type1 := reflect.TypeOf(i).Kind() == reflect.Ptr
	fmt.Printf("Is it a pointer? %v\n", _type1)

	_type2 := reflect.TypeOf(i).Kind() == reflect.Slice
	fmt.Printf("Is it a slice? %v\n", _type2)

	_type3 := reflect.ValueOf(i).Type().Kind()
	// _type3 := reflect.ValueOf(i).Elem().Type().Kind()
	fmt.Printf("Is it pointer to a slice? %v\n", _type3)

	_type4 := reflect.ValueOf(i).Elem().Type().Elem().Kind()
	fmt.Printf("Is it pointer to a slice of struct? %v\n", _type4)

}

func test() {

	var i int64

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "&Person{}", IsPtr2Slice(&Person{}))
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "&Person{}", IsPtr2Struct(&Person{}))

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "&[]string{}", IsPtr2Slice(&[]string{}))
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "&[]string{}", IsPtr2Struct(&[]string{}))

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "&i", IsPtr2Slice(&i))
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "&i", IsPtr2Struct(&i))

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "Person{}", IsPtr2Slice(Person{}))
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "Person{}", IsPtr2Struct(Person{}))

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "[]string{}", IsPtr2Slice([]string{}))
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "[]string{}", IsPtr2Struct([]string{}))

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "i", IsPtr2Slice(i))
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "i", IsPtr2Struct(i))

}
