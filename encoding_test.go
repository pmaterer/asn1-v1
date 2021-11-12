package asn1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		options     string
		errExpected bool
	}{
		{
			name:        "Test encode printable string",
			expected:    []byte{0x13, 0x02, 0x68, 0x69},
			value:       "hi",
			options:     "printable",
			errExpected: false,
		},
		{
			name:        "Test encode ia5string",
			expected:    []byte{0x16, 0x02, 0x68, 0x69},
			value:       "hi",
			options:     "ia5",
			errExpected: false,
		},
		{
			name:        "Test encode utf8string",
			expected:    []byte{0x0c, 0x04, 0xf0, 0x9f, 0x98, 0x8e},
			value:       "😎",
			options:     "utf8",
			errExpected: false,
		},
		{
			name:        "Test encode utc time",
			expected:    []byte{0x17, 0x11, 0x31, 0x39, 0x31, 0x32, 0x31, 0x35, 0x31, 0x39, 0x30, 0x32, 0x31, 0x30, 0x2d, 0x30, 0x38, 0x30, 0x30},
			value:       time.Date(2019, time.December, 15, 19, 02, 10, 0, time.FixedZone("UTC-8", -8*60*60)),
			options:     "utc",
			errExpected: false,
		},
		{
			name:        "Test encode oid",
			expected:    []byte{0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x0b},
			value:       ObjectIdentifier{Root: 1, Subidentifiers: []uint64{2, 840, 113549, 1, 1, 11}},
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode octet string",
			expected:    []byte{0x04, 0x04, 0x03, 0x02, 0x06, 0xa0},
			value:       []byte{0x03, 0x02, 0x06, 0xa0},
			options:     "",
			errExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := MarshalWithOptions(tt.value, tt.options)
			if tt.errExpected {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, enc)
			}

		})
	}
}
