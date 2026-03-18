package fn

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/secp256k1/fn_inv.rs

func TestInvOne(t *testing.T) {
	input := [4]uint64{1, 0, 0, 0}
	expected := [4]uint64{1, 0, 0, 0}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestInv(t *testing.T) {
	input := [4]uint64{0xf9ee4256a589409f, 0xa21a3985f17502d0, 0xb3eb393d00dc480c, 0x142def02c537eced}
	expected := [4]uint64{0x32fe23e91aa741a1, 0x204b2da7afd93e75, 0x39b0bef6b00ec8b0, 0x7a0f1a7146326666}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
