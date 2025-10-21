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

type Complex struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallBn254ComplexAddParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBn254ComplexMulParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBn254ComplexSubParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBn254AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [4]uint64{0x0000000000000001, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
	Y: [4]uint64{0x0000000000000002, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
}

var Neutral = Point256{
	X: [4]uint64{0, 0, 0, 0},
	Y: [4]uint64{0, 0, 0, 0},
}

//go:noinline
func ZiskBn254Add(arr *SyscallBn254AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bn254_add(ptr)
}

//go:noinline
func ZiskBn254Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	bn254_dbl(ptr)
}

//go:noinline
func ZiskBn254ComplexAdd(arr *SyscallBn254ComplexAddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bn254_complex_add(ptr)
}

//go:noinline
func ZiskBn254ComplexMul(arr *SyscallBn254ComplexMulParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bn254_complex_mul(ptr)
}

//go:noinline
func ZiskBn254ComplexSub(arr *SyscallBn254ComplexSubParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bn254_complex_sub(ptr)
}

func bn254_add(ptr uintptr)
func bn254_dbl(ptr uintptr)
func bn254_complex_add(ptr uintptr)
func bn254_complex_mul(ptr uintptr)
func bn254_complex_sub(ptr uintptr)

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
# Define the bn254 curve parameters
p = 21888242871839275222246405745257275088696311157297823662689037894645226208583
a = 0
b = 3

# Create the elliptic curve
E = EllipticCurve(GF(p), [a, b])

# Define the generator point G
G_x = 1
G_y = 2
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

	params := SyscallBn254AddParams{
		P1: &p1,
		P2: &p2,
	}

	ZiskBn254Dbl(&p1)
	ZiskBn254Add(&params)

	print_memory_repr(&p1)

	f1 := Complex{
		X: [4]uint64{2, 0, 0, 0},
		Y: [4]uint64{3, 0, 0, 0},
	}
	f2 := Complex{
		X: [4]uint64{4, 0, 0, 0},
		Y: [4]uint64{5, 0, 0, 0},
	}

	ZiskBn254ComplexAdd(&SyscallBn254ComplexAddParams{
		F1: &f1,
		F2: &f2,
	})

	print_memory_repr(&f1)

	ZiskBn254ComplexMul(&SyscallBn254ComplexMulParams{
		F1: &f1,
		F2: &f2,
	})

	print_memory_repr(&f1)

	ZiskBn254ComplexSub(&SyscallBn254ComplexSubParams{
		F1: &f1,
		F2: &f2,
	})

	print_memory_repr(&f1)

	if p2 == Generator {
		fmt.Println("Success")
	} else {
		fmt.Println("Failure")
	}
}
