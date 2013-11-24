package gohst

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Clause struct {
	Field    string
	Operator string
	Values   interface{}
}

func (e *Clause) Bake(object interface{}) (caluse string) {

	// Work the field
	field := e.Field
	switch field {
	case "Id":
		field = "id"
	case "CreatedAt":
		field = "created_at"
	case "UpdatedAt":
		field = "updated_at"
	default:
		if strings.Contains(field, "::") {
			fieldParts := strings.Split(field, "::")
			field = fmt.Sprintf("(data->>'%s')::%s", fieldParts[0], fieldParts[1])
		} else {
			_type := TypeOf(object, field)
			field = JsonName(object, field)
			switch _type {
			case "[]string":
				field = fmt.Sprintf("_array(data,'%s')", field)
			case "time.Time":
				field = fmt.Sprintf("_date(data,'%s')", field)
			case "int":
				field = fmt.Sprintf("(data->>'%s')::int", field)
			default:
				field = fmt.Sprintf("data->>'%s'", field)
			}
		}
	}

	// Work the values
	var value string

	isSlice := reflect.TypeOf(e.Values).Kind() == reflect.Slice

	if isSlice {
		value = IN(e.Values)
	} else {
		value = process(e.Values)
	}

	if e.Operator == "<@" || e.Operator == "@>" {
		caluse = fmt.Sprintf("%s %s %s", value, "<@", field)
	} else {
		caluse = fmt.Sprintf("%s %s %s", field, e.Operator, value)
	}

	return
}

func process(value interface{}) string {

	var _type reflect.Type
	_type = reflect.TypeOf(value)

	_kind := _type.Kind()

	switch {
	case _kind >= reflect.Int && _kind <= reflect.Float64:
		return fmt.Sprintf("%v", value)
	case _kind == reflect.String:
		if strings.Contains(value.(string), "$") {
			return fmt.Sprintf("%v", value)
		} else {
			return fmt.Sprintf("'%v'", value)
		}
	case _kind == reflect.Bool:
		return fmt.Sprintf("%t", value)
	case _type == reflect.TypeOf(time.Now()):
		return "'" + value.(time.Time).Format(time.RFC3339) + "'"
	case _kind == reflect.Struct:
		_id := reflect.ValueOf(value).FieldByName("Id")
		if _id.Kind() == reflect.Invalid {
			return ""
		} else {
			return fmt.Sprintf("%v", _id.Interface())
		}
	}

	return ""

}

type Requester interface {
	Bake(object interface{}) string
}

// Prepared Statement Caller
type PreparedStatement struct {
	name string
}

func (ps *PreparedStatement) Bake(object interface{}) string {
	return ps.name
}

// SQL Request Builder
type RequestChain struct {
	clauses     []Clause
	operations  []string
	isGroupOpen bool
}

func (rc *RequestChain) Where(e Clause) *RequestChain {
	rc.clauses = append(rc.clauses, e)
	rc.operations = append(rc.operations, "WHERE")
	return rc
}

func (rc *RequestChain) And(e Clause) *RequestChain {
	rc.clauses = append(rc.clauses, e)
	rc.operations = append(rc.operations, "AND")
	return rc
}

func (rc *RequestChain) Or(e Clause) *RequestChain {
	rc.clauses = append(rc.clauses, e)
	rc.operations = append(rc.operations, "Or")
	return rc
}

func (rc *RequestChain) WhereGroup(e Clause) *RequestChain {
	rc.clauses = append(rc.clauses, e)
	rc.operations = append(rc.operations, "WHERE (")
	rc.isGroupOpen = true
	return rc
}

func (rc *RequestChain) AndGroup(e Clause) *RequestChain {
	rc.clauses = append(rc.clauses, e)
	if rc.isGroupOpen {
		rc.operations = append(rc.operations, ") AND (")
	} else {
		rc.operations = append(rc.operations, "AND (")
		rc.isGroupOpen = true
	}
	return rc
}

func (rc *RequestChain) OrGroup(e Clause) *RequestChain {
	rc.clauses = append(rc.clauses, e)
	if rc.isGroupOpen {
		rc.operations = append(rc.operations, ") OR (")
	} else {
		rc.operations = append(rc.operations, "OR (")
		rc.isGroupOpen = true
	}
	return rc
}

func (rc *RequestChain) Bake(object interface{}) string {
	conditions := ""
	for idx, op := range rc.operations {
		conditions += fmt.Sprintf("%s %s", op, rc.clauses[idx].Bake(object))
	}

	if rc.isGroupOpen {
		conditions += ")"
	}

	return conditions
}

// This takes a slice and an optional boolean, the slice get converted to a sql list,
// if the boolean is based as true, then the list will be constructed as an sql array,
// if any other type is based than a slice and empty string is returned.
func IN(_slice interface{}, params ...bool) (values string) {

	_type := reflect.TypeOf(_slice)
	_typeName := _type.String()
	if !strings.Contains(_typeName, "[]") {
		return
	}

	// Time thus convert to a slice of strings of time.RFC3339
	if _typeName == "[]time.Time" {
		times := _slice.([]time.Time)
		fomats := make([]string, len(times))
		for index, value := range times {
			fomats[index] = value.Format(time.RFC3339)
		}
		_slice = fomats
	}

	// Objects thus get their Ids as []int64
	if _typeName != "[]time.Time" && _type.Elem().Kind() == reflect.Struct {
		length := reflect.ValueOf(_slice).Len()
		ids := make([]interface{}, length)
		for index := 0; index < length; index++ {
			object := reflect.ValueOf(_slice).Index(index)
			ids[index] = object.FieldByName("Id").Interface()
		}
		_slice = ids
	}

	isArray := len(params) > 0 && params[0]

	values = fmt.Sprintf("%v", _slice)
	values = strings.Trim(values, "[]")

	seperator := "','"
	surrounder := "'%s'"
	brackets := "(%s)"
	if isArray {
		seperator = "\",\""
		surrounder = "\"%s\""
		brackets = "{%s}"
	}

	if _typeName != "[]string" && _typeName != "[]time.Time" {
		seperator = ","
		surrounder = "%s"
	}

	values = strings.Replace(values, " ", seperator, -1)
	values = fmt.Sprintf(surrounder, values)
	values = fmt.Sprintf(brackets, values)

	return

}
