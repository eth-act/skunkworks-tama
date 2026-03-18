package keccakf

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/precompiles/helpers/src/keccak/keccak_f/mod.rs

func TestKeccakZeroState(t *testing.T) {
	var data [25]uint64
	Keccak(&data)

	expected := [25]uint64{
		0xF1258F7940E1DDE7, 0x84D5CCF933C0478A, 0xD598261EA65AA9EE,
		0xBD1547306F80494D, 0x8B284E056253D057, 0xFF97A42D7F8E6FD4,
		0x90FEE5A0A44647C4, 0x8C5BDA0CD6192E76, 0xAD30A6F71B19059C,
		0x30935AB7D08FFC64, 0xEB5AA93F2317D635, 0xA9A6E6260D712103,
		0x81A57C16DBCF555F, 0x43B831CD0347C826, 0x01F22F1A11A5569F,
		0x05E5635A21D9AE61, 0x64BEFEF28CC970F2, 0x613670957BC46611,
		0xB87C5A554FD00ECB, 0x8C3EE88A1CCF32C8, 0x940C7922AE3A2614,
		0x1841F924A2C509E4, 0x16F53526E70465C2, 0x75F644E97F30A13B,
		0xEAF1FF7B5CECA249,
	}
	if data != expected {
		t.Errorf("got %v, want %v", data, expected)
	}
}

func TestKeccakNonzeroState(t *testing.T) {
	data := [25]uint64{
		0xF1258F7940E1DDE7, 0x84D5CCF933C0478A, 0xD598261EA65AA9EE,
		0xBD1547306F80494D, 0x8B284E056253D057, 0xFF97A42D7F8E6FD4,
		0x90FEE5A0A44647C4, 0x8C5BDA0CD6192E76, 0xAD30A6F71B19059C,
		0x30935AB7D08FFC64, 0xEB5AA93F2317D635, 0xA9A6E6260D712103,
		0x81A57C16DBCF555F, 0x43B831CD0347C826, 0x01F22F1A11A5569F,
		0x05E5635A21D9AE61, 0x64BEFEF28CC970F2, 0x613670957BC46611,
		0xB87C5A554FD00ECB, 0x8C3EE88A1CCF32C8, 0x940C7922AE3A2614,
		0x1841F924A2C509E4, 0x16F53526E70465C2, 0x75F644E97F30A13B,
		0xEAF1FF7B5CECA249,
	}
	Keccak(&data)

	expected := [25]uint64{
		0x2D5C954DF96ECB3C, 0x6A332CD07057B56D, 0x093D8D1270D76B6C,
		0x8A20D9B25569D094, 0x4F9C4F99E5E7F156, 0xF957B9A2DA65FB38,
		0x85773DAE1275AF0D, 0xFAF4F247C3D810F7, 0x1F1B9EE6F79A8759,
		0xE4FECC0FEE98B425, 0x68CE61B6B9CE68A1, 0xDEEA66C4BA8F974F,
		0x33C43D836EAFB1F5, 0xE00654042719DBD9, 0x7CF8A9F009831265,
		0xFD5449A6BF174743, 0x97DDAD33D8994B40, 0x48EAD5FC5D0BE774,
		0xE3B8C8EE55B7B03C, 0x91A0226E649E42E9, 0x900E3129E7BADD7B,
		0x202A9EC5FAA3CCE8, 0x5B3402464E1C3DB6, 0x609F4E62A44C1059,
		0x20D06CD26A8FBF5C,
	}
	if data != expected {
		t.Errorf("got %v, want %v", data, expected)
	}
}

// Reference implementation from: https://github.com/ethereum/go-ethereum/blob/v1.8.20/crypto/sha3/keccakf.go

var rc = [24]uint64{
	0x0000000000000001,
	0x0000000000008082,
	0x800000000000808A,
	0x8000000080008000,
	0x000000000000808B,
	0x0000000080000001,
	0x8000000080008081,
	0x8000000000008009,
	0x000000000000008A,
	0x0000000000000088,
	0x0000000080008009,
	0x000000008000000A,
	0x000000008000808B,
	0x800000000000008B,
	0x8000000000008089,
	0x8000000000008003,
	0x8000000000008002,
	0x8000000000000080,
	0x000000000000800A,
	0x800000008000000A,
	0x8000000080008081,
	0x8000000000008080,
	0x0000000080000001,
	0x8000000080008008,
}

