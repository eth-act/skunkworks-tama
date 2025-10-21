//go:build tamago && riscv64

package main

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
	"unsafe"
)

// SyscallSha256Params represents the parameters for sha256 syscall
// see https://github.com/0xPolygonHermez/zisk/blob/main/ziskos/entrypoint/src/syscalls/sha256.rs
type SyscallSha256Params struct {
	State *[4]uint64
	Input *[8]uint64
}

//go:noinline
func ZiskSha256(arr *SyscallSha256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	sha256(ptr)
}

func sha256(ptr uintptr)

func main() {
	// Sample 256-bit values
	state := [4]uint64{0xbb67ae856a09e667, 0xa54ff53a3c6ef372, 0x9b05688c510e527f, 0x5be0cd191f83d9ab}
	input := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}

	// Create the struct instance
	params := SyscallSha256Params{State: &state, Input: &input}

	// Call the syscall function
	ZiskSha256(&params)

	// TODO: verify the result
	fmt.Println("state: {}, input: {}", state, input)
}
