package main

import (
	"fmt"
	// "reflect"
	"time"
)

// Get
// GetRaw
// Put
// Delete
// Execute
// ExecuteRaw
// ExecutePrepared

func main() {

	modify("Hello")
	modify("Hello", "World")
	modify("Hello", time.Now())

}

func modify(params ...interface{}) {
	if len(params) > 0 {

		fmt.Println(params[0])

		if len(params) > 1 {
			fmt.Println(params[1])
		}

	}
}
