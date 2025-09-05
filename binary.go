// Package binary provides binary serialization and deserialization of Go structs
// similar to json.Marshal and json.Unmarshal but using binary format.
package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BinaryEncoder is the interface implemented by types that can encode themselves.
type BinaryEncoder interface {
	Encode() ([]byte, error)
}

// BinaryDecoder is the interface implemented by types that can decode themselves.
type BinaryDecoder interface {
	Decode([]byte) error
}

// Encode serializes a struct into binary format
func Encode(v interface{}) ([]byte, error) {
	// Check if the value implements BinaryEncoder
	if encoder, ok := v.(BinaryEncoder); ok {
		return encoder.Encode()
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("only structs are supported")
	}

	var buf bytes.Buffer
	err := encodeStruct(val, &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decode deserializes binary data into a struct
func Decode(data []byte, v interface{}) error {
	// Check if the value implements BinaryDecoder
	if decoder, ok := v.(BinaryDecoder); ok {
		return decoder.Decode(data)
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("only pointers to structs are supported")
	}

	buf := bytes.NewReader(data)
	return decodeStruct(buf, val.Elem())
}

// encodeStruct handles serialization of a struct
func encodeStruct(val reflect.Value, buf *bytes.Buffer) error {
	typ := val.Type()
	numField := val.NumField()

	for i := 0; i < numField; i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Check if field implements BinaryEncoder
		if encoder, ok := field.Interface().(BinaryEncoder); ok {
			fieldData, err := encoder.Encode()
			if err != nil {
				return fmt.Errorf("error encoding field %s: %w", fieldType.Name, err)
			}
			// Write length + data for the field
			length := uint32(len(fieldData))
			if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
				return err
			}
			_, err = buf.Write(fieldData)
			if err != nil {
				return err
			}
			continue
		}

		tag := fieldType.Tag.Get("binary")
		// If tag is "-", skip this field entirely
		if tag == "-" {
			continue
		}

		if err := encodeField(field, buf, tag); err != nil {
			return fmt.Errorf("error encoding field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// encodeField handles serialization of a single field
func encodeField(field reflect.Value, buf *bytes.Buffer, tag string) error {
	switch field.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		return binary.Write(buf, binary.LittleEndian, field.Interface())

	case reflect.String:
		return encodeString(field.String(), buf, tag)

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return encodeBytes(field.Bytes(), buf, tag)
		}
		// Other slices
		return encodeSlice(field, buf, tag)

	case reflect.Array:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			// [N]byte - convert to []byte
			length := field.Len()
			data := make([]byte, length)
			for i := 0; i < length; i++ {
				data[i] = byte(field.Index(i).Uint())
			}
			return encodeBytes(data, buf, tag)
		}
		// Other arrays
		return encodeArray(field, buf, tag)

	case reflect.Struct:
		return encodeStruct(field, buf)

	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
}

// encodeString handles serialization of strings
func encodeString(s string, buf *bytes.Buffer, tag string) error {
	data := []byte(s)

	// Check if tag specifies length
	if tag != "" {
		if length, err := parseTag(tag); err == nil {
			if uint32(len(data)) > length {
				// Truncate data if it's longer than specified length
				data = data[:length]
			} else if uint32(len(data)) < length {
				// Pad with zeros if data is shorter than specified length
				padded := make([]byte, length)
				copy(padded, data)
				data = padded
			}
			// For fixed-length strings, we don't write the length prefix
			_, err := buf.Write(data)
			return err
		}
	}

	// Default format: len(data) + data
	length := uint32(len(data))
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(data)
	return err
}

// encodeBytes handles serialization of []byte and [N]byte
func encodeBytes(b []byte, buf *bytes.Buffer, tag string) error {
	// Check if tag specifies length
	if tag != "" {
		if length, err := parseTag(tag); err == nil {
			if uint32(len(b)) > length {
				// Truncate data if it's longer than specified length
				b = b[:length]
			} else if uint32(len(b)) < length {
				// Pad with zeros if data is shorter than specified length
				padded := make([]byte, length)
				copy(padded, b)
				b = padded
			}
			// For fixed-length bytes, we don't write the length prefix
			_, err := buf.Write(b)
			return err
		}
	}

	// Default format: len(data) + data
	length := uint32(len(b))
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(b)
	return err
}

// encodeSlice handles serialization of slices (except []byte)
func encodeSlice(slice reflect.Value, buf *bytes.Buffer, tag string) error {
	// Check if tag specifies length
	if tag != "" {
		if length, err := parseTag(tag); err == nil {
			// For fixed-length slices, we don't write the length prefix
			// Write elements, padding with zero values if necessary
			sliceLen := uint32(slice.Len())
			elemType := slice.Type().Elem()

			for i := uint32(0); i < length; i++ {
				var elem reflect.Value
				if i < sliceLen {
					// Use actual element from slice
					elem = slice.Index(int(i))
				} else {
					// Create zero value for element
					elem = reflect.Zero(elemType)
				}

				if err := encodeField(elem, buf, ""); err != nil {
					return err
				}
			}
			return nil
		} else if tag == "-" {
			// If tag is "-", use default format
		}
	}

	// Default format: len(slice) + elements
	length := uint32(slice.Len())
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return err
	}

	// Write each element
	for i := 0; i < int(length); i++ {
		elem := slice.Index(i)
		if err := encodeField(elem, buf, ""); err != nil {
			return err
		}
	}

	return nil
}

// encodeArray handles serialization of arrays (except [N]byte)
func encodeArray(array reflect.Value, buf *bytes.Buffer, tag string) error {
	// Check if tag specifies length
	if tag != "" {
		if length, err := parseTag(tag); err == nil {
			// For fixed-length arrays, we don't write the length prefix
			// Write elements, padding with zero values if necessary
			arrayLen := uint32(array.Len())
			elemType := array.Type().Elem()

			for i := uint32(0); i < length; i++ {
				var elem reflect.Value
				if i < arrayLen {
					// Use actual element from array
					elem = array.Index(int(i))
				} else {
					// Create zero value for element
					elem = reflect.Zero(elemType)
				}

				if err := encodeField(elem, buf, ""); err != nil {
					return err
				}
			}
			return nil
		} else if tag == "-" {
			// If tag is "-", use default format
		}
	}

	// Default format: len(array) + elements
	length := uint32(array.Len())
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return err
	}

	// Write each element
	for i := 0; i < int(length); i++ {
		elem := array.Index(i)
		if err := encodeField(elem, buf, ""); err != nil {
			return err
		}
	}

	return nil
}

