package arith256_mod

import (
	"math/big"
	_ "tamagotest/tamaboards/zkvm"
	"tamagotest/zisk_precompiles/internal"
	"testing"
)

func verifyComputation(params *SyscallArith256ModParams, t *testing.T) {
	a := internal.Uint64ArrayToBigInt(params.A)
	b := internal.Uint64ArrayToBigInt(params.B)
	c := internal.Uint64ArrayToBigInt(params.C)
	module := internal.Uint64ArrayToBigInt(params.Modulus)
	d := internal.Uint64ArrayToBigInt(params.D)

	// Compute expected result: (a * b + c) mod module
	expected := new(big.Int)
	expected.Mul(a, b)             // a * b
	expected.Add(expected, c)      // + c
	expected.Mod(expected, module) // mod module

	match := expected.Cmp(d) == 0
	if !match {
		t.Errorf("Failure")
	}
}

func TestArith256(t *testing.T) {
	a := [4]uint64{0x112210F4B2D230D2, 0x0000000000000019, 0x0000000000000000, 0x0000000000000000}

	b := [4]uint64{0x84B2C7C5BC5C597A, 0x00000000000000CB, 0x0000000000000000, 0x0000000000000000}

	c := [4]uint64{0xC19E6C2252E377C7, 0x0000000000000017, 0x0000000000000000, 0x0000000000000000}

	// module = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
	// This is the secp256k1 field prime (used in Bitcoin/Ethereum)
	module := [4]uint64{
		0xFFFFFFFEFFFFFC2F, // least significant
		0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, // most significant
	}

	d := [4]uint64{0, 0, 0, 0}

	params := SyscallArith256ModParams{
		A:       &a,
		B:       &b,
		C:       &c,
		Modulus: &module,
		D:       &d,
	}

	Arith256_mod(&params)

	verifyComputation(&params, t)
}
