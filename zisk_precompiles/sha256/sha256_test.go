package sha256

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

func TestSha256(t *testing.T) {
	// The initial state in `state` looks like the initial state described in the standard.
	// Zisk changes the endianess of the state and input so in reality in this test `state` is just random data
	// https://github.com/0xPolygonHermez/zisk/blob/0207f975d0f03724f6f5784692ee2725e0d340d3/ziskclib/src/helpers.rs
	state := [4]uint64{0xbb67ae856a09e667, 0xa54ff53a3c6ef372, 0x9b05688c510e527f, 0x5be0cd191f83d9ab}
	input := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}

	// Create the struct instance
	params := SyscallSha256Params{State: &state, Input: &input}

	// Call the syscall function
	Sha256Update(&params)

	expected_state := [4]uint64{0x9e00af9da10d573d, 0xc8720af1ddd16faf, 0x5151dc0044476aa6, 0x535ad14d40420688}

	// If the endianess in helpers.rs wasn't changed then the expected_state would be
	// expected_state := [4]uint64{0x51f0168a9933c35b, 0x2e85fff39b6ad4d1, 0x5299d77714321264, 0xdb0bd60059557c70}

	if state != expected_state {
		t.Errorf("Failure")
	}
}
