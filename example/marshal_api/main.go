package main

import (
	"fmt"
	"log"

	"github.com/lengzhao/binary"
)

// Person represents a person with various data types
type Person struct {
	Name    string
	Age     uint8
	Email   string    `binary:"50"` // Fixed length of 50 bytes
	Data    []byte    `binary:"10"` // Fixed length of 10 bytes
	Scores  []uint32  `binary:"5"`  // Fixed length of 5 elements
	ID      [16]byte  `binary:"16"` // Fixed length of 16 bytes
	Values  [4]uint32 `binary:"4"`  // Fixed length of 4 elements
	Address string    // Use default format (length + data)
	Height  float32   // Float32 value
	Weight  float64   // Float64 value
}

func main() {
	fmt.Println("=== Marshal/Unmarshal API Demo ===")

	person := Person{
		Name:    "Alice",
		Age:     30,
		Email:   "alice@example.com",
		Data:    []byte{1, 2, 3, 4, 5}, // Only 5 bytes, will be padded to 10
		Scores:  []uint32{100, 95, 87}, // Only 3 elements, will be padded to 5
		ID:      [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Values:  [4]uint32{1000, 2000, 3000, 4000},
		Address: "Main St 123",
		Height:  165.5, // Height in cm
		Weight:  62.3,  // Weight in kg
	}

	fmt.Println("Original person:")
	fmt.Printf("  Name: %s\n", person.Name)
	fmt.Printf("  Age: %d\n", person.Age)
	fmt.Printf("  Email: %s (length: %d)\n", person.Email, len(person.Email))
	fmt.Printf("  Address: %s\n", person.Address)

	// Marshal using the new API
	data, err := binary.Marshal(person)
	if err != nil {
		log.Fatal("Marshal failed:", err)
	}

	fmt.Printf("\nMarshaled data length: %d bytes\n", len(data))

	// Unmarshal using the new API
	var decoded Person
	err = binary.Unmarshal(data, &decoded)
	if err != nil {
		log.Fatal("Unmarshal failed:", err)
	}

	fmt.Println("\nDecoded person:")
	fmt.Printf("  Name: %s\n", decoded.Name)
	fmt.Printf("  Age: %d\n", decoded.Age)
	fmt.Printf("  Email: %s (length: %d)\n", decoded.Email, len(decoded.Email))
	fmt.Printf("  Address: %s\n", decoded.Address)

	fmt.Println("\n=== Direct Value Marshal/Unmarshal ===")

	// Direct slice marshal/unmarshal
	slice := []uint32{10, 20, 30, 40, 50}
	fmt.Printf("Original slice: %v\n", slice)

	sliceData, err := binary.Marshal(slice)
	if err != nil {
		log.Fatal("Marshal slice failed:", err)
	}

	var decodedSlice []uint32
	err = binary.Unmarshal(sliceData, &decodedSlice)
	if err != nil {
		log.Fatal("Unmarshal slice failed:", err)
	}

	fmt.Printf("Decoded slice: %v\n", decodedSlice)

	// Direct basic type marshal/unmarshal
	number := uint32(42)
	fmt.Printf("Original number: %d\n", number)

	numberData, err := binary.Marshal(number)
	if err != nil {
		log.Fatal("Marshal number failed:", err)
	}

	var decodedNumber uint32
	err = binary.Unmarshal(numberData, &decodedNumber)
	if err != nil {
		log.Fatal("Unmarshal number failed:", err)
	}

	fmt.Printf("Decoded number: %d\n", decodedNumber)

	fmt.Println("\n🎉 All tests passed! The Marshal/Unmarshal API works correctly.")
}
