package ecdsa

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/secp256k1/ecdsa.rs

func TestSecp256k1EcdsaVerify(t *testing.T) {
	params := EcdsaVerifyParams{
		Pk: [8]uint64{
			0x59F2815B16F81798, 0x029BFCDB2DCE28D9, 0x55A06295CE870B07, 0x79BE667EF9DCBBAC,
			0x9C47D08FFB10D4B8, 0xFD17B448A6855419, 0x5DA4FBFC0E1108A8, 0x483ADA7726A3C465,
		},
		Z: [4]uint64{2, 0, 0, 0},
		R: [4]uint64{3, 0, 0, 0},
		S: [4]uint64{1, 0, 0, 0},
	}

	expected := Point256{
		X: [4]uint64{0xcba8d569b240efe4, 0xe88b84bddc619ab7, 0x55b4a7250a5c5128, 0x2f8bde4d1a072093},
		Y: [4]uint64{0xdca87d3aa6ac62d6, 0xf788271bab0d6840, 0xd4dba9dda6c9c426, 0xd8ac222636e5e3d6},
	}

	result := FcallVerify(&params)

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
