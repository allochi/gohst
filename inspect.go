package gohst

import (
	"reflect"
)

type Kind uint

// Pointer are odds
const (
	Nil         = 0
	Pointer     = 1
	Primitive   = 2
	Struct      = 4
	Slice       = 8
	OfPrimitive = 16
	OfStruct    = 32

	Pointer2Primitive        = Pointer | Primitive
	Pointer2Struct           = Pointer | Struct
	Pointer2Slice            = Pointer | Slice
	Pointer2SliceOfPrimitive = Pointer | Slice | OfPrimitive
	Pointer2SliceOfStruct    = Pointer | Slice | OfStruct

	SliceOfPrimitive = Slice | OfPrimitive
	SliceOfStruct    = Slice | OfStruct
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
		kind |= Pointer
	} else {
		_type = _value.Type()
	}

	// Is it a slice or a struct?
	_kind := _type.Kind()
	var _elemKind reflect.Kind
	if _kind == reflect.Slice {
		_elemKind = _type.Elem().Kind()
		kind |= Slice
	} else if _kind == reflect.Struct {
		kind |= Struct
	}

	// Is it a slice of struct or primitives
	if (_elemKind > reflect.Invalid) && (_elemKind < reflect.Array || _elemKind == reflect.String) {
		if kind == Slice || kind == Pointer2Slice {
			kind |= OfPrimitive
		}
	} else if _elemKind == reflect.Struct {
		if kind == Slice || kind == Pointer2Slice {
			kind |= OfStruct
		}
	}

	return
}