// parseTag parses the tag to extract length specification
func parseTag(tag string) (uint32, error) {
	if tag == "" {
		return 0, fmt.Errorf("empty tag")
	}

	// If tag is "-", it means to ignore the tag
	if tag == "-" {
		return 0, fmt.Errorf("ignore tag")
	}

	// Try to parse as integer
	if length, err := strconv.ParseUint(tag, 10, 32); err == nil {
		return uint32(length), nil
	}

	// Try to parse as "len:N" format
	if strings.HasPrefix(tag, "len:") {
		parts := strings.Split(tag, ":")
		if len(parts) == 2 {
			if length, err := strconv.ParseUint(parts[1], 10, 32); err == nil {
				return uint32(length), nil
			}
		}
	}

	return 0, fmt.Errorf("invalid tag format: %s", tag)
}

// decodeStruct handles deserialization of a struct
func decodeStruct(buf *bytes.Reader, val reflect.Value) error {
	typ := val.Type()
	numField := val.NumField()

	for i := 0; i < numField; i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Check if field implements BinaryDecoder
		if field.Kind() == reflect.Struct {
			// Create a pointer to the field for interface check
			fieldPtr := reflect.New(field.Type())
			fieldPtr.Elem().Set(field)

			if decoder, ok := fieldPtr.Interface().(BinaryDecoder); ok {
				// Read length
				var length uint32
				if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
					return err
				}
				// Read data
				data := make([]byte, length)
				if _, err := buf.Read(data); err != nil {
					return err
				}
				// Decode the field
				if err := decoder.Decode(data); err != nil {
					return fmt.Errorf("error decoding field %s: %w", fieldType.Name, err)
				}
				// Set the field
				field.Set(fieldPtr.Elem())
				continue
			}
		}

		tag := fieldType.Tag.Get("binary")
		// If tag is "-", skip this field entirely
		if tag == "-" {
			continue
		}

		if err := decodeField(buf, field, tag); err != nil {
			return fmt.Errorf("error decoding field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// decodeField handles deserialization of a single field
func decodeField(buf *bytes.Reader, field reflect.Value, tag string) error {
	switch field.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64:
		// For basic numeric types, we need to pass a pointer to binary.Read
		if field.CanAddr() {
			return binary.Read(buf, binary.LittleEndian, field.Addr().Interface())
		} else {
			// For non-addressable values (like array elements), we need to read into a temporary variable
			temp := reflect.New(field.Type()).Elem()
			err := binary.Read(buf, binary.LittleEndian, temp.Addr().Interface())
			if err != nil {
				return err
			}
			field.Set(temp)
			return nil
		}

	case reflect.String:
		return decodeString(buf, field, tag)

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			// []byte
			return decodeBytes(buf, field, tag)
		}
		// Other slices
		return decodeSlice(buf, field, tag)

	case reflect.Array:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			// [N]byte
			return decodeByteArray(buf, field, tag)
		}
		// Other arrays
		return decodeArray(buf, field, tag)

	case reflect.Struct:
		return decodeStruct(buf, field)

	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
}

