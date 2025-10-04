# Binary

Golang binary serialization library, support struct with tag

This library provides binary serialization and deserialization of Go values similar to `json.Marshal` and `json.Unmarshal` but using binary format. It's based on reflection and struct tags.

## Features

- Serialize/deserialize any Go value to/from binary format using Marshal/Unmarshal functions
- Support for fixed-length types without tags
- Support for variable-length types (string, []byte, slices) with optional tags
- Support for arrays ([N]T types) with optional tags
- Default format for variable-length types: `len(data) + data`
- Tag support for specifying fixed lengths:
  - If tag length is greater than data length, pad with zeros
  - If tag length is less than data length, truncate extra data
- Support for custom BinaryMarshaler and BinaryUnmarshaler interfaces

## Installation

```bash
go get github.com/lengzhao/binary
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"github.com/lengzhao/binary"
)

type Person struct {
	Name    string
	Age     uint8
	Email   string     `binary:"50"`  // Fixed length of 50 bytes
	Data    []byte     `binary:"10"`  // Fixed length of 10 bytes
	Scores  []uint32   `binary:"5"`   // Fixed length of 5 elements
	ID      [16]byte   `binary:"16"`  // Fixed length of 16 bytes
	Values  [4]uint32  `binary:"4"`   // Fixed length of 4 elements
}

func main() {
	person := Person{
		Name:   "Alice",
		Age:    30,
		Email:  "alice@example.com",
		Data:   []byte{1, 2, 3, 4, 5},      // Only 5 bytes, will be padded to 10
		Scores: []uint32{100, 95, 87},      // Only 3 elements, will be padded to 5
		ID:     [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Values: [4]uint32{1000, 2000, 3000, 4000},
	}

	// Serialize
	data, err := binary.Marshal(person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Serialized data length: %d bytes\n", len(data))

	// Deserialize
	var decoded Person
	err = binary.Unmarshal(data, &decoded)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Decoded: %+v\n", decoded)
}
```

### Direct Value Encoding

You can now encode and decode any supported Go value directly:

```go
// Direct slice encoding/decoding
slice := []uint32{10, 20, 30, 40, 50}
data, err := binary.Marshal(slice)
// ... handle error
var decodedSlice []uint32
err = binary.Unmarshal(data, &decodedSlice)

// Direct array encoding/decoding
array := [5]uint32{100, 200, 300, 400, 500}
data, err := binary.Marshal(array)
// ... handle error
var decodedArray [5]uint32
err = binary.Unmarshal(data, &decodedArray)

// Direct basic type encoding/decoding
number := uint32(42)
data, err := binary.Marshal(number)
// ... handle error
var decodedNumber uint32
err = binary.Unmarshal(data, &decodedNumber)
```

### Partial Unmarshaling

The library now supports partial unmarshaling with `UnmarshalPartial`, which allows you to decode data and get information about remaining bytes:

```go
type Message struct {
    ID   uint32
    Text string
}

// Create data containing multiple messages
data := []byte{...} // Binary data with multiple messages

// Process messages sequentially
currentData := data
for len(currentData) > 0 {
    var msg Message
    remaining, err := binary.UnmarshalPartial(currentData, &msg)
    if err != nil {
        break
    }
    
    fmt.Printf("Decoded message: %+v, remaining bytes: %d\n", msg, remaining)
    
    // Move to next message
    if remaining == 0 {
        break
    }
    currentData = currentData[len(currentData)-remaining:]
}
```

**Key differences between `Unmarshal` and `UnmarshalPartial`:**

- `Unmarshal(data []byte, v interface{}) error`:
  - Expects all data to be consumed
  - Returns error if there are remaining bytes after unmarshaling
  - Best for cases where you expect exact data match

- `UnmarshalPartial(data []byte, v interface{}) (remaining int, error)`:
  - Allows partial data consumption
  - Returns the number of bytes remaining after unmarshaling
  - Perfect for processing data streams or multiple consecutive structures
  - Useful when input data may contain more information than needed

### Custom Encoder/Decoder

Structs can implement the BinaryMarshaler and BinaryUnmarshaler interfaces for custom serialization:

```go
type CustomType struct {
	Value string
}

func (c *CustomType) MarshalBinary() ([]byte, error) {
	return []byte("custom:" + c.Value), nil
}

func (c *CustomType) UnmarshalBinary(data []byte) error {
	if len(data) < 7 || string(data[:7]) != "custom:" {
		return nil // Not in our custom format
	}
	c.Value = string(data[7:])
	return nil
}
```

### Tag Format

Tags can be specified in the following formats:

1. Simple length: `binary:"50"` - Fixed length of 50 bytes
2. Length specifier: `binary:"len:50"` - Fixed length of 50 bytes
3. Ignore tag: `binary:"-"` - Ignore the field

For variable-length types without tags, the library uses the default format: `len(data) + data`

For fixed-length types with tags:
- If data is shorter than specified length, pad with zeros (or zero values for slices/arrays)
- If data is longer than specified length, truncate extra data

This applies to:
- String types: Pad with zeros or truncate
- []byte types: Pad with zeros or truncate
- Slice types: Pad with zero values or truncate
- Array types: Pad with zero values or truncate

### Supported Types

- Integer types: `uint8`, `uint16`, `uint32`, `uint64`, `int8`, `int16`, `int32`, `int64`
- Floating point types: `float32`, `float64`
- String
- Byte slice (`[]byte`)
- Byte arrays (`[N]byte`)
- Other slices
- Other arrays
- Structs
- Nested structs

## Implementation Details

- Uses little-endian encoding for numeric types
- For fixed-length types with tags:
  - If data is shorter than specified length, pad with zeros (or zero values for slices/arrays)
  - If data is longer than specified length, truncate extra data
  - No length prefix is written when using tags
- For variable-length types without tags, uses `len(data) + data` format where len is a `uint32`
- Slices with tags are serialized as `elements` (no length prefix)
- Slices without tags are serialized as `len(slice) + elements` where len is a `uint32`
- Arrays with tags are serialized as `elements` (no length prefix)
- Arrays without tags are serialized as `len(array) + elements` where len is a `uint32`
- If a struct implements BinaryMarshaler/BinaryUnmarshaler, those methods are used instead of the default reflection-based approach
- Direct value encoding is now supported for all supported types