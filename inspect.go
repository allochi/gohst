package gohst

import (
	"reflect"
)

func IsPtr2Slice(i interface{}) bool {
	_value := reflect.ValueOf(i)
	return (_value.Type().Kind() == reflect.Ptr) && (_value.Elem().Type().Kind() == reflect.Slice)
}

func IsPtr2Struct(i interface{}) bool {
	_value := reflect.ValueOf(i)
	return (_value.Type().Kind() == reflect.Ptr) && (_value.Elem().Type().Kind() == reflect.Struct)
}

func IsStructOrPtr2Struct(i interface{}) bool {
	_value := reflect.ValueOf(i)
	return (_value.Type().Kind() == reflect.Struct) || ((_value.Type().Kind() == reflect.Ptr) && (_value.Elem().Type().Kind() == reflect.Struct))
}

func IsPtr2SliceOfStruct(i interface{}) (result bool) {
	_value := reflect.ValueOf(i)
	if _value.Type().Kind() == reflect.Ptr {
		_elemType := _value.Elem().Type()
		if (_elemType.Kind() == reflect.Slice) && (_elemType.Elem().Kind() == reflect.Struct) {
			result = true
		}
	}
	return
}

func IsPtr2SliceOfPrimitive(i interface{}) (result bool) {
	_value := reflect.ValueOf(i)
	if _value.Type().Kind() == reflect.Ptr {
		_elem := _value.Elem()
		if _elem.Type().Kind() == reflect.Slice {
			_elemKind := _value.Elem().Type().Elem().Kind()
			if (_elemKind > reflect.Invalid) && (_elemKind < reflect.Array || _elemKind == reflect.String) {
				result = true
			}
		}
	}
	return
}

func IsSliceOrPtr2SliceOfPrimitive(i interface{}) (result bool) {
	_value := reflect.ValueOf(i)
	if _value.Type().Kind() == reflect.Ptr {
		_elem := _value.Elem()
		if _elem.Type().Kind() == reflect.Slice {
			_elemKind := _value.Elem().Type().Elem().Kind()
			if (_elemKind > reflect.Invalid) && (_elemKind < reflect.Array || _elemKind == reflect.String) {
				result = true
			}
		}
	} else if _value.Type().Kind() == reflect.Slice {
		_elemKind := _value.Type().Elem().Kind()
		if (_elemKind > reflect.Invalid) && (_elemKind < reflect.Array || _elemKind == reflect.String) {
			result = true
		}
	}
	return
}