// decodeString handles deserialization of strings
func decodeString(buf *bytes.Reader, field reflect.Value, tag string) error {
	var data []byte
	var err error

	// Check if tag specifies length
	if tag != "" {
		if length, parseErr := parseTag(tag); parseErr == nil {
			data = make([]byte, length)
			if _, err = buf.Read(data); err != nil {
				return err
			}
			// Trim trailing zeros
			data = bytes.TrimRight(data, "\x00")
			field.SetString(string(data))
			return nil
		}
	}

	// Default format: len(data) + data
	var length uint32
	if err = binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return err
	}
	data = make([]byte, length)
	if _, err = buf.Read(data); err != nil {
		return err
	}

	field.SetString(string(data))
	return nil
}

// decodeBytes handles deserialization of []byte
func decodeBytes(buf *bytes.Reader, field reflect.Value, tag string) error {
	var data []byte
	var err error

	// Check if tag specifies length
	if tag != "" {
		if length, parseErr := parseTag(tag); parseErr == nil {
			data = make([]byte, length)
			if _, err = buf.Read(data); err != nil {
				return err
			}
			field.SetBytes(data)
			return nil
		}
	}

	// Default format: len(data) + data
	var length uint32
	if err = binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return err
	}
	data = make([]byte, length)
	if _, err = buf.Read(data); err != nil {
		return err
	}

	field.SetBytes(data)
	return nil
}

