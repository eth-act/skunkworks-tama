package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"tamagotest/tamaboards/serialization"
)

func main() {
	// Hardcoded values - change these as needed
	x := int64(123)
	y := int64(456)

	// Output file - assumes running from project root
	output := "tama-programs/addition/witness.bin"

	// Create the witness file
	file, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// For generic Read[int64](), we need:
	// 1. Length (8 bytes) - length of serialized data
	// 2. Serialized data

	// Serialize first int64 using the serialization package
	xData := serialization.MustSerializeData(x)

	// Write length prefix (8 bytes, little-endian)
	err = binary.Write(file, binary.LittleEndian, uint64(len(xData)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing length for x: %v\n", err)
		os.Exit(1)
	}

	// Write serialized data
	_, err = file.Write(xData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing x data: %v\n", err)
		os.Exit(1)
	}

	// Serialize second int64 using the serialization package
	yData := serialization.MustSerializeData(y)

	// Write length prefix
	err = binary.Write(file, binary.LittleEndian, uint64(len(yData)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing length for y: %v\n", err)
		os.Exit(1)
	}

	// Write serialized data
	_, err = file.Write(yData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing y data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created witness file: %s\n", output)
	fmt.Printf("Values: x=%d, y=%d\n", x, y)
	fmt.Printf("Expected sum: %d\n", x+y)
	fmt.Printf("Serialized lengths: x=%d bytes, y=%d bytes\n", len(xData), len(yData))
	totalSize := 8 + len(xData) + 8 + len(yData)
	fmt.Printf("Total file size: %d bytes\n", totalSize)
}
