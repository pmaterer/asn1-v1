package asn1

import "reflect"

func encodeString(value reflect.Value) ([]byte, error) {
	if value.Kind() != reflect.String {
		return nil, invalidTypeError("string", value)
	}

	return []byte(value.String()), nil
}
