package main

import (
	"bytes"
	"fmt"

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
	// TempField string   `binary:"-"`      // Example of ignored field
}

func main() {
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
	fmt.Printf("  Data: %v (length: %d)\n", person.Data, len(person.Data))
	fmt.Printf("  Scores: %v (length: %d)\n", person.Scores, len(person.Scores))
	fmt.Printf("  ID: %v (length: %d)\n", person.ID, len(person.ID))
	fmt.Printf("  Values: %v (length: %d)\n", person.Values, len(person.Values))
	fmt.Printf("  Address: %s\n", person.Address)
	fmt.Printf("  Height: %.1f cm\n", person.Height)
	fmt.Printf("  Weight: %.1f kg\n", person.Weight)

	// Serialize
	data, err := binary.Marshal(person)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSerialized data length: %d bytes\n", len(data))

	// Deserialize
	var decoded Person
	err = binary.Unmarshal(data, &decoded)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nDecoded person:")
	fmt.Printf("  Name: %s\n", decoded.Name)
	fmt.Printf("  Age: %d\n", decoded.Age)
	fmt.Printf("  Email: %s (length: %d)\n", decoded.Email, len(decoded.Email))
	fmt.Printf("  Data: %v (length: %d)\n", decoded.Data, len(decoded.Data))
	fmt.Printf("  Scores: %v (length: %d)\n", decoded.Scores, len(decoded.Scores))
	fmt.Printf("  ID: %v (length: %d)\n", decoded.ID, len(decoded.ID))
	fmt.Printf("  Values: %v (length: %d)\n", decoded.Values, len(decoded.Values))
	fmt.Printf("  Address: %s\n", decoded.Address)
	fmt.Printf("  Height: %.1f cm\n", decoded.Height)
	fmt.Printf("  Weight: %.1f kg\n", decoded.Weight)

	// Check if original and decoded persons are compatible
	// Use epsilon comparison for floating point values
	heightEqual := abs32(person.Height-decoded.Height) < 1e-6
	weightEqual := abs64(person.Weight-decoded.Weight) < 1e-12

	// For tagged fields, we need to compare the original values only
	dataEqual := bytes.Equal(person.Data, decoded.Data[:len(person.Data)])         // Compare only original data
	scoresEqual := slicesEqual(person.Scores, decoded.Scores[:len(person.Scores)]) // Compare only original scores

	if person.Name == decoded.Name &&
		person.Age == decoded.Age &&
		person.Email == decoded.Email &&
		dataEqual &&
		scoresEqual &&
		arraysEqual(person.ID[:], decoded.ID[:]) &&
		arraysEqualUint32(person.Values[:], decoded.Values[:]) &&
		person.Address == decoded.Address &&
		heightEqual &&
		weightEqual {
		fmt.Println("\nSuccess: Original and decoded persons are compatible!")
	} else {
		fmt.Println("\nError: Original and decoded persons are not compatible!")
		fmt.Printf("  Name: %v == %v\n", person.Name, decoded.Name)
		fmt.Printf("  Age: %v == %v\n", person.Age, decoded.Age)
		fmt.Printf("  Email: %v == %v\n", person.Email, decoded.Email)
		fmt.Printf("  Data equal: %v\n", dataEqual)
		fmt.Printf("  Scores equal: %v\n", scoresEqual)
		fmt.Printf("  ID: %v == %v\n", person.ID, decoded.ID)
		fmt.Printf("  Values: %v == %v\n", person.Values, decoded.Values)
		fmt.Printf("  Address: %v == %v\n", person.Address, decoded.Address)
		fmt.Printf("  Height: %v == %v (diff: %v)\n", person.Height, decoded.Height, abs32(person.Height-decoded.Height))
		fmt.Printf("  Weight: %v == %v (diff: %v)\n", person.Weight, decoded.Weight, abs64(person.Weight-decoded.Weight))
	}
}

// Helper functions for comparison
func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func slicesEqual(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func arraysEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func arraysEqualUint32(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
