package asn1

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
