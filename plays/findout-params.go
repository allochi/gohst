package main

import (
	"fmt"
	"reflect"
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

	trxName := ""
	var query interface{}

	if len(params) > 0 {

		if reflect.TypeOf(params[0]).Kind() == reflect.String {
			trxName = params[0].(string)
		} else {
			query = params[0]
		}

		if len(params) > 1 {
			if reflect.TypeOf(params[1]).Kind() == reflect.String && trxName == "" {
				trxName = params[1].(string)
			} else {
				query = params[1]
			}
		}

	}

	fmt.Printf("found %s and %v\n", trxName, query)

}
