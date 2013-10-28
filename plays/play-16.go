package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	// "strings"
)

type Entry struct {
	field    string
	operator string
	values   interface{}
}

func (e *Entry) Bake() string {
	field := e.field
	switch field {
	case "Id":
		fallthrough
	case "CreatedAt":
		fallthrough
	case "UpdatedAt":
		field = fmt.Sprintf("'%s'", field)
	default:
		field = fmt.Sprintf("data->>'%s'", field)
	}
	return "..."
}

type Requester interface {
	Where(Entry) *Requester
	And(Entry) *Requester
	// Or(Entry) *Requester
	Bake() string
}

type RequestChain struct {
	entries    []Entry
	operations []string
}

func (rc *RequestChain) Where(e Entry) *RequestChain {
	rc.entries = append(rc.entries, e)
	rc.operations = append(rc.operations, "WHERE")
	return rc
}

func (rc *RequestChain) And(e Entry) *RequestChain {
	rc.entries = append(rc.entries, e)
	rc.operations = append(rc.operations, "AND")
	return rc
}

func (rc *RequestChain) Or(e Entry) *RequestChain {
	rc.entries = append(rc.entries, e)
	rc.operations = append(rc.operations, "Or")
	return rc
}

func (rc *RequestChain) Bake() string {
	conditions := ""
	for idx, op := range rc.operations {
		conditions += " " + op + " " + rc.entries[idx].Bake()
	}
	return conditions
}

func main() {

	e1 := Entry{"country_id", "IN", []int64{1, 2, 3, 4, 5}}
	e2 := Entry{"title", "LIKE", "Cheif %"}

	rc := &RequestChain{}
	rc.Where(e1).And(e2)

	// fmt.Println(rc.Bake())
	fmt.Printf("%s\n", rc.Bake())
	spew.Dump(rc)

}
