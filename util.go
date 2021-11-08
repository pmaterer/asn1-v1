package asn1

import "strings"

func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func reverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < len(b)/2; i++ {
		b[i], b[j-i] = b[j-i], b[i]
	}
	return b
}

func encodeUint(n uint64) []byte {
	length := uintLength(n)
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		shift := uint((length - 1 - i) * 8)
		buf[i] = byte(n >> int(shift))
	}
	return buf
}

func uintLength(i uint64) (length int) {
	length = 1
	for i > 255 {
		length++
		i >>= 8
	}
	return
}
