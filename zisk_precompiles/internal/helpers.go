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

func DecStringToUint64x4(s string) [4]uint64 {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("invalid decimal string: " + s)
	}
	mask := new(big.Int).SetUint64(0xffffffffffffffff)
	var result [4]uint64
	for i := 0; i < 4; i++ {
		result[i] = new(big.Int).And(n, mask).Uint64()
		n.Rsh(n, 64)
	}
	return result
}

func DecStringToUint64x6(s string) [6]uint64 {
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("invalid decimal string: " + s)
	}
	mask := new(big.Int).SetUint64(0xffffffffffffffff)
	var result [6]uint64
	for i := 0; i < 6; i++ {
		result[i] = new(big.Int).And(n, mask).Uint64()
		n.Rsh(n, 64)
	}
	return result
}

func HexStringToUint64x6(s string) [6]uint64 {
	if len(s) >= 2 && s[:2] == "0x" || len(s) >= 2 && s[:2] == "0X" {
		s = s[2:]
	}
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid hex string: " + s)
	}
	mask := new(big.Int).SetUint64(0xffffffffffffffff)
	var result [6]uint64
	for i := 0; i < 6; i++ {
		result[i] = new(big.Int).And(n, mask).Uint64()
		n.Rsh(n, 64)
	}
	return result
}

func PrintMemoryRepr[T any](point *T) {
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
