package fp2

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/bn254/fp2.rs

func TestInvOne(t *testing.T) {
	input := [8]uint64{1, 0, 0, 0, 0, 0, 0, 0}
	expected := [8]uint64{1, 0, 0, 0, 0, 0, 0, 0}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestInv(t *testing.T) {
	input := [8]uint64{
		0xa4528921da9661b8, 0xc13514a2f09d4f06, 0x52406705a0d612b8, 0x2b02b26b72efef38,
		0xb64cd3ecb5b08b28, 0xe29c6143da89de45, 0xdfa4f8b46115f7f6, 0x17abb41fc8d1b2c7,
	}
	expected := [8]uint64{
		0x163d11f5aa617bfc, 0x825bc78934e518e5, 0x31485988143cff2e, 0x0551d3643b94a0ba,
		0xbd2738b4b0c67843, 0xbed5ac50b31d3cef, 0x516d2e7c293eef52, 0x302d79e76ed154c1,
	}
	result := FcallInv(&input)
	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
