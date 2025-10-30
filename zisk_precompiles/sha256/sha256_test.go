package sha256

import (
	"fmt"
	_ "tamagotest/tamaboards/zkvm"
	"testing"
)

func TestSha256(t *testing.T) {
	// Sample 256-bit values
	state := [4]uint64{0xbb67ae856a09e667, 0xa54ff53a3c6ef372, 0x9b05688c510e527f, 0x5be0cd191f83d9ab}
	input := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}

	// Create the struct instance
	params := SyscallSha256Params{State: &state, Input: &input}

	// Call the syscall function
	Sha256Update(&params)

	// TODO: verify the result
	fmt.Println("state: {}, input: {}", state, input)
}
