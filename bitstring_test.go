package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBitString(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode bit string",
			expected:    []byte{0x07, 0x68, 0x80},
			value:       "b011010001",
			errExpected: false,
		},
		{
			name:        "Test encode bit string",
			expected:    []byte{0x06, 0x6e, 0x5d, 0xc0},
			value:       "b011011100101110111",
			errExpected: false,
		},
		{
			name:        "Test encode bit string error",
			expected:    []byte{},
			value:       "0x9a",
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := EncodeBitString(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
				assert.Nil(t, err)
			}
		})
	}
}

func TestIsValidBitString(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		value    string
	}{
		{
			name:     "Test encode proper bit string",
			expected: true,
			value:    "b0101101011",
		},
		{
			name:     "Test encode bad bit string",
			expected: false,
			value:    "01011a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := isValidBitString(tt.value)
			assert.Equal(t, tt.expected, ok)
		})
	}
}
