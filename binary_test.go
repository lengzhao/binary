package binary

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeSimpleStruct(t *testing.T) {
	type SimpleStruct struct {
		A uint32
		B int16
		C uint8
	}

	original := SimpleStruct{
		A: 12345,
		B: -100,
		C: 255,
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded SimpleStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeString(t *testing.T) {
	type StringStruct struct {
		Name string
	}

	original := StringStruct{
		Name: "Hello, 世界",
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded StringStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeStringWithTag(t *testing.T) {
	type StringWithTagStruct struct {
		Name string `binary:"20"`
	}

	original := StringWithTagStruct{
		Name: "Hello",
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded StringWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// For fixed length strings, we expect exactly 20 bytes
	expected := "Hello"
	if len(decoded.Name) > 20 {
		t.Errorf("Expected name to be at most 20 characters, got %d", len(decoded.Name))
	}

	// The decoded name should start with "Hello"
	if len(decoded.Name) < len(expected) || decoded.Name[:len(expected)] != expected {
		t.Errorf("Expected name to start with %s, got %s", expected, decoded.Name)
	}
}

func TestEncodeDecodeBytes(t *testing.T) {
	type BytesStruct struct {
		Data []byte
	}

	original := BytesStruct{
		Data: []byte{1, 2, 3, 4, 5},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded BytesStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeBytesWithTagTruncate(t *testing.T) {
	type BytesWithTagStruct struct {
		Data []byte `binary:"3"`
	}

	original := BytesWithTagStruct{
		Data: []byte{1, 2, 3, 4, 5}, // 5 bytes, but tag specifies 3
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Check that encoded data has exactly 3 bytes
	if len(data) != 3 {
		t.Errorf("Expected encoded data to be 3 bytes, got %d", len(data))
	}

	// Check that encoded data contains first 3 bytes
	expected := []byte{1, 2, 3}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected encoded data to be %v, got %v", expected, data)
	}

	var decoded BytesWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded data has exactly 3 bytes
	if len(decoded.Data) != 3 {
		t.Errorf("Expected decoded data to be 3 bytes, got %d", len(decoded.Data))
	}

	// Check that decoded data contains first 3 bytes
	if !reflect.DeepEqual(decoded.Data, expected) {
		t.Errorf("Expected decoded data to be %v, got %v", expected, decoded.Data)
	}
}

func TestEncodeDecodeBytesWithTagPad(t *testing.T) {
	type BytesWithTagStruct struct {
		Data []byte `binary:"7"`
	}

	original := BytesWithTagStruct{
		Data: []byte{1, 2, 3}, // 3 bytes, but tag specifies 7
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Check that encoded data has exactly 7 bytes
	if len(data) != 7 {
		t.Errorf("Expected encoded data to be 7 bytes, got %d", len(data))
	}

	// Check that encoded data contains first 3 bytes and 4 zero bytes
	expected := []byte{1, 2, 3, 0, 0, 0, 0}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("Expected encoded data to be %v, got %v", expected, data)
	}

	var decoded BytesWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded data has exactly 7 bytes
	if len(decoded.Data) != 7 {
		t.Errorf("Expected decoded data to be 7 bytes, got %d", len(decoded.Data))
	}

	// Trim trailing zeros for comparison
	trimmed := decoded.Data
	for len(trimmed) > 0 && trimmed[len(trimmed)-1] == 0 {
		trimmed = trimmed[:len(trimmed)-1]
	}

	originalTrimmed := original.Data
	if !reflect.DeepEqual(trimmed, originalTrimmed) {
		t.Errorf("Expected decoded data to be %v, got %v", originalTrimmed, trimmed)
	}
}

func TestEncodeDecodeSlice(t *testing.T) {
	type SliceStruct struct {
		Numbers []uint32
	}

	original := SliceStruct{
		Numbers: []uint32{10, 20, 30, 40, 50},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded SliceStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeSliceWithTagTruncate(t *testing.T) {
	type SliceWithTagStruct struct {
		Numbers []uint32 `binary:"3"`
	}

	original := SliceWithTagStruct{
		Numbers: []uint32{10, 20, 30, 40, 50}, // 5 elements, but tag specifies 3
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the data
	var decoded SliceWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded slice has exactly 3 elements
	if len(decoded.Numbers) != 3 {
		t.Errorf("Expected decoded slice to have 3 elements, got %d", len(decoded.Numbers))
	}

	// Check that decoded slice contains first 3 elements
	expected := []uint32{10, 20, 30}
	if !reflect.DeepEqual(decoded.Numbers, expected) {
		t.Errorf("Expected decoded slice to be %v, got %v", expected, decoded.Numbers)
	}
}

func TestDecodeSliceWithTag(t *testing.T) {
	type TestStruct struct {
		Data []uint32 `binary:"5"`
	}

	// Test case 1: Data shorter than fixed length (should pad with zeros)
	original1 := TestStruct{
		Data: []uint32{100, 200, 300}, // 3 elements, should pad to 5
	}

	data1, err := Marshal(original1)
	assert.NoError(t, err)

	var decoded1 TestStruct
	err = Unmarshal(data1, &decoded1)
	assert.NoError(t, err)

	// Should have 5 elements: [100, 200, 300, 0, 0]
	expected1 := []uint32{100, 200, 300, 0, 0}
	assert.Equal(t, expected1, decoded1.Data)

	// Test case 2: Data longer than fixed length (should truncate)
	original2 := TestStruct{
		Data: []uint32{100, 200, 300, 400, 500, 600, 700}, // 7 elements, should truncate to 5
	}

	data2, err := Marshal(original2)
	assert.NoError(t, err)

	var decoded2 TestStruct
	err = Unmarshal(data2, &decoded2)
	assert.NoError(t, err)

	// Should have 5 elements: [100, 200, 300, 400, 500]
	expected2 := []uint32{100, 200, 300, 400, 500}
	assert.Equal(t, expected2, decoded2.Data)

	// Test case 3: Data exactly matching fixed length
	original3 := TestStruct{
		Data: []uint32{100, 200, 300, 400, 500}, // 5 elements, should remain 5
	}

	data3, err := Marshal(original3)
	assert.NoError(t, err)

	var decoded3 TestStruct
	err = Unmarshal(data3, &decoded3)
	assert.NoError(t, err)

	// Should have 5 elements: [100, 200, 300, 400, 500]
	expected3 := []uint32{100, 200, 300, 400, 500}
	assert.Equal(t, expected3, decoded3.Data)
}

func TestEncodeDecodeSliceWithTagPad(t *testing.T) {
	type SliceWithTagStruct struct {
		Numbers []uint32 `binary:"7"`
	}

	original := SliceWithTagStruct{
		Numbers: []uint32{10, 20, 30}, // 3 elements, but tag specifies 7
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the data
	var decoded SliceWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded slice has exactly 7 elements
	if len(decoded.Numbers) != 7 {
		t.Errorf("Expected decoded slice to have 7 elements, got %d", len(decoded.Numbers))
	}

	// Check that decoded slice contains first 3 elements and 4 zero elements
	expected := []uint32{10, 20, 30, 0, 0, 0, 0}
	if !reflect.DeepEqual(decoded.Numbers, expected) {
		t.Errorf("Expected decoded slice to be %v, got %v", expected, decoded.Numbers)
	}
}

func TestEncodeDecodeArray(t *testing.T) {
	type ArrayStruct struct {
		Numbers [5]uint32
	}

	original := ArrayStruct{
		Numbers: [5]uint32{10, 20, 30, 40, 50},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded ArrayStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeArrayWithTagTruncate(t *testing.T) {
	type ArrayWithTagStruct struct {
		Numbers [5]uint32 `binary:"3"`
	}

	original := ArrayWithTagStruct{
		Numbers: [5]uint32{10, 20, 30, 40, 50}, // 5 elements, but tag specifies 3
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the data
	var decoded ArrayWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded array has exactly 5 elements
	if len(decoded.Numbers) != 5 {
		t.Errorf("Expected decoded array to have 5 elements, got %d", len(decoded.Numbers))
	}

	// Check that first 3 elements are as expected and remaining are zero
	expected := [5]uint32{10, 20, 30, 0, 0}
	if !reflect.DeepEqual(decoded.Numbers, expected) {
		t.Errorf("Expected decoded array to be %v, got %v", expected, decoded.Numbers)
	}
}

func TestEncodeDecodeArrayWithTagPad(t *testing.T) {
	type ArrayWithTagStruct struct {
		Numbers [3]uint32 `binary:"5"`
	}

	original := ArrayWithTagStruct{
		Numbers: [3]uint32{10, 20, 30}, // 3 elements, but tag specifies 5
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the data
	var decoded ArrayWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded array has exactly 3 elements
	if len(decoded.Numbers) != 3 {
		t.Errorf("Expected decoded array to have 3 elements, got %d", len(decoded.Numbers))
	}

	// Check that elements are as expected
	expected := [3]uint32{10, 20, 30}
	if !reflect.DeepEqual(decoded.Numbers, expected) {
		t.Errorf("Expected decoded array to be %v, got %v", expected, decoded.Numbers)
	}
}

func TestEncodeDecodeByteArray(t *testing.T) {
	type ByteArrayStruct struct {
		Data [5]byte
	}

	original := ByteArrayStruct{
		Data: [5]byte{1, 2, 3, 4, 5},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded ByteArrayStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeByteArrayWithTagTruncate(t *testing.T) {
	type ByteArrayWithTagStruct struct {
		Data [5]byte `binary:"3"`
	}

	original := ByteArrayWithTagStruct{
		Data: [5]byte{1, 2, 3, 4, 5}, // 5 bytes, but tag specifies 3
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the data
	var decoded ByteArrayWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded array has exactly 5 elements
	if len(decoded.Data) != 5 {
		t.Errorf("Expected decoded array to have 5 elements, got %d", len(decoded.Data))
	}

	// Check that first 3 elements are as expected and remaining are zero
	expected := [5]byte{1, 2, 3, 0, 0}
	if !reflect.DeepEqual(decoded.Data, expected) {
		t.Errorf("Expected decoded array to be %v, got %v", expected, decoded.Data)
	}
}

func TestEncodeDecodeByteArrayWithTagPad(t *testing.T) {
	type ByteArrayWithTagStruct struct {
		Data [3]byte `binary:"5"`
	}

	original := ByteArrayWithTagStruct{
		Data: [3]byte{1, 2, 3}, // 3 bytes, but tag specifies 5
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode the data
	var decoded ByteArrayWithTagStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Check that decoded array has exactly 3 elements
	if len(decoded.Data) != 3 {
		t.Errorf("Expected decoded array to have 3 elements, got %d", len(decoded.Data))
	}

	// Check that elements are as expected
	expected := [3]byte{1, 2, 3}
	if !reflect.DeepEqual(decoded.Data, expected) {
		t.Errorf("Expected decoded array to be %v, got %v", expected, decoded.Data)
	}
}

func TestEncodeDecodeNestedStruct(t *testing.T) {
	type Address struct {
		Street string
		Number uint16
	}

	type Person struct {
		Name    string
		Age     uint8
		Address Address
	}

	original := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			Street: "Main St",
			Number: 123,
		},
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded Person
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

func TestEncodeDecodeFloats(t *testing.T) {
	type FloatStruct struct {
		Float32Value float32
		Float64Value float64
	}

	original := FloatStruct{
		Float32Value: 3.14159,
		Float64Value: 2.718281828459045,
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	var decoded FloatStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Use InEpsilon for floating point comparison
	assert.InEpsilon(t, original.Float32Value, decoded.Float32Value, 1e-6)
	assert.InEpsilon(t, original.Float64Value, decoded.Float64Value, 1e-12)
}

func TestParseTag(t *testing.T) {
	tests := []struct {
		tag      string
		expected uint32
		hasError bool
	}{
		{"", 0, true},
		{"-", 0, true},
		{"10", 10, false},
		{"len:20", 20, false},
		{"len:abc", 0, true},
		{"invalid", 0, true},
	}

	for _, test := range tests {
		result, err := parseTag(test.tag)
		if test.hasError {
			assert.Error(t, err, "Expected error for tag: %s", test.tag)
		} else {
			assert.NoError(t, err, "Unexpected error for tag: %s", test.tag)
			assert.Equal(t, test.expected, result, "Expected %d for tag: %s", test.expected, test.tag)
		}
	}
}

func TestIgnoreTag(t *testing.T) {
	type TestStruct struct {
		Data []uint32 `binary:"-"`
	}

	// Test that the "-" tag causes the field to be skipped entirely
	original := TestStruct{
		Data: []uint32{100, 200, 300}, // Should be ignored completely
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be empty since it was skipped
	assert.Equal(t, []uint32(nil), decoded.Data)
}

func TestIgnoreTagSkipField(t *testing.T) {
	type TestStruct struct {
		Data   []uint32 `binary:"-"`
		Number uint32
		Name   string
	}

	// Test that the "-" tag causes the field to be skipped
	original := TestStruct{
		Data:   []uint32{100, 200, 300}, // This field should be ignored
		Number: 42,
		Name:   "test",
	}

	data, err := Marshal(original)
	assert.NoError(t, err)

	var decoded TestStruct
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Data field should be empty since it was skipped
	assert.Equal(t, []uint32(nil), decoded.Data)
	// Other fields should be preserved
	assert.Equal(t, original.Number, decoded.Number)
	assert.Equal(t, original.Name, decoded.Name)
}

// Test custom BinaryMarshaler and BinaryUnmarshaler implementation
type CustomType struct {
	Value string
}

func (c CustomType) MarshalBinary() ([]byte, error) {
	return []byte("custom:" + c.Value), nil
}

func (c *CustomType) UnmarshalBinary(data []byte) error {
	if len(data) < 7 || string(data[:7]) != "custom:" {
		return nil // Not in our custom format
	}
	c.Value = string(data[7:])
	return nil
}

func TestCustomMarshalerUnmarshaler(t *testing.T) {
	type StructWithCustomType struct {
		Custom CustomType
		Number uint32
	}

	original := StructWithCustomType{
		Custom: CustomType{Value: "test"},
		Number: 42,
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	var decoded StructWithCustomType
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if original.Custom.Value != decoded.Custom.Value {
		t.Errorf("Expected custom value %s, got %s", original.Custom.Value, decoded.Custom.Value)
	}

	if original.Number != decoded.Number {
		t.Errorf("Expected number %d, got %d", original.Number, decoded.Number)
	}
}

// Test the Marshal/Unmarshal API
func TestMarshalUnmarshal(t *testing.T) {
	type SimpleStruct struct {
		A uint32
		B int16
		C uint8
	}

	original := SimpleStruct{
		A: 12345,
		B: -100,
		C: 255,
	}

	data, err := Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded SimpleStruct
	err = Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !reflect.DeepEqual(original, decoded) {
		t.Errorf("Expected %+v, got %+v", original, decoded)
	}
}

// TestRemainingDataHandling tests how the decoder handles remaining data
func TestRemainingDataHandling(t *testing.T) {
	// Test struct decoding with extra data
	type TestStruct struct {
		A int32
		B string
	}

	// Create test data: normal struct data + extra garbage data
	buf := new(bytes.Buffer)

	// Write A field (int32)
	binary.Write(buf, binary.LittleEndian, int32(42))

	// Write B field (string): length + data
	testString := "hello"
	binary.Write(buf, binary.LittleEndian, uint32(len(testString)))
	buf.WriteString(testString)

	// Add extra garbage data
	extraData := []byte{0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA}
	buf.Write(extraData)

	data := buf.Bytes()
	t.Logf("Total data length: %d bytes", len(data))
	t.Logf("Data content: %x", data)

	// Try to decode
	var result TestStruct
	err := Unmarshal(data, &result)
	assert.Error(t, err, "Expected error due to remaining data")
	assert.Contains(t, err.Error(), "bytes of data remaining", "Error should mention remaining data")
	t.Logf("Decode error (expected): %v", err)

	// Test array decoding with extra elements being skipped
	type ArrayStruct struct {
		Numbers [3]int32 `binary:"5"` // Fixed read 5 elements, but array can only hold 3
	}

	t.Log("Testing array remaining data handling:")
	buf2 := new(bytes.Buffer)
	// Write 5 int32 data
	for i := int32(1); i <= 5; i++ {
		binary.Write(buf2, binary.LittleEndian, i)
	}

	data2 := buf2.Bytes()
	t.Logf("Array test data length: %d bytes", len(data2))

	var arrayResult ArrayStruct
	err = Unmarshal(data2, &arrayResult)
	assert.NoError(t, err, "Array decoding should succeed (extra elements are skipped)")
	assert.Equal(t, [3]int32{1, 2, 3}, arrayResult.Numbers, "Array should contain first 3 elements")
	t.Logf("Array decode result: %+v", arrayResult)
}

// TestUnmarshalPartial tests the partial unmarshaling functionality
func TestUnmarshalPartial(t *testing.T) {
	// Test basic partial unmarshaling
	type SimpleStruct struct {
		A int32
		B string
	}

	// Create test data: struct + extra data
	buf := new(bytes.Buffer)

	// Write struct data
	binary.Write(buf, binary.LittleEndian, int32(42))
	testString := "hello"
	binary.Write(buf, binary.LittleEndian, uint32(len(testString)))
	buf.WriteString(testString)

	// Add extra data
	extraData := []byte{0xFF, 0xEE, 0xDD, 0xCC}
	buf.Write(extraData)

	data := buf.Bytes()
	t.Logf("Total test data: %d bytes", len(data))

	// Test partial unmarshaling
	var result SimpleStruct
	remaining, err := UnmarshalPartial(data, &result)
	assert.NoError(t, err, "UnmarshalPartial should succeed")
	assert.Equal(t, len(extraData), remaining, "Should have %d bytes remaining", len(extraData))
	assert.Equal(t, int32(42), result.A, "Field A should be decoded correctly")
	assert.Equal(t, "hello", result.B, "Field B should be decoded correctly")
	t.Logf("Decoded struct: %+v, remaining bytes: %d", result, remaining)

	// Test multiple sequential parsing
	type TwoInts struct {
		X int32
		Y int32
	}

	// Create data for two structs
	buf2 := new(bytes.Buffer)
	binary.Write(buf2, binary.LittleEndian, int32(100))
	binary.Write(buf2, binary.LittleEndian, int32(200))
	binary.Write(buf2, binary.LittleEndian, int32(300))
	binary.Write(buf2, binary.LittleEndian, int32(400))

	data2 := buf2.Bytes()
	t.Logf("Sequential test data: %d bytes", len(data2))

	// Parse first struct
	var first TwoInts
	remaining, err = UnmarshalPartial(data2, &first)
	assert.NoError(t, err, "First UnmarshalPartial should succeed")
	assert.Equal(t, 8, remaining, "Should have 8 bytes remaining after first struct")
	assert.Equal(t, int32(100), first.X, "First struct X should be 100")
	assert.Equal(t, int32(200), first.Y, "First struct Y should be 200")

	// Parse second struct from remaining data
	remainingData := data2[len(data2)-remaining:]
	var second TwoInts
	remaining2, err := UnmarshalPartial(remainingData, &second)
	assert.NoError(t, err, "Second UnmarshalPartial should succeed")
	assert.Equal(t, 0, remaining2, "Should have no bytes remaining after second struct")
	assert.Equal(t, int32(300), second.X, "Second struct X should be 300")
	assert.Equal(t, int32(400), second.Y, "Second struct Y should be 400")

	t.Logf("First struct: %+v, remaining: %d", first, remaining)
	t.Logf("Second struct: %+v, remaining: %d", second, remaining2)

	// Test compatibility with regular Unmarshal
	type OnlyOneInt struct {
		Value int32
	}

	// Create data with extra bytes
	buf3 := new(bytes.Buffer)
	binary.Write(buf3, binary.LittleEndian, int32(999))
	buf3.Write([]byte{0x01, 0x02, 0x03}) // extra bytes

	data3 := buf3.Bytes()

	// Regular Unmarshal should fail due to extra bytes
	var regularResult OnlyOneInt
	err = Unmarshal(data3, &regularResult)
	assert.Error(t, err, "Regular Unmarshal should fail with extra bytes")

	// UnmarshalPartial should succeed and report remaining bytes
	var partialResult OnlyOneInt
	remaining3, err := UnmarshalPartial(data3, &partialResult)
	assert.NoError(t, err, "UnmarshalPartial should succeed")
	assert.Equal(t, 3, remaining3, "Should have 3 bytes remaining")
	assert.Equal(t, int32(999), partialResult.Value, "Value should be decoded correctly")

	t.Logf("Compatibility test - partial result: %+v, remaining: %d", partialResult, remaining3)

	// Test with BinaryUnmarshaler interface
	var customType CustomType
	customData := []byte("custom:test123")
	remaining4, err := UnmarshalPartial(customData, &customType)
	assert.NoError(t, err, "UnmarshalPartial with BinaryUnmarshaler should succeed")
	assert.Equal(t, 0, remaining4, "BinaryUnmarshaler should consume all data")
	assert.Equal(t, "test123", customType.Value, "Custom type should be decoded correctly")

	t.Logf("BinaryUnmarshaler test - result: %+v, remaining: %d", customType, remaining4)
}
