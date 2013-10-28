package main

import (
	"fmt"
	// "strings"
)

type Entry struct {
	field    string
	operator string
	values   interface{}
}

type Requester interface {
	Where(Entry) *Requester
	// And(Entry) *Requester
	// Or(Entry) *Requester
	Bake() string
}

type RequestChain []Entry

func (rc *RequestChain) Where(e Entry) *RequestChain {
	// rc = append(rc, e)
	return rc
}

func (rc *RequestChain) Bake() string {
	// return strings.Join(rc[0], ",")
	return ""
}

func main() {

	// e1 := Entry{"country_id", Operator.IN, []int64{1, 2, 3, 4, 5}}
	// e2 := Entry{"title", "Like", "Cheif %"}

	// var rc Requester = make(RequestChain)
	var rc Requester = make(RequestChain)
	// rc.Where(e2)

	// fmt.Println(rc.Bake())
	fmt.Println(rc)

}