// decodeByteArray handles deserialization of [N]byte
func decodeByteArray(buf *bytes.Reader, field reflect.Value, tag string) error {
	var data []byte
	var err error

	// Check if tag specifies length
	if tag != "" {
		if length, parseErr := parseTag(tag); parseErr == nil {
			data = make([]byte, length)
			if _, err = buf.Read(data); err != nil {
				return err
			}

			// Copy data to array, truncating or padding as necessary
			arrayLen := field.Len()
			copyLen := len(data)
			if copyLen > arrayLen {
				copyLen = arrayLen
			}

			// Copy data to array
			for i := 0; i < copyLen; i++ {
				field.Index(i).SetUint(uint64(data[i]))
			}

			// Zero out remaining elements if data is shorter than array
			for i := copyLen; i < arrayLen; i++ {
				field.Index(i).SetUint(0)
			}

			return nil
		}
	}

	// Default format: len(data) + data
	var length uint32
	if err = binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return err
	}
	data = make([]byte, length)
	if _, err = buf.Read(data); err != nil {
		return err
	}

	// Copy data to array, truncating or padding as necessary
	arrayLen := field.Len()
	copyLen := len(data)
	if copyLen > arrayLen {
		copyLen = arrayLen
	}

	// Copy data to array
	for i := 0; i < copyLen; i++ {
		field.Index(i).SetUint(uint64(data[i]))
	}

	// Zero out remaining elements if data is shorter than array
	for i := copyLen; i < arrayLen; i++ {
		field.Index(i).SetUint(0)
	}

	return nil
}

// decodeSlice handles deserialization of slices (except []byte)
func decodeSlice(buf *bytes.Reader, field reflect.Value, tag string) error {
	// Check if tag specifies length
	if tag != "" {
		if length, err := parseTag(tag); err == nil {
			// Get slice type and element type
			sliceType := field.Type()

			// For fixed-length slices, we don't read a length prefix
			// Create slice with the specified fixed length
			newSlice := reflect.MakeSlice(sliceType, int(length), int(length))

			// Read elements directly
			for i := uint32(0); i < length; i++ {
				elem := newSlice.Index(int(i))
				if err := decodeField(buf, elem, ""); err != nil {
					return err
				}
			}

			field.Set(newSlice)
			return nil
		}
	}

	// Default format: len(slice) + elements
	var length uint32
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return err
	}

	// Create slice
	sliceType := field.Type()
	newSlice := reflect.MakeSlice(sliceType, int(length), int(length))

	// Read each element
	for i := 0; i < int(length); i++ {
		elem := newSlice.Index(i)
		if err := decodeField(buf, elem, ""); err != nil {
			return err
		}
	}

	field.Set(newSlice)
	return nil
}

// decodeArray handles deserialization of arrays (except [N]byte)
func decodeArray(buf *bytes.Reader, field reflect.Value, tag string) error {
	// Check if tag specifies length
	if tag != "" {
		if length, err := parseTag(tag); err == nil {
			// Get array type and length
			arrayType := field.Type()
			arrayLen := uint32(arrayType.Len())

			// For fixed-length arrays, we don't read a length prefix
			// Read elements directly
			for i := uint32(0); i < length; i++ {
				if i < arrayLen {
					// Read actual element into array
					elem := field.Index(int(i))
					if err := decodeField(buf, elem, ""); err != nil {
						return err
					}
				} else {
					// Skip extra elements by reading into a temporary value
					temp := reflect.New(arrayType.Elem()).Elem()
					if err := decodeField(buf, temp, ""); err != nil {
						return err
					}
				}
			}

			// Zero out remaining elements if data is shorter than array
			for i := length; i < arrayLen; i++ {
				field.Index(int(i)).Set(reflect.Zero(arrayType.Elem()))
			}

			return nil
		}
	}

	// Default format: len(array) + elements
	var length uint32
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return err
	}

	// Get array type and length
	arrayType := field.Type()
	arrayLen := uint32(arrayType.Len())
	elemType := arrayType.Elem()

	// Read elements, padding with zero values if necessary
	for i := uint32(0); i < length; i++ {
		if i < arrayLen {
			// Read actual element into array
			elem := field.Index(int(i))
			if err := decodeField(buf, elem, ""); err != nil {
				return err
			}
		} else {
			// Skip extra elements by reading into a temporary value
			temp := reflect.New(elemType).Elem()
			if err := decodeField(buf, temp, ""); err != nil {
				return err
			}
		}
	}

	// Zero out remaining elements if data is shorter than array
	for i := length; i < arrayLen; i++ {
		field.Index(int(i)).Set(reflect.Zero(elemType))
	}

	return nil
}
