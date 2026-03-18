package blake2br

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/precompiles/helpers/src/blake2/blake2b/mod.rs

func TestBlake2bRound(t *testing.T) {
	state := [16]uint64{
		0x6a09e667f3bcc908, 0xbb67ae8584caa73b,
		0x3c6ef372fe94f82b, 0xa54ff53a5f1d36f1,
		0x510e527fade682d1, 0x9b05688c2b3e6c1f,
		0x1f83d9abfb41bd6b, 0x5be0cd19137e2179,
		0x6a09e667f3bcc908, 0xbb67ae8584caa73b,
		0x3c6ef372fe94f82b, 0xa54ff53a5f1d36f1,
		0x510e527fade682d1, 0x9b05688c2b3e6c1f,
		0x1f83d9abfb41bd6b, 0x5be0cd19137e2179,
	}
	input := [16]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	params := SyscallBlake2bRoundParams{
		Index: 0,
		State: &state,
		Input: &input,
	}

	Blake2bRound(&params)

	expected := [16]uint64{
		0xd55b46925ed55b3a, 0x4eb04e43e11afb88,
		0x89d6f2a774fb03dd, 0x346cd307546da6db,
		0x622f3384b788b9aa, 0xb8ec4601a4587f7e,
		0xb0c8c2e4c07c3cf8, 0x528ae96b5b858171,
		0xede323999d36b7df, 0x55ba81615222d385,
		0x5f47fa7a7a33116e, 0xe42220840ad78c8c,
		0x0b495b872f7842a2, 0x60266c4cf9e8168e,
		0xbd8c2e25fff60a30, 0xc3501e54781c3097,
	}

	if state != expected {
		t.Errorf("got %v, want %v", state, expected)
	}
}
