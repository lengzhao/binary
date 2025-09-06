package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeSliceDirectly(t *testing.T) {
	// Test slice of integers
	original := []uint32{10, 20, 30, 40, 50}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded []uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEncodeDecodeArrayDirectly(t *testing.T) {
	// Test array of integers
	original := [5]uint32{10, 20, 30, 40, 50}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded [5]uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEncodeDecodeByteArrayDirectly(t *testing.T) {
	// Test byte array
	original := [5]byte{1, 2, 3, 4, 5}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded [5]byte
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEncodeDecodeByteSliceDirectly(t *testing.T) {
	// Test byte slice
	original := []byte{1, 2, 3, 4, 5}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded []byte
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEncodeDecodeStringSliceDirectly(t *testing.T) {
	// Test slice of strings
	original := []string{"hello", "world", "test"}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded []string
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}
