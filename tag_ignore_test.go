package binary

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSliceIgnoreTag tests that slice with tag "-" is properly ignored
func TestSliceIgnoreTag(t *testing.T) {
	type TestStruct struct {
		Data   []uint32 `binary:"-"`
		Number uint32
		Name   string
	}

	// Test that the "-" tag causes the field to be skipped entirely
	original := TestStruct{
		Data:   []uint32{100, 200, 300}, // This field should be ignored
		Number: 42,
		Name:   "test",
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	// Check that data only contains Number and Name (not Data)
	// Number (4 bytes) + Name length (4 bytes) + Name data (4 bytes) = 12 bytes
	assert.Equal(t, 12, len(data))

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be empty since it was skipped
	assert.Equal(t, []uint32(nil), decoded.Data)
	// Other fields should be preserved
	assert.Equal(t, original.Number, decoded.Number)
	assert.Equal(t, original.Name, decoded.Name)
}

// TestArrayIgnoreTag tests that array with tag "-" is properly ignored
func TestArrayIgnoreTag(t *testing.T) {
	type TestStruct struct {
		Data   [3]uint32 `binary:"-"`
		Number uint32
		Name   string
	}

	// Test that the "-" tag causes the field to be skipped entirely
	original := TestStruct{
		Data:   [3]uint32{100, 200, 300}, // This field should be ignored
		Number: 42,
		Name:   "test",
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	// Check that data only contains Number and Name (not Data)
	// Number (4 bytes) + Name length (4 bytes) + Name data (4 bytes) = 12 bytes
	assert.Equal(t, 12, len(data))

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be zero values since it was skipped
	assert.Equal(t, [3]uint32{0, 0, 0}, decoded.Data)
	// Other fields should be preserved
	assert.Equal(t, original.Number, decoded.Number)
	assert.Equal(t, original.Name, decoded.Name)
}

// TestStringIgnoreTag tests that string with tag "-" is properly ignored
func TestStringIgnoreTag(t *testing.T) {
	type TestStruct struct {
		Data   string `binary:"-"`
		Number uint32
		Name   string
	}

	// Test that the "-" tag causes the field to be skipped entirely
	original := TestStruct{
		Data:   "should be ignored", // This field should be ignored
		Number: 42,
		Name:   "test",
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	// Check that data only contains Number and Name (not Data)
	// Number (4 bytes) + Name length (4 bytes) + Name data (4 bytes) = 12 bytes
	assert.Equal(t, 12, len(data))

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be empty since it was skipped
	assert.Equal(t, "", decoded.Data)
	// Other fields should be preserved
	assert.Equal(t, original.Number, decoded.Number)
	assert.Equal(t, original.Name, decoded.Name)
}

// TestBytesIgnoreTag tests that []byte with tag "-" is properly ignored
func TestBytesIgnoreTag(t *testing.T) {
	type TestStruct struct {
		Data   []byte `binary:"-"`
		Number uint32
		Name   string
	}

	// Test that the "-" tag causes the field to be skipped entirely
	original := TestStruct{
		Data:   []byte{1, 2, 3, 4, 5}, // This field should be ignored
		Number: 42,
		Name:   "test",
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	// Check that data only contains Number and Name (not Data)
	// Number (4 bytes) + Name length (4 bytes) + Name data (4 bytes) = 12 bytes
	assert.Equal(t, 12, len(data))

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be empty since it was skipped
	assert.Equal(t, []byte(nil), decoded.Data)
	// Other fields should be preserved
	assert.Equal(t, original.Number, decoded.Number)
	assert.Equal(t, original.Name, decoded.Name)
}

// TestByteArrayIgnoreTag tests that [N]byte with tag "-" is properly ignored
func TestByteArrayIgnoreTag(t *testing.T) {
	type TestStruct struct {
		Data   [5]byte `binary:"-"`
		Number uint32
		Name   string
	}

	// Test that the "-" tag causes the field to be skipped entirely
	original := TestStruct{
		Data:   [5]byte{1, 2, 3, 4, 5}, // This field should be ignored
		Number: 42,
		Name:   "test",
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	// Check that data only contains Number and Name (not Data)
	// Number (4 bytes) + Name length (4 bytes) + Name data (4 bytes) = 12 bytes
	assert.Equal(t, 12, len(data))

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be zero values since it was skipped
	assert.Equal(t, [5]byte{0, 0, 0, 0, 0}, decoded.Data)
	// Other fields should be preserved
	assert.Equal(t, original.Number, decoded.Number)
	assert.Equal(t, original.Name, decoded.Name)
}

// TestDirectSliceIgnore tests ignoring tag when directly encoding/decoding slices
func TestDirectSliceIgnore(t *testing.T) {
	original := []int32{1, 2, 3, 4, 5}

	// Test that we can still encode/decode normally (no tag)
	data, err := Marshal(original)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(data)) // Should have actual data

	var decoded []int32
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(original, decoded))
}

// TestDirectArrayIgnore tests ignoring tag when directly encoding/decoding arrays
func TestDirectArrayIgnore(t *testing.T) {
	original := [5]int32{1, 2, 3, 4, 5}

	// Test that we can still encode/decode normally (no tag)
	data, err := Marshal(original)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(data)) // Should have actual data

	var decoded [5]int32
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(original, decoded))
}
