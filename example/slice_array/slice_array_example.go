package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/binary"
)

func main() {
	// Example 1: Directly encode/decode a slice
	fmt.Println("=== Example 1: Slice ===")
	slice := []uint32{10, 20, 30, 40, 50}
	fmt.Printf("Original slice: %v\n", slice)

	// Encode the slice directly
	data, err := binary.Encode(slice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded data length: %d bytes\n", len(data))

	// Decode the slice directly
	var decodedSlice []uint32
	err = binary.Decode(data, &decodedSlice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded slice: %v\n", decodedSlice)
	fmt.Println()

	// Example 2: Directly encode/decode an array
	fmt.Println("=== Example 2: Array ===")
	array := [5]uint32{100, 200, 300, 400, 500}
	fmt.Printf("Original array: %v\n", array)

	// Encode the array directly
	data, err = binary.Encode(array)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded data length: %d bytes\n", len(data))

	// Decode the array directly
	var decodedArray [5]uint32
	err = binary.Decode(data, &decodedArray)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded array: %v\n", decodedArray)
	fmt.Println()

	// Example 3: Directly encode/decode a byte slice
	fmt.Println("=== Example 3: Byte Slice ===")
	byteSlice := []byte{1, 2, 3, 4, 5}
	fmt.Printf("Original byte slice: %v\n", byteSlice)

	// Encode the byte slice directly
	data, err = binary.Encode(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded data length: %d bytes\n", len(data))

	// Decode the byte slice directly
	var decodedByteSlice []byte
	err = binary.Decode(data, &decodedByteSlice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded byte slice: %v\n", decodedByteSlice)
	fmt.Println()

	// Example 4: Directly encode/decode a byte array
	fmt.Println("=== Example 4: Byte Array ===")
	byteArray := [5]byte{10, 20, 30, 40, 50}
	fmt.Printf("Original byte array: %v\n", byteArray)

	// Encode the byte array directly
	data, err = binary.Encode(byteArray)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded data length: %d bytes\n", len(data))

	// Decode the byte array directly
	var decodedByteArray [5]byte
	err = binary.Decode(data, &decodedByteArray)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded byte array: %v\n", decodedByteArray)
	fmt.Println()

	// Example 5: Directly encode/decode a slice of strings
	fmt.Println("=== Example 5: String Slice ===")
	stringSlice := []string{"hello", "world", "golang", "binary"}
	fmt.Printf("Original string slice: %v\n", stringSlice)

	// Encode the string slice directly
	data, err = binary.Encode(stringSlice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encoded data length: %d bytes\n", len(data))

	// Decode the string slice directly
	var decodedStringSlice []string
	err = binary.Decode(data, &decodedStringSlice)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decoded string slice: %v\n", decodedStringSlice)
}
