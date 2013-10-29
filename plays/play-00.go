package main

import (
	"fmt"
	"reflect"
)

func main() {

	i := []int{4}
	fmt.Println(reflect.ValueOf(i).Type().In(i))

}
