package asn1

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeList(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encode list slice",
			expected:    []byte{0x0c, 0x01, 0x61, 0x0c, 0x01, 0x62, 0x0c, 0x01, 0x63},
			value:       []string{"a", "b", "c"},
			errExpected: false,
		},
		{
			name:        "Test encode list slice",
			expected:    []byte{0x02, 0x01, 0x01, 0x02, 0x01, 0x02, 0x02, 0x01, 0x03},
			value:       []int{1, 2, 3},
			errExpected: false,
		},
		{
			name:        "Test encode list array",
			expected:    []byte{0x0c, 0x01, 0x61, 0x0c, 0x01, 0x62, 0x0c, 0x01, 0x63},
			value:       [3]string{"a", "b", "c"},
			errExpected: false,
		},
		{
			name:        "Test encode list array",
			expected:    []byte{0x02, 0x01, 0x01, 0x02, 0x01, 0x02, 0x02, 0x01, 0x03},
			value:       [3]int{1, 2, 3},
			errExpected: false,
		},
		{
			name:        "Test encode list empty",
			expected:    nil,
			value:       []int{},
			errExpected: false,
		},
		{
			name:        "Test encode list error",
			expected:    []byte{},
			value:       5,
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeList(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
			}
		})
	}
}
