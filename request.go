package gohst

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Entry struct {
	Field    string
	Operator string
	Values   interface{}
}

func (e *Entry) Bake() string {

	// == Work the field
	field := e.Field
	switch field {
	case "Id":
		field = "id"
	case "CreatedAt":
		field = "created_at"
	case "UpdatedAt":
		field = "updated_at"
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

	// isSlice := gohst.KindOf(e.Values) == gohst.SliceOfPrimitive
	isSlice := reflect.TypeOf(e.Values).Kind() == reflect.Slice
	var _type reflect.Type

	if isSlice {
		_type = reflect.TypeOf(e.Values).Elem()
	} else {
		_type = reflect.TypeOf(e.Values)
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
		slice := reflect.ValueOf(e.Values)
		values = make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			values[i] = process(slice.Index(i).Interface())
		}
		value = fmt.Sprintf("(%s)", strings.Join(values, ","))
	} else {
		value = process(e.Values)
	}

	return field + " " + e.Operator + " " + value
}

type Requester interface {
	// Where(Entry) *Requester
	// And(Entry) *Requester
	// Or(Entry) *Requester
	// WhereGroup(Entry) *Requester
	// AndGroup(Entry) *Requester
	// OrGroup(Entry) *Requester
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
