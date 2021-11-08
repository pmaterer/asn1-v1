package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOctetString(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode octet string slice",
			expected:    []byte{0x22, 0xff, 0xa2, 0x02},
			value:       []byte{0x22, 0xff, 0xa2, 0x02},
			errExpected: false,
		},
		{
			name:        "Test encode octet string array",
			expected:    []byte{0x22, 0xff, 0xa2, 0x02},
			value:       [4]byte{0x22, 0xff, 0xa2, 0x02},
			errExpected: false,
		},
		{
			name:        "Test encode octet string error",
			expected:    []byte{0x22, 0xff, 0xa2, 0x02},
			value:       []int{1, 20, 5},
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeOctetString(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
				assert.Nil(t, err)
			}
		})
	}
}
