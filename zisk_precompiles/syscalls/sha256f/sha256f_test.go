package sha256f

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestSha256(t *testing.T) {
	state := [4]uint64{0xbb67ae856a09e667, 0xa54ff53a3c6ef372, 0x9b05688c510e527f, 0x5be0cd191f83d9ab}
	input := [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}

	params := SyscallSha256Params{State: &state, Input: &input}
	Sha256Update(&params)

	expected := [4]uint64{0x51f0168a9933c35b, 0x2e85fff39b6ad4d1, 0x5299d77714321264, 0xdb0bd60059557c70}
	if state != expected {
		t.Errorf("got %v, want %v", state, expected)
	}
}
