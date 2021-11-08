package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultBool bool

func TestEncodeBool(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode bool",
			expected:    []byte{0xff},
			value:       true,
			errExpected: false,
		},
		{
			name:        "Test encode bool error",
			expected:    []byte{0xff},
			value:       66,
			errExpected: true,
		},
		{
			name:        "Test encode bool default",
			expected:    []byte{0x00},
			value:       defaultBool,
			errExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeBool(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
				assert.Nil(t, err)
			}
		})
	}
}
