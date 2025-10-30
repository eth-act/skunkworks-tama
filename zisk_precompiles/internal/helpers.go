package internal

import (
	"fmt"
	"math/big"
	"unsafe"
)

// uint64ArrayToBigInt converts a [4]uint64 array to big.Int (little-endian)
func Uint64ArrayToBigInt(arr *[4]uint64) *big.Int {
	result := new(big.Int)
	for i := 3; i >= 0; i-- {
		result.Lsh(result, 64)
		result.Or(result, new(big.Int).SetUint64(arr[i]))
	}
	return result
}

// printUint256 prints a 256-bit number from [4]uint64 array
func PrintUint256(name string, arr *[4]uint64) {
	value := Uint64ArrayToBigInt(arr)
	fmt.Printf("%s: %s\n", name, value.String())
}

func Print_memory_repr[T any](point *T) {
	// Get pointer to the struct
	ptr := unsafe.Pointer(point)

	// Convert to byte slice
	size := unsafe.Sizeof(*point)
	bytes := unsafe.Slice((*byte)(ptr), size)

	// Print bytes in hex
	fmt.Printf("Size: %d bytes\n\n", size)

	for i, b := range bytes {
		fmt.Printf("%02x ", b)
		if (i+1)%8 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
