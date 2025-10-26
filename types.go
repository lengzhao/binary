// Package binary provides binary encoding and decoding functionality
// similar to Go's standard encoding packages.
//
// Main APIs:
//   - Marshal(v interface{}) ([]byte, error): Serialize any Go value to binary data
//   - Unmarshal(data []byte, v interface{}) error: Deserialize binary data to Go value
//   - UnmarshalPartial(data []byte, v interface{}) (remaining int, error): Partial deserialization with remaining byte count
//
// The UnmarshalPartial function allows for partial parsing of data streams,
// returning the number of bytes that remain unprocessed. This is useful for:
//   - Processing multiple consecutive structures from a single data buffer
//   - Handling data streams where you need to process data incrementally
//   - Scenarios where the input data may contain more information than needed
//
// Supported data types:
//   - Integer types: uint8, uint16, uint32, uint64, int8, int16, int32, int64
//   - Boolean type: bool
//   - Floating point types: float32, float64
//   - String
//   - Byte slice ([]byte)
//   - Byte arrays ([N]byte)
//   - Other slices
//   - Other arrays
//   - Structs
//   - Nested structs
//
// Custom types can implement BinaryMarshaler and BinaryUnmarshaler interfaces
// for custom serialization behavior.
package binary

// BinaryMarshaler is the interface implemented by types that can marshal themselves into binary form.
type BinaryMarshaler interface {
	MarshalBinary() ([]byte, error)
}

// BinaryUnmarshaler is the interface implemented by types that can unmarshal themselves from binary form.
type BinaryUnmarshaler interface {
	UnmarshalBinary([]byte) error
}
