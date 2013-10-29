package main

import (
	// "allochi/gohst"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"reflect"
	"strings"
	"time"
)

type Entry struct {
	field    string
	operator string
	values   interface{}
}

func (e *Entry) Bake() string {

	// == Work the field
	field := e.field
	switch field {
	case "Id":
		fallthrough
	case "CreatedAt":
		fallthrough
	case "UpdatedAt":
		// == Direct fields don't need "data->>'...'"
		field = fmt.Sprintf("'%s'", field)
	default:
		fields := strings.Split(field, ":")
		if len(fields) > 1 {
			// == SQL casting included
			field = fmt.Sprintf("(data->>'%s')::%s", fields[0], fields[1])
		} else {
			// == Use default
			field = fmt.Sprintf("data->>'%s'", field)
		}
	}

	// == Work the values
	var values []string
	var value string

	// isSlice := gohst.KindOf(e.values) == gohst.SliceOfPrimitive
	isSlice := reflect.TypeOf(e.values).Kind() == reflect.Slice
	var _type reflect.Type

	if isSlice {
		_type = reflect.TypeOf(e.values).Elem()
	} else {
		_type = reflect.TypeOf(e.values)
	}

	_kind := _type.Kind()

	var process func(interface{}) string

	switch {
	case _kind >= reflect.Int && _kind <= reflect.Float64:
		process = func(i interface{}) string {
			return fmt.Sprintf("%v", i)
		}
	case _kind == reflect.String:
		process = func(i interface{}) string {
			return fmt.Sprintf("'%v'", i)
		}
	case _kind == reflect.Bool:
		process = func(i interface{}) string {
			return fmt.Sprintf("%t", i)
		}
	case _type == reflect.TypeOf(time.Now()):
		process = func(i interface{}) string {
			return "'" + i.(time.Time).Format("2006-01-02 15:04:05") + "'"
		}
	case _kind == reflect.Struct:
		process = func(i interface{}) string {
			_id := reflect.ValueOf(i).FieldByName("Id")
			if _id.Kind() == reflect.Invalid {
				return ""
			} else {
				return fmt.Sprintf("%v", _id.Interface())
			}
		}
	}

	if isSlice {
		slice := reflect.ValueOf(e.values)
		values = make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			values[i] = process(slice.Index(i).Interface())
		}
		value = fmt.Sprintf("(%s)", strings.Join(values, ","))
	} else {
		value = process(e.values)
	}

	return field + " " + e.operator + " " + value
}

type Requester interface {
	Where(Entry) *Requester
	And(Entry) *Requester
	Or(Entry) *Requester
	WhereGroup(Entry) *Requester
	AndGroup(Entry) *Requester
	OrGroup(Entry) *Requester
	Bake() string
}

type RequestChain struct {
	entries     []Entry
	operations  []string
	isGroupOpen bool
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

func (rc *RequestChain) WhereGroup(e Entry) *RequestChain {
	rc.entries = append(rc.entries, e)
	rc.operations = append(rc.operations, "WHERE (")
	rc.isGroupOpen = true
	return rc
}

func (rc *RequestChain) AndGroup(e Entry) *RequestChain {
	rc.entries = append(rc.entries, e)
	if rc.isGroupOpen {
		rc.operations = append(rc.operations, ") AND (")
	} else {
		rc.operations = append(rc.operations, "AND (")
		rc.isGroupOpen = true
	}
	return rc
}

func (rc *RequestChain) OrGroup(e Entry) *RequestChain {
	rc.entries = append(rc.entries, e)
	if rc.isGroupOpen {
		rc.operations = append(rc.operations, ") OR (")
	} else {
		rc.operations = append(rc.operations, "OR (")
		rc.isGroupOpen = true
	}
	return rc
}

func (rc *RequestChain) Bake() string {
	conditions := ""
	for idx, op := range rc.operations {
		conditions += fmt.Sprintf(" %s %s", op, rc.entries[idx].Bake())
	}

	if rc.isGroupOpen {
		conditions += ")"
	}

	return conditions
}

type Person struct {
	Id   int64
	Name string
}

func main() {

	asyadi := []Person{
		Person{1, "Mohammad"},
		Person{2, "Ali"},
		Person{3, "Hasan"},
		Person{4, "Husain"},
		Person{12, "Hujjah"},
	}

	e1 := Entry{"country_id:int", "IN", []int64{56, 79, 28, 90, 10}}
	e2 := Entry{"title", "LIKE", "Cheif %"}
	e3 := Entry{"colors:text", "IN", []string{"Red", "Green", "Blue"}}
	e4 := Entry{"days:datetime", "IN", []time.Time{time.Now(), time.Now().AddDate(-1, 0, 0)}}
	e5 := Entry{"Imam", "IN", asyadi}

	rc := &RequestChain{}
	// rc.Where(e1).And(e2).And(e3)
	rc.Where(e1).And(e2).And(e3).Or(e4).And(e5)

	request := []string{
		"(", e1.Bake(), "AND", e2.Bake(), ")",
		"OR",
		"(", e3.Bake(), "AND", e4.Bake(), ")",
	}

	fmt.Printf("%s\n\n", rc.Bake())
	fmt.Printf("%s\n\n", strings.Join(request, " "))

	rc2 := &RequestChain{}
	rc2.WhereGroup(e1).And(e2).AndGroup(e3).Or(e4).AndGroup(e5)
	fmt.Printf("%s\n\n", rc2.Bake())

	contactsRQ := &RequestChain{}
	contactsRQ.Where(Entry{"country_id:int", "=", 20}).And(Entry{"title_id:int", "=", 4}).And(Entry{"telephone", "=", ""})
	fmt.Printf("%s\n\n", contactsRQ.Bake())

	// fmt.Println(rc.Bake())
	// spew.Dump(rc)

}
