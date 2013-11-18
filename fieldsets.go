package gohst

import (
	"reflect"
)

var FieldSets map[string]FieldTypes

type FieldTypes struct {
	FieldType map[string]string
	JsonType  map[string]string
	JsonName  map[string]string
}

var SQLTypes map[string]string

func init() {
	FieldSets = make(map[string]FieldTypes)
	SQLTypes = map[string]string{
		"bool":        "bool",
		"int":         "integer",
		"int8":        "smallint",
		"int16":       "smallint",
		"int32":       "integer",
		"int64":       "integer",
		"uint":        "integer",
		"uint8":       "integer",
		"uint16":      "integer",
		"uint32":      "bigint",
		"uint64":      "bigint",
		"float32":     "float4",
		"float64":     "float8",
		"string":      "text",
		"time.Time":   "timestamp",
		"[]bool":      "bool[]",
		"[]int":       "integer[]",
		"[]int8":      "smallint[]",
		"[]int16":     "smallint[]",
		"[]int32":     "integer[]",
		"[]int64":     "integer[]",
		"[]uint":      "integer[]",
		"[]uint8":     "integer[]",
		"[]uint16":    "integer[]",
		"[]uint32":    "bigint[]",
		"[]uint64":    "bigint[]",
		"[]float32":   "float4[]",
		"[]float64":   "float8[]",
		"[]string":    "text[]",
		"[]time.Time": "timestamp[]",
	}
}

func ScanFields(object interface{}) {

	_typeName, _type, _value := TypeName(object)

	FieldSets[_typeName] = FieldTypes{map[string]string{}, map[string]string{}, map[string]string{}}

	for i := 0; i < _value.NumField(); i++ {
		_field := _type.Field(i)
		_fieldType := _field.Type.String()
		_name := _field.Name
		_jsonName := _field.Tag.Get("json")

		FieldSets[_typeName].FieldType[_name] = _fieldType
		if !(_jsonName == "-" || _jsonName == "") {
			FieldSets[_typeName].JsonName[_name] = _jsonName
			FieldSets[_typeName].JsonType[_jsonName] = _fieldType
		}

	}
}

func TypeOf(object interface{}, field string) string {

	_typeName, _, _ := TypeName(object)

	if _fieldSets, ok := FieldSets[_typeName]; ok {
		if _type, ok := _fieldSets.FieldType[field]; ok {
			return _type
		} else {
			return _fieldSets.JsonType[field]
		}
	} else {
		ScanFields(object)
		return TypeOf(object, field)
	}

}

// This function returns the json name of a field.
// The Object has to be scanned before this function can return anything,
// TypeOf() or ScanFields() should be called at least once before JsonName()
func JsonName(object interface{}, field string) string {
	_typeName, _, _ := TypeName(object)
	if _, ok := FieldSets[_typeName]; ok {
		_json := FieldSets[_typeName].JsonName[field]
		if _json != "" {
			return _json
		} else {
			return field
		}
	} else {
		ScanFields(object)
		return JsonName(object, field)
	}

}

func TypeName(object interface{}) (_name string, _type reflect.Type, _value reflect.Value) {

	// var _value reflect.Value
	// var _type reflect.Type
	_kind := KindOf(object)
	switch _kind {
	case Struct:
		_value = reflect.ValueOf(object)
		_type = _value.Type()
	case Pointer2Struct:
		_value = reflect.Indirect(reflect.ValueOf(object))
		_type = _value.Type()
	case SliceOfStruct:
		_value = reflect.ValueOf(object)
		_type = _value.Type().Elem()
	case Pointer2SliceOfStruct:
		_value = reflect.Indirect(reflect.ValueOf(object))
		_type = reflect.TypeOf(object).Elem().Elem()
	}

	_name = _type.Name()

	return

}
