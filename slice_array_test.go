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

// 异常场景测试用例

func TestEncodeNonSupportedType(t *testing.T) {
	// Test encoding a channel (not supported)
	ch := make(chan int)
	_, err := Encode(ch)
	assert.Error(t, err)
}

func TestDecodeToNonPointer(t *testing.T) {
	// Test decoding to a non-pointer value
	data := []byte{1, 2, 3, 4}
	var decoded []uint32
	err := Decode(data, decoded) // 注意：这里传入的是值而不是指针
	assert.Error(t, err)
}

func TestDecodeToUnsupportedType(t *testing.T) {
	// Test decoding to an unsupported type
	data := []byte{1, 2, 3, 4}
	var decoded chan int
	err := Decode(data, &decoded)
	assert.Error(t, err)
}

func TestDecodeWithInsufficientData(t *testing.T) {
	// Test decoding with insufficient data
	data := []byte{1, 2} // 不足以解码一个uint32
	var decoded []uint32
	err := Decode(data, &decoded)
	assert.Error(t, err)
}

func TestEncodeNilSlice(t *testing.T) {
	// Test encoding a nil slice
	var original []uint32 = nil

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded []uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	// Nil slice should decode to empty slice
	assert.Empty(t, decoded)
}

func TestEncodeEmptySlice(t *testing.T) {
	// Test encoding an empty slice
	original := []uint32{}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded []uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestEncodeEmptyArray(t *testing.T) {
	// Test encoding an empty array
	original := [0]uint32{}

	data, err := Encode(original)
	assert.NoError(t, err)

	var decoded [0]uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original, decoded)
}

func TestDecodeToNilPointer(t *testing.T) {
	// Test decoding to a nil pointer
	data := []byte{1, 2, 3, 4}
	var decoded []uint32 = nil
	err := Decode(data, &decoded)
	// Should return an error because we can't decode into a nil pointer
	assert.Error(t, err)
}

// 更多异常场景测试用例

func TestEncodeUnsupportedChannelType(t *testing.T) {
	// Test encoding a channel (not supported)
	ch := make(chan int)
	_, err := Encode(ch)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestEncodeUnsupportedFuncType(t *testing.T) {
	// Test encoding a function (not supported)
	fn := func() {}
	_, err := Encode(fn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestEncodeUnsupportedMapType(t *testing.T) {
	// Test encoding a map (not supported)
	m := make(map[string]int)
	_, err := Encode(m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestEncodeUnsupportedPointerType(t *testing.T) {
	// Test encoding a pointer to unsupported type
	// But pointer to channel should fail
	ch := make(chan int)
	_, err := Encode(&ch)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")

	// Pointer to function should fail
	fn := func() {}
	_, err = Encode(&fn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")

	// Pointer to map should fail
	m := make(map[string]int)
	_, err = Encode(&m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestDecodeToUnsupportedChannelType(t *testing.T) {
	// Test decoding to a channel (not supported)
	data := []byte{1, 2, 3, 4}
	var ch chan int
	err := Decode(data, &ch)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestDecodeToUnsupportedFuncType(t *testing.T) {
	// Test decoding to a function (not supported)
	data := []byte{1, 2, 3, 4}
	var fn func()
	err := Decode(data, &fn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestDecodeToUnsupportedMapType(t *testing.T) {
	// Test decoding to a map (not supported)
	data := []byte{1, 2, 3, 4}
	var m map[string]int
	err := Decode(data, &m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestDecodeWithMalformedData(t *testing.T) {
	// Test decoding with malformed data
	data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF} // Malformed data
	var decoded []uint32
	err := Decode(data, &decoded)
	assert.Error(t, err)
}

func TestEncodeLargeSlice(t *testing.T) {
	// Test encoding a large slice
	original := make([]uint32, 10000)
	for i := range original {
		original[i] = uint32(i)
	}

	data, err := Encode(original)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded []uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}

func TestEncodeNestedSlice(t *testing.T) {
	// Test encoding a nested slice structure
	original := [][]uint32{{1, 2}, {3, 4, 5}, {6}}

	data, err := Encode(original)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	var decoded [][]uint32
	err = Decode(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}
