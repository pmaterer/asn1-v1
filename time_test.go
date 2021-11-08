package asn1

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncodeUTCTime(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encoding",
			expected:    []byte{0x32, 0x31, 0x30, 0x32, 0x32, 0x31, 0x30, 0x31, 0x31, 0x30, 0x33, 0x30, 0x5a},
			value:       time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			errExpected: false,
		},
		{
			name:        "Test encoding error",
			expected:    []byte{},
			value:       0,
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeUTCTime(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
			}

		})
	}
}

func TestEncodeGeneralizedTime(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		errExpected bool
	}{
		{
			name:        "Test encoding",
			expected:    []byte{0x32, 0x30, 0x32, 0x31, 0x30, 0x32, 0x32, 0x31, 0x30, 0x31, 0x31, 0x30, 0x33, 0x30, 0x5a},
			value:       time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			errExpected: false,
		},
		{
			name:        "Test encoding error",
			expected:    []byte{},
			value:       0,
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := encodeGeneralizedTime(reflect.ValueOf(tt.value))
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, b)
			}
		})
	}
}
