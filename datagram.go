package asn1

import "bytes"

// identifier represents the identifier octets of a BER datagram
type identifier struct {
	class         TagClass
	isConstructed bool
	tag           Tag
}

// datagram represents a BER encoded data structure
type datagram struct {
	identifier
	length uint
	body   []byte
}

func newDatagram() *datagram {
	return &datagram{
		identifier: identifier{
			class:         TagClassUniversal,
			isConstructed: false,
		},
		length: 0,
		body:   nil,
	}
}

func (d *datagram) encode() []byte {
	buf := new(bytes.Buffer)

	buf.Write(d.identifier.encode())
	d.length = uint(len(d.body))
	buf.Write(encodeLength(uint(d.length)))
	buf.Write(d.body)
	return buf.Bytes()
}

func (i *identifier) encode() []byte {
	b := []byte{0x00}

	b[0] |= byte(i.class << 6)

	if i.isConstructed {
		b[0] |= byte(1 << 5)
	} else {
		b[0] |= byte(0 << 5)
	}

	// universal tags 0-30
	if i.tag <= 30 {
		b[0] |= byte(i.tag)
	} else {
		b[0] |= byte(0x1f)
		b = append(b, encodeBase128(uint64(i.tag))...)
	}

	return b
}

func encodeLength(length uint) []byte {
	// only definite form supported
	// length encoded as unsigned binary integers

	lengthBytes := encodeUint(uint64(length))

	// short form
	if length <= 0x7f {
		return lengthBytes
	}

	// long form
	b := new(bytes.Buffer)
	header := len(lengthBytes) | 0x80

	b.Write(encodeUint(uint64(header)))
	b.Write(lengthBytes)
	return b.Bytes()
}
