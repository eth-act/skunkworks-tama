//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
	"unsafe"
)

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallSecp256k1AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [4]uint64{0x59F2815B16F81798, 0x029BFCDB2DCE28D9, 0x55A06295CE870B07, 0x79BE667EF9DCBBAC},
	Y: [4]uint64{0x9C47D08FFB10D4B8, 0xFD17B448A6855419, 0x5DA4FBFC0E1108A8, 0x483ADA7726A3C465},
}

var Neutral = Point256{
	X: [4]uint64{0, 0, 0, 0},
	Y: [4]uint64{0, 0, 0, 0},
}

// ziskos/entrypoint/src/zisklib/lib/secp256k1/curve.rs:44
// Given points `p1` and `p2`, performs the point addition `p1 + p2` and assigns the result to `p1`.
// It assumes that `p1` and `p2` are from the Secp256k1 curve, that `p1,p2 != 𝒪` and that `p2 != p1,-p1`

//go:noinline
func ZiskSecp256k1Add(arr *SyscallSecp256k1AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	secp256k1_add(ptr)
}

//go:noinline
func ZiskSecp256k1Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	secp256k1_dbl(ptr)
}

func secp256k1_add(ptr uintptr)
func secp256k1_dbl(ptr uintptr)

func print_memory_repr[T any](point *T) {
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

/*
# Define the secp256k1 curve parameters
p = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
a = 0
b = 7

# Create the elliptic curve
E = EllipticCurve(GF(p), [a, b])

# Define the generator point G
G_x = 0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798
G_y = 0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8
G = E(G_x, G_y)

# Add the generator to itself: G + G = 2G
result = G + G

# Get the affine coordinates
x_coord = result[0]
y_coord = result[1]

# Print the result
print("Generator point G:")
print(f"x = {hex(G_x)}")
print(f"y = {hex(G_y)}")
print("\nG + G (2G) result:")
print(f"x = {hex(x_coord)}")
print(f"y = {hex(y_coord)}")
*/

func main() {
	// Sample 256-bit values
	p1 := Generator
	p2 := Generator

	params := SyscallSecp256k1AddParams{
		P1: &p1,
		P2: &p2,
	}

	ZiskSecp256k1Dbl(&p1)
	ZiskSecp256k1Add(&params)

	print_memory_repr(&p1)

	if p2 == Generator {
		fmt.Println("Success")
	} else {
		fmt.Println("Failure")
	}
}
