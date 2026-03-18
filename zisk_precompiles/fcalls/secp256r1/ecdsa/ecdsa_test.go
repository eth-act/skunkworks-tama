package ecdsa

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/secp256r1/ecdsa.rs

func TestSecp256r1EcdsaVerify(t *testing.T) {
	params := EcdsaVerifyParams{
		Pk: [8]uint64{
			0xF4A13945D898C296, 0x77037D812DEB33A0, 0xF8BCE6E563A440F2, 0x6B17D1F2E12C4247,
			0xCBB6406837BF51F5, 0x2BCE33576B315ECE, 0x8EE7EB4A7C0F9E16, 0x4FE342E2FE1A7F9B,
		},
		Z: [4]uint64{2, 0, 0, 0},
		R: [4]uint64{3, 0, 0, 0},
		S: [4]uint64{1, 0, 0, 0},
	}

	expected := Point256{
		X: [4]uint64{0x21554a0dc3d033ed, 0xef8c82fd1f5be524, 0xd784c85608668fdf, 0x51590b7a515140d2},
		Y: [4]uint64{0xd1d0bb44fda16da4, 0x0d012f00d4d80888, 0x8ae1bf36bf8a7926, 0xe0c17da8904a727d},
	}

	result := FcallVerify(&params)

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