func keccakF1600(a *[25]uint64) {
	var t, bc0, bc1, bc2, bc3, bc4, d0, d1, d2, d3, d4 uint64

	for i := 0; i < 24; i += 4 {
		// Round 1
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[6] ^ d1
		bc1 = t<<44 | t>>(64-44)
		t = a[12] ^ d2
		bc2 = t<<43 | t>>(64-43)
		t = a[18] ^ d3
		bc3 = t<<21 | t>>(64-21)
		t = a[24] ^ d4
		bc4 = t<<14 | t>>(64-14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i]
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc2 = t<<3 | t>>(64-3)
		t = a[16] ^ d1
		bc3 = t<<45 | t>>(64-45)
		t = a[22] ^ d2
		bc4 = t<<61 | t>>(64-61)
		t = a[3] ^ d3
		bc0 = t<<28 | t>>(64-28)
		t = a[9] ^ d4
		bc1 = t<<20 | t>>(64-20)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc4 = t<<18 | t>>(64-18)
		t = a[1] ^ d1
		bc0 = t<<1 | t>>(64-1)
		t = a[7] ^ d2
		bc1 = t<<6 | t>>(64-6)
		t = a[13] ^ d3
		bc2 = t<<25 | t>>(64-25)
		t = a[19] ^ d4
		bc3 = t<<8 | t>>(64-8)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc1 = t<<36 | t>>(64-36)
		t = a[11] ^ d1
		bc2 = t<<10 | t>>(64-10)
		t = a[17] ^ d2
		bc3 = t<<15 | t>>(64-15)
		t = a[23] ^ d3
		bc4 = t<<56 | t>>(64-56)
		t = a[4] ^ d4
		bc0 = t<<27 | t>>(64-27)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc3 = t<<41 | t>>(64-41)
		t = a[21] ^ d1
		bc4 = t<<2 | t>>(64-2)
		t = a[2] ^ d2
		bc0 = t<<62 | t>>(64-62)
		t = a[8] ^ d3
		bc1 = t<<55 | t>>(64-55)
		t = a[14] ^ d4
		bc2 = t<<39 | t>>(64-39)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		// Round 2
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[16] ^ d1
		bc1 = t<<44 | t>>(64-44)
		t = a[7] ^ d2
		bc2 = t<<43 | t>>(64-43)
		t = a[23] ^ d3
		bc3 = t<<21 | t>>(64-21)
		t = a[14] ^ d4
		bc4 = t<<14 | t>>(64-14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i+1]
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc2 = t<<3 | t>>(64-3)
		t = a[11] ^ d1
		bc3 = t<<45 | t>>(64-45)
		t = a[2] ^ d2
		bc4 = t<<61 | t>>(64-61)
		t = a[18] ^ d3
		bc0 = t<<28 | t>>(64-28)
		t = a[9] ^ d4
		bc1 = t<<20 | t>>(64-20)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc4 = t<<18 | t>>(64-18)
		t = a[6] ^ d1
		bc0 = t<<1 | t>>(64-1)
		t = a[22] ^ d2
		bc1 = t<<6 | t>>(64-6)
		t = a[13] ^ d3
		bc2 = t<<25 | t>>(64-25)
		t = a[4] ^ d4
		bc3 = t<<8 | t>>(64-8)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc1 = t<<36 | t>>(64-36)
		t = a[1] ^ d1
		bc2 = t<<10 | t>>(64-10)
		t = a[17] ^ d2
		bc3 = t<<15 | t>>(64-15)
		t = a[8] ^ d3
		bc4 = t<<56 | t>>(64-56)
		t = a[24] ^ d4
		bc0 = t<<27 | t>>(64-27)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc3 = t<<41 | t>>(64-41)
		t = a[21] ^ d1
		bc4 = t<<2 | t>>(64-2)
		t = a[12] ^ d2
		bc0 = t<<62 | t>>(64-62)
		t = a[3] ^ d3
		bc1 = t<<55 | t>>(64-55)
		t = a[19] ^ d4
		bc2 = t<<39 | t>>(64-39)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		// Round 3
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[11] ^ d1
		bc1 = t<<44 | t>>(64-44)
		t = a[22] ^ d2
		bc2 = t<<43 | t>>(64-43)
		t = a[8] ^ d3
		bc3 = t<<21 | t>>(64-21)
		t = a[19] ^ d4
		bc4 = t<<14 | t>>(64-14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i+2]
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc2 = t<<3 | t>>(64-3)
		t = a[1] ^ d1
		bc3 = t<<45 | t>>(64-45)
		t = a[12] ^ d2
		bc4 = t<<61 | t>>(64-61)
		t = a[23] ^ d3
		bc0 = t<<28 | t>>(64-28)
		t = a[9] ^ d4
		bc1 = t<<20 | t>>(64-20)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc4 = t<<18 | t>>(64-18)
		t = a[16] ^ d1
		bc0 = t<<1 | t>>(64-1)
		t = a[2] ^ d2
		bc1 = t<<6 | t>>(64-6)
		t = a[13] ^ d3
		bc2 = t<<25 | t>>(64-25)
		t = a[24] ^ d4
		bc3 = t<<8 | t>>(64-8)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc1 = t<<36 | t>>(64-36)
		t = a[6] ^ d1
		bc2 = t<<10 | t>>(64-10)
		t = a[17] ^ d2
		bc3 = t<<15 | t>>(64-15)
		t = a[3] ^ d3
		bc4 = t<<56 | t>>(64-56)
		t = a[14] ^ d4
		bc0 = t<<27 | t>>(64-27)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc3 = t<<41 | t>>(64-41)
		t = a[21] ^ d1
		bc4 = t<<2 | t>>(64-2)
		t = a[7] ^ d2
		bc0 = t<<62 | t>>(64-62)
		t = a[18] ^ d3
		bc1 = t<<55 | t>>(64-55)
		t = a[4] ^ d4
		bc2 = t<<39 | t>>(64-39)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		// Round 4
		bc0 = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		bc1 = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		bc2 = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		bc3 = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		bc4 = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]
		d0 = bc4 ^ (bc1<<1 | bc1>>63)
		d1 = bc0 ^ (bc2<<1 | bc2>>63)
		d2 = bc1 ^ (bc3<<1 | bc3>>63)
		d3 = bc2 ^ (bc4<<1 | bc4>>63)
		d4 = bc3 ^ (bc0<<1 | bc0>>63)

		bc0 = a[0] ^ d0
		t = a[1] ^ d1
		bc1 = t<<44 | t>>(64-44)
		t = a[2] ^ d2
		bc2 = t<<43 | t>>(64-43)
		t = a[3] ^ d3
		bc3 = t<<21 | t>>(64-21)
		t = a[4] ^ d4
		bc4 = t<<14 | t>>(64-14)
		a[0] = bc0 ^ (bc2 &^ bc1) ^ rc[i+3]
		a[1] = bc1 ^ (bc3 &^ bc2)
		a[2] = bc2 ^ (bc4 &^ bc3)
		a[3] = bc3 ^ (bc0 &^ bc4)
		a[4] = bc4 ^ (bc1 &^ bc0)

		t = a[5] ^ d0
		bc2 = t<<3 | t>>(64-3)
		t = a[6] ^ d1
		bc3 = t<<45 | t>>(64-45)
		t = a[7] ^ d2
		bc4 = t<<61 | t>>(64-61)
		t = a[8] ^ d3
		bc0 = t<<28 | t>>(64-28)
		t = a[9] ^ d4
		bc1 = t<<20 | t>>(64-20)
		a[5] = bc0 ^ (bc2 &^ bc1)
		a[6] = bc1 ^ (bc3 &^ bc2)
		a[7] = bc2 ^ (bc4 &^ bc3)
		a[8] = bc3 ^ (bc0 &^ bc4)
		a[9] = bc4 ^ (bc1 &^ bc0)

		t = a[10] ^ d0
		bc4 = t<<18 | t>>(64-18)
		t = a[11] ^ d1
		bc0 = t<<1 | t>>(64-1)
		t = a[12] ^ d2
		bc1 = t<<6 | t>>(64-6)
		t = a[13] ^ d3
		bc2 = t<<25 | t>>(64-25)
		t = a[14] ^ d4
		bc3 = t<<8 | t>>(64-8)
		a[10] = bc0 ^ (bc2 &^ bc1)
		a[11] = bc1 ^ (bc3 &^ bc2)
		a[12] = bc2 ^ (bc4 &^ bc3)
		a[13] = bc3 ^ (bc0 &^ bc4)
		a[14] = bc4 ^ (bc1 &^ bc0)

		t = a[15] ^ d0
		bc1 = t<<36 | t>>(64-36)
		t = a[16] ^ d1
		bc2 = t<<10 | t>>(64-10)
		t = a[17] ^ d2
		bc3 = t<<15 | t>>(64-15)
		t = a[18] ^ d3
		bc4 = t<<56 | t>>(64-56)
		t = a[19] ^ d4
		bc0 = t<<27 | t>>(64-27)
		a[15] = bc0 ^ (bc2 &^ bc1)
		a[16] = bc1 ^ (bc3 &^ bc2)
		a[17] = bc2 ^ (bc4 &^ bc3)
		a[18] = bc3 ^ (bc0 &^ bc4)
		a[19] = bc4 ^ (bc1 &^ bc0)

		t = a[20] ^ d0
		bc3 = t<<41 | t>>(64-41)
		t = a[21] ^ d1
		bc4 = t<<2 | t>>(64-2)
		t = a[22] ^ d2
		bc0 = t<<62 | t>>(64-62)
		t = a[23] ^ d3
		bc1 = t<<55 | t>>(64-55)
		t = a[24] ^ d4
		bc2 = t<<39 | t>>(64-39)
		a[20] = bc0 ^ (bc2 &^ bc1)
		a[21] = bc1 ^ (bc3 &^ bc2)
		a[22] = bc2 ^ (bc4 &^ bc3)
		a[23] = bc3 ^ (bc0 &^ bc4)
		a[24] = bc4 ^ (bc1 &^ bc0)
	}
}

func TestKeccakReferenceComparison(t *testing.T) {
	var data [25]uint64

	for i := range data {
		data[i] = uint64(i * 100)
	}

	copied := data

	Keccak(&data)
	keccakF1600(&copied)

	if data != copied {
		t.Errorf("got %v, want %v", data, copied)
	}
}
