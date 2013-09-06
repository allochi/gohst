package gohst

import (
	"reflect"
)

type Kind uint

// Pointer are odds
const (
	Nil Kind = iota
	Pointer
	Struct
	Pointer2Struct
	Slice
	Pointer2Slice
	SliceOfPrimitive
	Pointer2SliceOfPrimitive
	SliceOfStruct
	Pointer2SliceOfStruct
)

func KindOf(i interface{}) (kind Kind) {
	_value := reflect.ValueOf(i)

	// Is it nil?
	if _value.Kind() == reflect.Invalid {
		return
	}

	// Is it a pointer?
	var _type reflect.Type
	if _value.Type().Kind() == reflect.Ptr {
		_type = _value.Elem().Type()
		kind += Pointer
	} else {
		_type = _value.Type()
	}

	// Is it a slice or a struct?
	_kind := _type.Kind()
	var _elemKind reflect.Kind
	if _kind == reflect.Slice {
		_elemKind = _type.Elem().Kind()
		kind += Slice
	} else if _kind == reflect.Struct {
		kind += Struct
	}

	// Is it a slice of struct or primitives
	if (_elemKind > reflect.Invalid) && (_elemKind < reflect.Array || _elemKind == reflect.String) {
		if kind == Slice || kind == Pointer2Slice {
			kind += 2
		}
	} else if _elemKind == reflect.Struct {
		if kind == Slice || kind == Pointer2Slice {
			kind += 4
		}
	}

	return
}
