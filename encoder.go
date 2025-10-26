package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// Marshal serializes a value into binary format
func Marshal(v interface{}) ([]byte, error) {
	// Check if the value implements BinaryMarshaler
	if marshaler, ok := v.(BinaryMarshaler); ok {
		return marshaler.MarshalBinary()
	}

	val := reflect.ValueOf(v)

	// Marshal any type by calling encodeField directly
	var buf bytes.Buffer
	tag := "" // No tag for direct encoding
	if err := encodeField(val, &buf, tag); err != nil {
		return nil, fmt.Errorf("error marshaling value: %w", err)
	}

	return buf.Bytes(), nil
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

		// Check if field implements BinaryMarshaler
		if marshaler, ok := field.Interface().(BinaryMarshaler); ok {
			fieldData, err := marshaler.MarshalBinary()
			if err != nil {
				return fmt.Errorf("error marshaling field %s: %w", fieldType.Name, err)
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
	case reflect.Ptr:
		// Handle pointer types by dereferencing them
		if field.IsNil() {
			return fmt.Errorf("cannot encode nil pointer")
		}
		return encodeField(field.Elem(), buf, tag)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int, reflect.Bool:
		return binary.Write(buf, binary.LittleEndian, field.Interface())

	case reflect.Float32, reflect.Float64:
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
			sliceLen := uint32(slice.Len())
			elemType := slice.Type().Elem()

			for i := uint32(0); i < length; i++ {
				var elem reflect.Value
				if i < sliceLen {
					elem = slice.Index(int(i))
				} else {
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
			arrayLen := uint32(array.Len())
			elemType := array.Type().Elem()

			for i := uint32(0); i < length; i++ {
				var elem reflect.Value
				if i < arrayLen {
					elem = array.Index(int(i))
				} else {
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
