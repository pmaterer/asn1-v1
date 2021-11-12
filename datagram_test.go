package asn1

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createBytes(size int) []byte {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func TestEncodeLength(t *testing.T) {
	tests := []struct {
		name     string
		expected []byte
		value    []byte
	}{
		{
			name:     "Test length short",
			expected: []byte{0x19},
			value:    createBytes(0x19),
		},
		{
			name:     "Test length long",
			expected: []byte{0x82, 0x01, 0xb3},
			value:    createBytes(0x1b3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			assert.Equal(t, tt.expected, encodeLength(uint(len(tt.value))))

		})
	}
}
