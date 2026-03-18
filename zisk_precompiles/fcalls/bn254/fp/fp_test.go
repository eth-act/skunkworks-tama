package fp

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/bn254/fp.rs

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
	expected := [4]uint64{0x7258dab6e90d1680, 0x779f7ec5cad25c1d, 0xb9c114d250bcaa3c, 0x2525db1f6832d97d}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
