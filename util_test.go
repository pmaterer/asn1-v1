package asn1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBase128(t *testing.T) {
	tests := []struct {
		name     string
		expected []byte
		value    uint64
	}{
		{
			name:     "Test base 128 encoding",
			expected: []byte{0x7f},
			value:    127,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0x81, 0x00},
			value:    128,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0xC0, 0x00},
			value:    8192,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0xFF, 0x7F},
			value:    16383,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0x81, 0x80, 0x00},
			value:    16384,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0xFF, 0xFF, 0x7F},
			value:    2097151,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0x81, 0x80, 0x80, 0x00},
			value:    2097152,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0xC0, 0x80, 0x80, 0x00},
			value:    134217728,
		},
		{
			name:     "Test base 128 encoding",
			expected: []byte{0xFF, 0xFF, 0xFF, 0x7F},
			value:    268435455,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := encodeBase128(tt.value)
			assert.Equal(t, tt.expected, b)
		})
	}
}
