package main

import (
	"allochi/gohst"
	"fmt"
	"reflect"
	// "time"
)

// Get
// GetRaw
// Put
// Delete
// Execute
// ExecuteRaw
// ExecutePrepared

func main() {

	// trx1 := gohst.Trx{}
	// trx2 := gohst.Trx{}
	// fmt.Printf("%t\n", &trx1 == &trx2)

	request := &gohst.RequestChain{}
	trx := gohst.Trx{}

	var r gohst.Requester
	var t gohst.Trx
	var err error

	r, t, err = modify(request)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%v %v\n", r, t)

	r, t, err = modify(trx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%v %v\n", r, t)

	r, t, err = modify(request, trx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%v %v\n", r, t)

	// One Request One Trx!
	r, t, err = modify(request, trx, trx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%v %v\n", r, t)

	r, t, err = modify(request, request, trx)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("%v %v\n", r, t)

}

func modify(params ...interface{}) (request gohst.Requester, transaction gohst.Trx, err error) {

	var r, t bool

	if len(params) > 2 {
		return nil, gohst.Trx{}, fmt.Errorf("More than two parameters has been passed!")
	}

	for _, param := range params {
		if reflect.TypeOf(param).Name() == "Trx" {
			if t {
				return nil, gohst.Trx{}, fmt.Errorf("More than a transaction has been passed!")
			}
			transaction = param.(gohst.Trx)
			t = true
		}

		if req, ok := param.(gohst.Requester); ok {
			if r {
				return nil, gohst.Trx{}, fmt.Errorf("More than a request has been passed!")
			}
			request = req
			r = true
		}
	}

	return

}
