package asn1

import (
	"reflect"
)

func empty(value reflect.Value) bool {
	defaultValue := reflect.Zero(value.Type())
	return reflect.DeepEqual(value.Interface(), defaultValue.Interface())
}
