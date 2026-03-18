package poseidon2

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestPoseidon2ZeroState(t *testing.T) {
	var state [16]uint64
	Poseidon2(&state)

	expected := [16]uint64{
		0xf2b2442ea4d72b98, 0x8367625af002a12, 0x41d794a3d56b9451, 0x533967a2f0a214c8,
		0x9b10cb9aecef64c2, 0x3af18efb76e71cc4, 0x20d42b106f3cd4d6, 0x537149275a93e1b9,
		0xe48c755b2541ac33, 0xd88485c5e6be8ad5, 0xf864699c52b2d651, 0x3bb13e057d4f33c6,
		0x7530b7e50b638c15, 0x4664c38414614b49, 0x267451ae2a8b9c47, 0x6a683ad447354817,
	}
	if state != expected {
		t.Errorf("got %v, want %v", state, expected)
	}
}

func TestPoseidon2Sequential(t *testing.T) {
	state := [16]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	Poseidon2(&state)

	expected := [16]uint64{
		0x85c54702470d9756, 0xaa53c7a7d52d9898, 0x285128096efb0dd7, 0xf3fde5edd3050ac8,
		0xc7b65efd040df908, 0x4be3f6c467f57ae9, 0x274e9a67b41754fb, 0xf7d39cd5de94dac,
		0xd0224b9794d0b78c, 0x372f6139570042e1, 0xce6e8a93dc4ec26c, 0xace65e30a4daf7af,
		0x16f2824cc1ba3db, 0x2e8f3af37c434dec, 0xc80831bb6e09da01, 0x3a7d670bf1a86ee8,
	}
	if state != expected {
		t.Errorf("got %v, want %v", state, expected)
	}
}
