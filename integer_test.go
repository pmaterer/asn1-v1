package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultInt int

func TestEncodeInteger(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode 1 byte integer",
			expected:    []byte{0x2d},
			value:       45,
			errExpected: false,
		},
		{
			name:        "Test encode multi-byte integer",
			expected:    []byte{0x55, 0xcc},
			value:       21964,
			errExpected: false,
		},
		{
			name:        "Test encode multi-byte integer",
			expected:    []byte{0x0d, 0x17, 0x97, 0x31},
			value:       219649841,
			errExpected: false,
		},
		{
			name:        "Test encode negative integer",
			expected:    []byte{0xfd, 0x6f},
			value:       -657,
			errExpected: false,
		},
		{
			name:        "Test encode integer error",
			expected:    []byte{0xfd, 0x6f},
			value:       "55",
			errExpected: true,
		},
		{
			name:        "Test encode integer default",
			expected:    []byte{0x00},
			value:       defaultInt,
			errExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeInt(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
			}
		})
	}
}
