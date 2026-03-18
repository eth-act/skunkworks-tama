package fp

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/bls12_381/fp_inv.rs

func TestInvOne(t *testing.T) {
	input := [6]uint64{1, 0, 0, 0, 0, 0}
	expected := [6]uint64{1, 0, 0, 0, 0, 0}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestInv(t *testing.T) {
	input := [6]uint64{
		0x2d5f30c1d0577c56, 0x29aabf4bbbb4b60a, 0xf65faa3d6bda5044,
		0xa56da205ae4bf114, 0x6ad30a8453e66eac, 0x10a97e50d00668c,
	}
	expected := [6]uint64{
		0x1d8053f2aed3d017, 0x2912c6d8d7c59be0, 0xea3af967ab741430,
		0xdc3cb17c3b332919, 0x52a4afd74a0b5b20, 0x12be47b0938a6ee1,
	}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/bls12_381/fp_sqrt.rs

func TestSqrtOne(t *testing.T) {
	input := [6]uint64{1, 0, 0, 0, 0, 0}
	expectedValue := [6]uint64{1, 0, 0, 0, 0, 0}
	result := FcallSqrt(&input)
	if !result.Exists || result.Value != expectedValue {
		t.Errorf("exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}
}

func TestSqrt(t *testing.T) {
	input := [6]uint64{
		0xf22cb1516a067d13, 0x3e46be206ab02de6, 0x93153c30d0917c98,
		0x597d68ca77b5fa6d, 0x44a50733df914e5e, 0xf7377b1bb431d82,
	}
	expectedValue := [6]uint64{
		0x516e9b68ec7e4040, 0x4b1f0de82104d372, 0x7e742e30000909d7,
		0x44051766a1553492, 0xe7043ea4bffc292f, 0x3efcb69d6bf0ce0,
	}
	result := FcallSqrt(&input)
	if !result.Exists || result.Value != expectedValue {
		t.Errorf("exists=%v value=%v, want exists=true value=%v", result.Exists, result.Value, expectedValue)
	}
}

func TestNoSqrt(t *testing.T) {
	input := [6]uint64{
		0x361799ccd540a764, 0xf606e6b453a13bd8, 0x8880bd6a4b0b963a,
		0x8c9a8b3ba67f6d02, 0x922d30923791c733, 0x1975e3ccd03944ca,
	}
	expectedValue := [6]uint64{
		0x5514d9e1a2faebf1, 0x391ed94dec028013, 0x5a8c79b17991fded,
		0x56207337f5f736d0, 0xc6f1181533cc4b6, 0xb1d40edb0c1fec0,
	}
	result := FcallSqrt(&input)
	if result.Exists || result.Value != expectedValue {
		t.Errorf("exists=%v value=%v, want exists=false value=%v", result.Exists, result.Value, expectedValue)
	}
}
