package fp2

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/bls12_381/fp2_inv.rs

func TestInvOne(t *testing.T) {
	input := [12]uint64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expected := [12]uint64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestInv(t *testing.T) {
	input := [12]uint64{
		0x49b4b9e2ffd3bf5a, 0x6bc7632c9e4047a7, 0x805d19211a7dc450, 0x41c84ac8cfa40667,
		0xcbc8271a6d95e07f, 0x167ed52ad9b8dc52, 0x9919e620d143515b, 0x808a3f274c49a6c7,
		0xd65c110346cb2c1b, 0x8cd2c11ad5206061, 0x791b9ace70502ab1, 0x7f958516727acdd,
	}
	expected := [12]uint64{
		0x55aa5b187f77e83e, 0xff523f3ab3ac46a6, 0xf686d520afbeb578, 0xb1664497d371019b,
		0xcfef6ce72c61e835, 0x1474b2da727c6dfe, 0x5730d5d619884057, 0xd42b3decc96db687,
		0x8abb9a0eed22a8a3, 0xd2f92c46b24958f7, 0x8ab323bd7384ca05, 0x1859d94eddac5b45,
	}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/bls12_381/fp2_sqrt.rs

func TestFp2SqrtOne(t *testing.T) {
	input := Fp2Element{
		C0: [6]uint64{1, 0, 0, 0, 0, 0},
		C1: [6]uint64{0, 0, 0, 0, 0, 0},
	}

	expectedValue := Fp2Element{
		C0: [6]uint64{
			0xb9feffffffffaaaa,
			0x1eabfffeb153ffff,
			0x6730d2a0f6b0f624,
			0x64774b84f38512bf,
			0x4b1ba7b6434bacd7,
			0x1a0111ea397fe69a,
		},
		C1: [6]uint64{0, 0, 0, 0, 0, 0},
	}

	result := FcallSqrt(&input)

	if !result.Exists || result.Value != expectedValue {
		t.Errorf("got exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}
}

func TestFp2SqrtExists(t *testing.T) {
	input := Fp2Element{
		C0: [6]uint64{
			0x10486089be1876e9,
			0xcf0c3012bf0c13ef,
			0x51621421d2c37a8d,
			0xd52db71259449a47,
			0x370fd7a0a4be29da,
			0xc3d4fd75c076215,
		},
		C1: [6]uint64{
			0x3e6ff1a3151b0959,
			0x9f0b2a8dea2c9f82,
			0xb83d47ccb71501e2,
			0xa8c917818d857f05,
			0xc48150d1cd95e0c6,
			0x112ca78116187cc8,
		},
	}

	expectedValue := Fp2Element{
		C0: [6]uint64{
			0xcca66dfc0d7f69c9,
			0xaf22cf40d2f4555,
			0x92a6870798aff4d7,
			0xe595438fb87ee1fc,
			0x6f5e96c633b39798,
			0x215675032da3de5,
		},
		C1: [6]uint64{
			0x1ef8b538e151e6f3,
			0x94b37a0021182ef6,
			0xea0d1db797288ba2,
			0x567c72d5af34be56,
			0x5470d2ed597db716,
			0x10b61243878d0170,
		},
	}

	result := FcallSqrt(&input)

	if !result.Exists || result.Value != expectedValue {
		t.Errorf("got exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}
}

func TestFp2SqrtNotExists(t *testing.T) {
	input := Fp2Element{
		C0: [6]uint64{
			0x5531f66e0c366bf8,
			0x35f8f154ff2974e6,
			0xaa81eb7e92ae7b5e,
			0x8a521c9ff4654bc0,
			0xa224f0e84356bba8,
			0xffbbc4bdd5425cb,
		},
		C1: [6]uint64{
			0xf16972261c97a569,
			0xbf071b2a52d05a68,
			0xbaa99b2bc5260f74,
			0xedbd0c20e26eb5e5,
			0x6f3229e291d1d67a,
			0x119353ab08784f06,
		},
	}

	expectedValue := Fp2Element{
		C0: [6]uint64{
			0x6d8e1fc1edb82644,
			0xa6964afc770dab5d,
			0x37d90a0e925a572d,
			0x3547fbc3f051b409,
			0xd3cdef010df23067,
			0x159b8fd2cca0a180,
		},
		C1: [6]uint64{
			0xe0c163a5a7441092,
			0xf61c7202d7c3af80,
			0xf80c7aa929cb1e62,
			0xa076467c356a64cf,
			0x695e3d70b6a86704,
			0xb1ecd8ecdb0e8d2,
		},
	}

	result := FcallSqrt(&input)

	if result.Exists || result.Value != expectedValue {
		t.Errorf("got exists=%v value=%v, want exists=false value=%v", result.Exists, result.Value, expectedValue)
	}
}
