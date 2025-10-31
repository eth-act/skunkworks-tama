package main

import (
	"encoding/binary"
	"os"
	"tamagotest/tamaboards/serialization"
)

func main() {
	// Hardcoded values - change these as needed
	x := int64(123)
	y := int64(456)

	// Serialize first int64 using the serialization package
	xData := serialization.MustSerializeData(x)

	// Write length prefix (8 bytes, little-endian)
	err := binary.Write(os.Stdout, binary.LittleEndian, uint64(len(xData)))
	if err != nil {
		os.Exit(1)
	}

	// Write serialized data
	_, err = os.Stdout.Write(xData)
	if err != nil {
		os.Exit(1)
	}

	// Serialize second int64 using the serialization package
	yData := serialization.MustSerializeData(y)

	// Write length prefix
	err = binary.Write(os.Stdout, binary.LittleEndian, uint64(len(yData)))
	if err != nil {
		os.Exit(1)
	}

	// Write serialized data
	_, err = os.Stdout.Write(yData)
	if err != nil {
		os.Exit(1)
	}
}
