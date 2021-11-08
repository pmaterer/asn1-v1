package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var (
// 	testOIDValid = &ObjectIdentifier{0, []uint64{0, 9, 2342, 19200300, 100, 1, 3}}
// )

func TestNewObjectIdentifier(t *testing.T) {

	// valid
	_, err := NewObjectIdentifier(0, []uint64{9, 2342, 19200300, 100, 1, 2})
	assert.Nil(t, err)

	// bad root node
	_, err = NewObjectIdentifier(5, []uint64{9, 2342, 19200300, 100, 1, 2})
	assert.Error(t, err)

	// empty node
	_, err = NewObjectIdentifier(0, []uint64{})
	assert.Error(t, err)
}

func TestObjectIdentifierToString(t *testing.T) {
	oid, _ := NewObjectIdentifier(0, []uint64{9, 2342, 19200300, 100, 1, 1})
	oidStr := oid.ToString()

	assert.Equal(t, "0.9.2342.19200300.100.1.1", oidStr)
}

func TestEncodeObjectIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode OID",
			expected:    []byte{0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x0b},
			value:       ObjectIdentifier{1, []uint64{2, 840, 113549, 1, 1, 11}},
			errExpected: false,
		},
		{
			name:        "Test encode OID error",
			expected:    []byte{},
			value:       "1.2.840",
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeObjectIdentifier(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
				assert.Nil(t, err)
			}
		})
	}
}
