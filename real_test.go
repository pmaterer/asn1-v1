package asn1

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeReal(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode real",
			expected:    []byte{0x02, 0x33, 0x2e, 0x31, 0x34, 0x35, 0x39},
			value:       3.1459,
			errExpected: false,
		},
		{
			name:        "Test encode infinity",
			expected:    []byte{0x40},
			value:       math.Inf(1),
			errExpected: false,
		},
		{
			name:        "Test encode -infinity",
			expected:    []byte{0x41},
			value:       math.Inf(-1),
			errExpected: false,
		},
		{
			name:        "Test encode NaN",
			expected:    []byte{0x42},
			value:       math.Log(-1.0),
			errExpected: false,
		},
		{
			name:        "Test encode -0",
			expected:    []byte{0x43},
			value:       math.Copysign(0, -1),
			errExpected: false,
		},
		{
			name:        "Test encode real error",
			expected:    []byte{},
			value:       4,
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeReal(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
			}
		})
	}
}
