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

	// == Work the values
	var values []string
	var value string

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
			if strings.Contains(i.(string), "$") {
				return fmt.Sprintf("%v", i)
			} else {
				return fmt.Sprintf("'%v'", i)
			}
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

	if e.Operator == "<@" {
		caluse = fmt.Sprintf("%s %s %s", value, e.Operator, field)
	} else {
		caluse = fmt.Sprintf("%s %s %s", field, e.Operator, value)
	}

	return
}

type Requester interface {
	Bake(object interface{}) string
}

// == Prepared Statement Caller
type PreparedStatement struct {
	name string
}

func (ps *PreparedStatement) Bake(object interface{}) string {
	return ps.name
}

// == SQL Request Builder
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

func Ints2String(ids []int64) (values string) {
	values = fmt.Sprintf("%v", ids)
	values = strings.Trim(values, "[]")
	values = strings.Replace(values, " ", ",", -1)
	return
}
