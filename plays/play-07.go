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
	fmt.Println("KindOf(person) == Struct: ", KindOf(person) == Struct)
	fmt.Println("KindOf(&person) == Pointer2Struct: ", KindOf(&person) == Pointer2Struct)
	fmt.Println("KindOf(person) == Pointer2Struct: ", KindOf(person) == Pointer2Struct)

	var something []uint64
	// fmt.Println(IsPtr2SliceOfPrimitive(&something))
	fmt.Println("KindOf(&something) == Pointer2SliceOfPrimitive: ", KindOf(&something) == Pointer2SliceOfPrimitive)
	// fmt.Println("KindOf(&something) == Pointer2SliceOfPrimitive: ", KindOf(&something))

	var people []Person
	// fmt.Println(IsPtr2SliceOfStruct(people))
	fmt.Println("KindOf(people) == Pointer2SliceOfStruct: ", KindOf(people) == Pointer2SliceOfStruct)
	fmt.Println("KindOf(people) == SliceOfStruct: ", KindOf(people) == SliceOfStruct)
	// fmt.Println(IsPtr2SliceOfStruct(&people))
	fmt.Println("KindOf(&people) == Pointer2SliceOfStruct: ", KindOf(&people) == Pointer2SliceOfStruct)

	// fmt.Println(IsPtr2SliceOfStruct(Person{}))
	fmt.Println("KindOf(Person{}) == Pointer2SliceOfStruct: ", KindOf(Person{}) == Pointer2SliceOfStruct)
	// fmt.Println(IsPtr2SliceOfStruct(&Person{}))
	fmt.Println("KindOf(&Person{}) == Pointer2SliceOfStruct: ", KindOf(&Person{}) == Pointer2SliceOfStruct)

	// fmt.Println(IsStructOrPtr2Struct(people))
	fmt.Println("KindOf(people) == Struct || KindOf(people) == Pointer2Struct: ", KindOf(people) == Struct || KindOf(people) == Pointer2Struct)
	// fmt.Println(IsStructOrPtr2Struct(&people))
	fmt.Println("KindOf(&people) == Struct || KindOf(&people) == Pointer2Struct: ", KindOf(&people) == Struct || KindOf(&people) == Pointer2Struct)

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

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "&Person{}", KindOf(&Person{}) == Pointer2Slice)
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "&Person{}", KindOf(&Person{}) == Pointer2Struct)

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "&[]string{}", KindOf(&[]string{}) == Pointer2Slice)
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "&[]string{}", KindOf(&[]string{}) == Pointer2Struct)

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "&i", KindOf(&i) == Pointer2Slice)
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "&i", KindOf(&i) == Pointer2Struct)

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "Person{}", KindOf(Person{}) == Pointer2Slice)
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "Person{}", KindOf(Person{}) == Pointer2Struct)

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "[]string{}", KindOf([]string{}) == Pointer2Slice)
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "[]string{}", KindOf([]string{}) == Pointer2Struct)

	fmt.Printf("Is %12s a pointer to  slice? %8t \n", "i", KindOf(i) == Pointer2Slice)
	fmt.Printf("Is %12s a pointer to struct? %8t \n", "i", KindOf(i) == Pointer2Struct)

}
