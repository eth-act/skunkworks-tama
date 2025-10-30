//go:build tamago && riscv64

package main

import (
  "fmt"
  "math/big"
  "unsafe"
  _ "tamagotest/tamaboards/zkvm"
)

// SyscallArith256ModParams represents the parameters for arith256_mod syscall
// see https://github.com/0xPolygonHermez/zisk/blob/main/ziskos/entrypoint/src/syscalls/arith256_mod.rs
type SyscallArith256ModParams struct {
  A      *[4]uint64
  B      *[4]uint64
  C      *[4]uint64
  Module *[4]uint64
  D      *[4]uint64
}

//go:noinline
func ZiskArith256_mod(arr *SyscallArith256ModParams) {
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
func verifyComputation(params *SyscallArith256ModParams) bool {
  a := uint64ArrayToBigInt(params.A)
  b := uint64ArrayToBigInt(params.B)
  c := uint64ArrayToBigInt(params.C)
  module := uint64ArrayToBigInt(params.Module)
  d := uint64ArrayToBigInt(params.D)

  // Compute expected result: (a * b + c) mod module
  expected := new(big.Int)
  expected.Mul(a, b)        // a * b
  expected.Add(expected, c) // + c
  expected.Mod(expected, module) // mod module

  fmt.Println("\nVerification:")
  fmt.Println("=============")
  fmt.Printf("Expected: %s\n", expected.String())
  fmt.Printf("Got:      %s\n", d.String())

  match := expected.Cmp(d) == 0
  if match {
    fmt.Println("Result is CORRECT!")
  } else {
    fmt.Println("Result is INCORRECT!")
  }

  return match
}

func main() {
  // Sample 256-bit values
  a := [4]uint64{0x112210F4B2D230D2, 0x0000000000000019, 0x0000000000000000, 0x0000000000000000}

  b := [4]uint64{0x84B2C7C5BC5C597A, 0x00000000000000CB, 0x0000000000000000, 0x0000000000000000}

  c := [4]uint64{0xC19E6C2252E377C7, 0x0000000000000017, 0x0000000000000000, 0x0000000000000000}

  // module = A large prime number (not a power of 2)
  // module = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
  // This is the secp256k1 field prime (used in Bitcoin/Ethereum)
  module := [4]uint64{
    0xFFFFFFFEFFFFFC2F, // least significant
    0xFFFFFFFFFFFFFFFF,
    0xFFFFFFFFFFFFFFFF,
    0xFFFFFFFFFFFFFFFF, // most significant
  }

  // d will hold the result
  d := [4]uint64{0, 0, 0, 0}

  // Create the struct instance
  params := SyscallArith256ModParams{
    A:      &a,
    B:      &b,
    C:      &c,
    Module: &module,
    D:      &d,
  }

  fmt.Println("Input values:")
  fmt.Println("=============")
  printUint256("a", params.A)
  printUint256("b", params.B)
  printUint256("c", params.C)
  printUint256("module", params.Module)
  fmt.Println()

  // Call the syscall function
  ZiskArith256_mod(&params)

  fmt.Println("Output:")
  fmt.Println("=======")
  printUint256("d (result)", params.D)

  // Verify the computation
  verifyComputation(&params)
}
