//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
	"unsafe"
)

type Point256 struct {
	X [6]uint64
	Y [6]uint64
}

type Complex struct {
	X [6]uint64
	Y [6]uint64
}

type SyscallBls12_381ComplexAddParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBls12_381ComplexMulParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBls12_381ComplexSubParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBls12_381AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [6]uint64{0xfb3af00adb22c6bb, 0x6c55e83ff97a1aef, 0xa14e3a3f171bac58, 0xc3688c4f9774b905, 0x2695638c4fa9ac0f, 0x17f1d3a73197d794},
	Y: [6]uint64{0x0caa232946c5e7e1, 0xd03cc744a2888ae4, 0x00db18cb2c04b3ed, 0xfcf5e095d5d00af6, 0xa09e30ed741d8ae4, 0x08b3f481e3aaa0f1},
}

var Neutral = Point256{
	X: [6]uint64{0, 0, 0, 0, 0, 0},
	Y: [6]uint64{0, 0, 0, 0, 0, 0},
}

//go:noinline
func ZiskBls12_381Add(arr *SyscallBls12_381AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bls12_381_add(ptr)
}

//go:noinline
func ZiskBls12_381Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	bls12_381_dbl(ptr)
}

//go:noinline
func ZiskBls12_381ComplexAdd(arr *SyscallBls12_381ComplexAddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bls12_381_complex_add(ptr)
}

//go:noinline
func ZiskBls12_381ComplexMul(arr *SyscallBls12_381ComplexMulParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bls12_381_complex_mul(ptr)
}

//go:noinline
func ZiskBls12_381ComplexSub(arr *SyscallBls12_381ComplexSubParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	bls12_381_complex_sub(ptr)
}

func bls12_381_add(ptr uintptr)
func bls12_381_dbl(ptr uintptr)
func bls12_381_complex_add(ptr uintptr)
func bls12_381_complex_mul(ptr uintptr)
func bls12_381_complex_sub(ptr uintptr)

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
# Define the bls12_381 curve parameters
p = 0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab
a = 0
b = 4

# Create the elliptic curve
E = EllipticCurve(GF(p), [a, b])

# Define the generator point G
G_x = 0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb
G_y = 0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1
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

	params := SyscallBls12_381AddParams{
		P1: &p1,
		P2: &p2,
	}

	ZiskBls12_381Dbl(&p1)
	ZiskBls12_381Add(&params)

	print_memory_repr(&p1)

	f1 := Complex{
		X: [6]uint64{2, 0, 0, 0, 0, 0},
		Y: [6]uint64{3, 0, 0, 0, 0, 0},
	}
	f2 := Complex{
		X: [6]uint64{4, 0, 0, 0, 0, 0},
		Y: [6]uint64{5, 0, 0, 0, 0, 0},
	}

	ZiskBls12_381ComplexAdd(&SyscallBls12_381ComplexAddParams{
		F1: &f1,
		F2: &f2,
	})

	print_memory_repr(&f1)

	ZiskBls12_381ComplexMul(&SyscallBls12_381ComplexMulParams{
		F1: &f1,
		F2: &f2,
	})

	print_memory_repr(&f1)

	ZiskBls12_381ComplexSub(&SyscallBls12_381ComplexSubParams{
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
