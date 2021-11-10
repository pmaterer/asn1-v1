package asn1

import (
	"reflect"
)

func encodeString(value reflect.Value) ([]byte, error) {
	if value.Kind() != reflect.String {
		return nil, invalidTypeError("string", value)
	}

	return []byte(value.String()), nil
}

// https://en.wikipedia.org/wiki/PrintableString
func isValidPrintableString(str string) bool {
	for _, c := range str {
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		default:
			switch c {
			case ' ', '\'', '(', ')', '+', ',', '-', '.', '/', ':', '=', '?':
			default:
				return false
			}
		}
	}
	return true
}
