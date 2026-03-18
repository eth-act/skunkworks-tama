package fp

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/secp256k1/fp_inv.rs

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
	expected := [4]uint64{0xc198809f72408ac9, 0xa8726302e84e0c65, 0xde970a9a3b70d025, 0xf70d37bc0fece9b8}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/secp256k1/fp_sqrt.rs

func TestSqrtOne(t *testing.T) {
	input := [4]uint64{1, 0, 0, 0}

	result := FcallSqrt(&input, 1)
	expectedValue := [4]uint64{1, 0, 0, 0}
	if !result.Exists || result.Value != expectedValue {
		t.Errorf("parity=1: exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}

	result = FcallSqrt(&input, 0)
	expectedValue = [4]uint64{0xfffffffefffffc2e, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff}
	if !result.Exists || result.Value != expectedValue {
		t.Errorf("parity=0: exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}
}

func TestSqrt(t *testing.T) {
	input := [4]uint64{0x643764b2faa1592a, 0x4ac3ab52286f702a, 0x6591d88c833ffd4f, 0xc6fb7a1e514eac26}

	result := FcallSqrt(&input, 0)
	expectedValue := [4]uint64{0xa3d2fb0160f29df6, 0x3ebce4d565b52649, 0x4cdec0bf5c968639, 0x123e42087c415355}
	if !result.Exists || result.Value != expectedValue {
		t.Errorf("parity=0: exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}

	result = FcallSqrt(&input, 1)
	expectedValue = [4]uint64{0x5c2d04fd9f0d5e39, 0xc1431b2a9a4ad9b6, 0xb3213f40a36979c6, 0xedc1bdf783beacaa}
	if !result.Exists || result.Value != expectedValue {
		t.Errorf("parity=1: exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}
}

func TestNoSqrt(t *testing.T) {
	input := [4]uint64{0x643764b2faa1592c, 0x4ac3ab52286f702a, 0x6591d88c833ffd4f, 0xc6fb7a1e514eac26}
	expectedValue := [4]uint64{0xdab2978e63122590, 0x5dc785c971480237, 0x87a60df9f92b07b9, 0x855b365e9f83d30d}

	result := FcallSqrt(&input, 0)
	if result.Exists || result.Value != expectedValue {
		t.Errorf("parity=0: exists=%v value=%v, want exists=false value=%v", result.Exists, result.Value, expectedValue)
	}

	result = FcallSqrt(&input, 1)
	if result.Exists || result.Value != expectedValue {
		t.Errorf("parity=1: exists=%v value=%v, want exists=false value=%v", result.Exists, result.Value, expectedValue)
	}
}
