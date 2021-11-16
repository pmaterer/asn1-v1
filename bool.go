package asn1

import "reflect"

var (
	asn1True  byte = 0xff
	asn1False byte = 0x00
)

func encodeBool(value reflect.Value) ([]byte, error) {
	switch value.Kind() {
	case reflect.Bool:
	default:
		return nil, invalidTypeError("bool", value)
	}

	if value.Bool() {
		return []byte{asn1True}, nil
	}
	return []byte{asn1False}, nil
}
