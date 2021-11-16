package asn1

import (
	"bytes"
	"reflect"
)

func encodeList(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return nil, invalidTypeError("array/slice", v)
	}

	buf := new(bytes.Buffer)
	for i := 0; i < v.Len(); i++ {
		b, err := Marshal(v.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}
