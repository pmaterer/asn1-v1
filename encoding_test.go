package asn1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SequenceTest struct {
	OctetStringValue []byte
	BooleanValue     bool
	IntegerValue     int
}

type AnotherSequenceTest struct {
	Algorithm  ObjectIdentifier
	Parameters Null
}

func TestEncoding(t *testing.T) {
	tests := []struct {
		name        string
		expected    []byte
		value       interface{}
		options     string
		errExpected bool
	}{
		{
			name:        "Test encode bool",
			expected:    []byte{0x01, 0x01, 0xff},
			value:       true,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode bool",
			expected:    []byte{0x01, 0x01, 0x00},
			value:       false,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode printable string",
			expected:    []byte{0x13, 0x02, 0x68, 0x69},
			value:       "hi",
			options:     "printable",
			errExpected: false,
		},
		{
			name:        "Test encode printable string",
			expected:    []byte{0x13, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x20, 0x55, 0x73, 0x65, 0x72, 0x20, 0x31},
			value:       "Test User 1",
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
			name:        "Test encode ia5string",
			expected:    []byte{0x16, 0xd, 0x74, 0x65, 0x73, 0x74, 0x31, 0x40, 0x72, 0x73, 0x61, 0x2e, 0x63, 0x6f, 0x6d},
			value:       "test1@rsa.com",
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
			name:        "Test encode utc time",
			expected:    []byte{0x17, 0x11, 0x39, 0x31, 0x30, 0x35, 0x30, 0x36, 0x31, 0x36, 0x34, 0x35, 0x34, 0x30, 0x2D, 0x30, 0x37, 0x30, 0x30},
			value:       time.Date(1991, time.May, 6, 16, 45, 40, 0, time.FixedZone("UTC-8", -7*60*60)),
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
			name:        "Test encode oid",
			expected:    []byte{0x06, 0x06, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d},
			value:       ObjectIdentifier{Root: 1, Subidentifiers: []uint64{2, 840, 113549}},
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
		{
			name:        "Test encode octet string",
			expected:    []byte{0x04, 0x08, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
			value:       []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode octet string",
			expected:    []byte{0x04, 0x06, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x21},
			value:       []byte("Hello!"),
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x01, 0x00},
			value:       0,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x01, 0x7f},
			value:       127,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x02, 0x00, 0x80},
			value:       128,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x02, 0x01, 0x00},
			value:       256,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x01, 0x80},
			value:       -128,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x02, 0xff, 0x7f},
			value:       -129,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x03, 0x00, 0xc3, 0x50},
			value:       50000,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode integer",
			expected:    []byte{0x02, 0x02, 0xcf, 0xc7},
			value:       -12345,
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode enumerated",
			expected:    []byte{0x0a, 0x01, 00},
			value:       0,
			options:     "enumerated",
			errExpected: false,
		},
		{
			name:        "Test encode null",
			expected:    []byte{0x05, 0x00},
			value:       Null{},
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode null",
			expected:    []byte{0x42, 0x00},
			value:       Null{},
			options:     "tag:2,application",
			errExpected: false,
		},
		{
			name:        "Test encode struct",
			expected:    []byte{0x30, 0x0e, 0x04, 0x06, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x21, 0x01, 0x01, 0xff, 0x02, 0x01, 0x05},
			value:       SequenceTest{OctetStringValue: []byte("Hello!"), BooleanValue: true, IntegerValue: 5},
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode struct",
			expected:    []byte{0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x0b, 0x05, 0x00},
			value:       AnotherSequenceTest{Algorithm: ObjectIdentifier{Root: 1, Subidentifiers: []uint64{2, 840, 113549, 1, 1, 11}}},
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode slice",
			expected:    []byte{0x30, 0x09, 0x02, 0x01, 0x07, 0x02, 0x01, 0x08, 0x02, 0x01, 0x09},
			value:       []int{7, 8, 9},
			options:     "",
			errExpected: false,
		},
		{
			name:        "Test encode bitstring",
			expected:    []byte{0x03, 0x04, 0x06, 0x6e, 0x5d, 0xc0},
			value:       "b011011100101110111",
			options:     "bitstring",
			errExpected: false,
		},
		{
			name:        "Test encode octetstring",
			expected:    []byte{0x04, 0x04, 0x03, 0x02, 0x06, 0xA0},
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
