package asn1

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func isValidBitString(str string) bool {
	if str[0] != 'b' {
		return false
	}
	for _, c := range str[1:] {
		if !(c == '0' || c == '1') {
			return false
		}
	}
	return true
}

func EncodeBitString(value reflect.Value) ([]byte, error) {
	if value.Kind() != reflect.String {
		return nil, invalidTypeError("string", value)
	}

	bitStr := value.String()
	if !isValidBitString(bitStr) {
		return nil, fmt.Errorf("%s not a valid bit string", bitStr)
	}
	bitStr = bitStr[1:]

	bitLength := float64(len(bitStr))
	paddedBitLength := uint(8 * (math.Ceil(bitLength / 8.0)))
	unusedBits := paddedBitLength - uint(bitLength)
	octetLength := paddedBitLength / 8

	buf := new(bytes.Buffer)
	buf.WriteByte(byte(unusedBits))

	padding := strings.Repeat("0", int(unusedBits))
	paddedBitStr := bitStr + padding

	for i := 1; i <= int(octetLength); i++ {
		index := i * 8

		parsed, err := strconv.ParseInt(paddedBitStr[index-8:index], 2, 64)
		if err != nil {
			return nil, err
		}

		enc := encodeUint(uint64(parsed))
		buf.Write(enc)
	}

	return buf.Bytes(), nil
}
