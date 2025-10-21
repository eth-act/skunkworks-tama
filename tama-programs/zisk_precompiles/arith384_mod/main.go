//go:build tamago && riscv64

package main

import (
  "fmt"
  "math/big"
  "unsafe"
  _ "tamagotest/tamaboards/zkvm"
)

// SyscallArith384ModParams represents the parameters for arith256_mod syscall
// see https://github.com/0xPolygonHermez/zisk/blob/main/ziskos/entrypoint/src/syscalls/arith256_mod.rs
type SyscallArith384ModParams struct {
  A      *[6]uint64
  B      *[6]uint64
  C      *[6]uint64
  Module *[6]uint64
  D      *[6]uint64
}

//go:noinline
func ZiskArith384_mod(arr *SyscallArith384ModParams) {
  ptr := uintptr(unsafe.Pointer(arr))
  arith384_mod(ptr)
}

func arith384_mod(ptr uintptr)

// uint64ArrayToBigInt converts a [6]uint64 array to big.Int (little-endian)
func uint64ArrayToBigInt(arr *[6]uint64) *big.Int {
  result := new(big.Int)
  for i := 5; i >= 0; i-- {
    result.Lsh(result, 64)
    result.Or(result, new(big.Int).SetUint64(arr[i]))
  }
  return result
}

// printUint384 prints a 384-bit number from [6]uint64 array
func printUint384(name string, arr *[6]uint64) {
  value := uint64ArrayToBigInt(arr)
  fmt.Printf("%s: %s\n", name, value.String())
}

// verifyComputation verifies that d = (a * b + c) mod module
func verifyComputation(params *SyscallArith384ModParams) bool {
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

  a := [6]uint64{0x112210F4B2D230D2, 0x74172ed0ab34612c, 0x5d46fa993452fef9, 0x657aaf58cb17a53c, 0x6bed8b7acd3824ee, 0xd7479084dc15a54e}

  b := [6]uint64{0x84B2C7C5BC5C597A, 0x2e6371d577af32ab, 0x3c8f0bc543b1286c, 0x037f0aa5d1c48341, 0xd7c8f233caee3c6d, 0x68ea856d2dbb74fd}

  c := [6]uint64{0xC19E6C2252E377C7, 0xc864ca220f556c33, 0xae15c7d0cd7774ba, 0x112ddde737d94162, 0x3771af4aa8234d7d, 0x2ceacd17137894c5}

  // module = A large prime number (not a power of 2)
  // module = 0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb15
  // This is the bls12-381 field prime (used in Bitcoin/Ethereum)
  module := [6]uint64{
   0x0f6241eabfffeb15,
   0x512bf6730d2a0f6b,
   0xbacd764774b84f38,
   0xfe69a4b1ba7b6434,
   0x000001a0111ea397,
   0x0000000000000000,
  }

  // d will hold the result
  d := [6]uint64{0, 0, 0, 0, 0, 0}

  // Create the struct instance
  params := SyscallArith384ModParams{
    A:      &a,
    B:      &b,
    C:      &c,
    Module: &module,
    D:      &d,
  }

  fmt.Println("Input values:")
  fmt.Println("=============")
  printUint384("a", params.A)
  printUint384("b", params.B)
  printUint384("c", params.C)
  printUint384("module", params.Module)
  fmt.Println()

  // Call the syscall function
  ZiskArith384_mod(&params)

  fmt.Println("Output:")
  fmt.Println("=======")
  printUint384("d (result)", params.D)

  // Verify the computation
  verifyComputation(&params)
}
