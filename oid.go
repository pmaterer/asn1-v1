package asn1

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

var (
	oidType = reflect.TypeOf(ObjectIdentifier{})
)

type ObjectIdentifier struct {
	Root           uint8
	Subidentifiers []uint64
}

func isValidRootNode(root uint8) bool {
	switch root {
	case 0, 1, 2:
		return true
	default:
		return false
	}
}

func NewObjectIdentifier(root uint8, node []uint64) (*ObjectIdentifier, error) {
	if !isValidRootNode(root) {
		return nil, fmt.Errorf("error creating oid: invalid root node")
	}
	if len(node) == 0 {
		return nil, fmt.Errorf("error creating oid: empty node")
	}

	switch root {
	case 0, 1:
		if node[0] >= 40 {
			return nil, fmt.Errorf("invalid node")
		}
	}
	return &ObjectIdentifier{root, node}, nil
}

func (o *ObjectIdentifier) ToString() string {
	oidStr := strconv.Itoa(int(o.Root)) + "."
	for j, i := range o.Subidentifiers {
		oidStr += strconv.Itoa(int(i))
		if j+1 < len(o.Subidentifiers) {
			oidStr += "."
		}
	}
	return oidStr
}

func encodeObjectIdentifier(value reflect.Value) ([]byte, error) {
	if value.Type() != oidType {
		return nil, invalidTypeError("ObjectIdentifier", value)
	}

	oid := value.Interface().(ObjectIdentifier)

	b := new(bytes.Buffer)

	initialID := int(oid.Root*40 + uint8(oid.Subidentifiers[0]))
	initialIDEnc, err := encodeInt(reflect.ValueOf(initialID))

	if err != nil {
		return nil, err
	}
	b.Write(initialIDEnc)

	for i := 1; i < len(oid.Subidentifiers); i++ {
		sid := oid.Subidentifiers[i]
		sidEnc := encodeSubidentifier(sid)
		if err != nil {
			return nil, err
		}
		b.Write(sidEnc)
	}

	return b.Bytes(), nil
}

func encodeSubidentifier(sid uint64) []byte {
	buf := new(bytes.Buffer)

	for sid != 0 {
		i := sid & 0x7f
		sid >>= 7

		// bit 8 = 1 in every byte except last
		if len(buf.Bytes()) != 0 {
			i |= 0x80
		}
		buf.WriteByte(byte(i))
	}

	return reverseBytes(buf.Bytes())
}
