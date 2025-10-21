//go:build tamago && riscv64

package main

import (
	"fmt"
	"math/big"
	_ "tamagotest/tamaboards/zkvm"
	"unsafe"
)

// SyscallArith256Params represents the parameters for arith256_mod syscall
// see https://github.com/0xPolygonHermez/zisk/blob/main/ziskos/entrypoint/src/syscalls/arith256_mod.rs
type SyscallArith256Params struct {
	A  *[4]uint64
	B  *[4]uint64
	C  *[4]uint64
	Dl *[4]uint64
	Dh *[4]uint64
}

//go:noinline
func ZiskArith256_mod(arr *SyscallArith256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	arith256_mod(ptr)
}

func arith256_mod(ptr uintptr)

// uint64ArrayToBigInt converts a [4]uint64 array to big.Int (little-endian)
func uint64ArrayToBigInt(arr *[4]uint64) *big.Int {
	result := new(big.Int)
	for i := 3; i >= 0; i-- {
		result.Lsh(result, 64)
		result.Or(result, new(big.Int).SetUint64(arr[i]))
	}
	return result
}

// printUint256 prints a 256-bit number from [4]uint64 array
func printUint256(name string, arr *[4]uint64) {
	value := uint64ArrayToBigInt(arr)
	fmt.Printf("%s: %s\n", name, value.String())
}

// verifyComputation verifies that d = (a * b + c) mod module
func verifyComputation(params *SyscallArith256Params) bool {
	a := uint64ArrayToBigInt(params.A)
	b := uint64ArrayToBigInt(params.B)
	c := uint64ArrayToBigInt(params.C)
	dl := uint64ArrayToBigInt(params.Dl)
	dh := uint64ArrayToBigInt(params.Dh)
	mask := uint64ArrayToBigInt(&[4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff})

	// Compute expected result: (a * b + c)
	expected_dhdl := new(big.Int)
	expected_dhdl.Mul(a, b)             // a * b
	expected_dhdl.Add(expected_dhdl, c) // + c

	expected_dl := new(big.Int)
	expected_dl.And(expected_dhdl, mask)

	expected_dh := new(big.Int)
	expected_dh.Rsh(expected_dhdl, 256)
	expected_dh.And(expected_dh, mask)

	fmt.Println("\nVerification:")
	fmt.Println("=============")
	fmt.Printf("Expected: %s, %s\n", expected_dh.String(), expected_dl.String())
	fmt.Printf("Got:      %s, %s\n", dh.String(), dl.String())

	match := expected_dl.Cmp(dl) == 0 && expected_dh.Cmp(dh) == 0
	if match {
		fmt.Println("Result is CORRECT!")
	} else {
		fmt.Println("Result is INCORRECT!")
	}

	return match
}

func main() {
	// Sample 256-bit values
	a := [4]uint64{0x112210F4B2D230D2, 0x74172ed0ab34612c, 0x5d46fa993452fef9, 0x657aaf58cb17a53c}

	b := [4]uint64{0x84B2C7C5BC5C597A, 0x2e6371d577af32ab, 0x3c8f0bc543b1286c, 0x037f0aa5d1c48341}

	c := [4]uint64{0xC19E6C2252E377C7, 0xc864ca220f556c33, 0xae15c7d0cd7774ba, 0x112ddde737d94162}

	dh := [4]uint64{0, 0, 0, 0}
	dl := [4]uint64{0, 0, 0, 0}

	// Create the struct instance
	params := SyscallArith256Params{
		A:  &a,
		B:  &b,
		C:  &c,
		Dh: &dh,
		Dl: &dl,
	}

	fmt.Println("Input values:")
	fmt.Println("=============")
	printUint256("a", params.A)
	printUint256("b", params.B)
	printUint256("c", params.C)
	fmt.Println()

	// Call the syscall function
	ZiskArith256_mod(&params)

	fmt.Println("Output:")
	fmt.Println("=======")
	printUint256("dh (result)", params.Dh)
	printUint256("dl (result)", params.Dl)

	// Verify the computation
	verifyComputation(&params)
}
