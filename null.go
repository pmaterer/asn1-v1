package asn1

import "reflect"

var nullType = reflect.TypeOf(Null{})

type Null struct{}

func encodeNull(value reflect.Value) ([]byte, error) {
	if value.Type() != nullType {
		return nil, invalidTypeError("Null", value)
	}
	return nil, nil
}
