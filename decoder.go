package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Unmarshal deserializes binary data into a value
// This function expects all data to be consumed and returns an error if there are remaining bytes
func Unmarshal(data []byte, v interface{}) error {
	remaining, err := UnmarshalPartial(data, v)
	if err != nil {
		return err
	}

	// Check for remaining data - this maintains backward compatibility
	if remaining > 0 {
		return fmt.Errorf("warning: %d bytes of data remaining after unmarshaling", remaining)
	}

	return nil
}

// UnmarshalPartial deserializes binary data into a value and returns the number of remaining bytes
// This allows for partial parsing of data streams where you might want to process multiple values
// sequentially or handle cases where the data contains more information than needed.
// Returns:
//   - remaining: number of bytes left unprocessed in the input data
//   - error: any error that occurred during unmarshaling
func UnmarshalPartial(data []byte, v interface{}) (remaining int, err error) {
	// Check if the value implements BinaryUnmarshaler
	if unmarshaler, ok := v.(BinaryUnmarshaler); ok {
		// For BinaryUnmarshaler, we consume all data and return 0 remaining
		// This maintains compatibility with existing implementations
		err = unmarshaler.UnmarshalBinary(data)
		return 0, err
	}

	val := reflect.ValueOf(v)

	// Check if v is a pointer
	if val.Kind() != reflect.Ptr {
		return len(data), fmt.Errorf("only pointers are supported for unmarshaling")
	}

	// Check if v is a nil pointer
	if val.IsNil() {
		return len(data), fmt.Errorf("cannot unmarshal into nil pointer")
	}

	// Get the element that the pointer points to
	elem := val.Elem()

	// Unmarshal any type by calling decodeField directly
	buf := bytes.NewReader(data)
	if err := decodeField(buf, elem, ""); err != nil {
		return buf.Len(), fmt.Errorf("error unmarshaling value: %w", err)
	}

	// Return the number of remaining bytes
	return buf.Len(), nil
}

// decodeField handles deserialization of a single field
func decodeField(buf *bytes.Reader, field reflect.Value, tag string) error {
	switch field.Kind() {
	case reflect.Ptr:
		// Handle pointer types by dereferencing them
		if field.IsNil() {
			// Create a new instance of the pointed-to type
			newValue := reflect.New(field.Type().Elem())
			field.Set(newValue)
		}
		return decodeField(buf, field.Elem(), tag)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
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

	case reflect.Float32, reflect.Float64:
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

		// Check if field implements BinaryUnmarshaler
		if field.Kind() == reflect.Struct {
			// Create a pointer to the field for interface check
			fieldPtr := reflect.New(field.Type())
			fieldPtr.Elem().Set(field)

			if unmarshaler, ok := fieldPtr.Interface().(BinaryUnmarshaler); ok {
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
				// Unmarshal the field
				if err := unmarshaler.UnmarshalBinary(data); err != nil {
					return fmt.Errorf("error unmarshaling field %s: %w", fieldType.Name, err)
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
