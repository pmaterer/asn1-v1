package asn1

import "bytes"

func reverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < len(b)/2; i++ {
		b[i], b[j-i] = b[j-i], b[i]
	}
	return b
}

// https://en.wikipedia.org/wiki/Variable-length_quantity
func encodeBase128(num uint64) []byte {
	buf := new(bytes.Buffer)

	for num != 0 {
		i := num & 0x7f
		num >>= 7

		if len(buf.Bytes()) != 0 {
			i |= 0x80
		}
		buf.WriteByte(byte(i))
	}

	return reverseBytes(buf.Bytes())
}
