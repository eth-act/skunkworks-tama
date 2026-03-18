package bin_decomp

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestBinDecomp5(t *testing.T) {
	x := []uint64{5}
	bits := FcallBinDecomp(x)
	expected := []uint64{1, 0, 1}
	if !sliceEqual(bits, expected) {
		t.Errorf("got %v, want %v", bits, expected)
	}
}

func TestBinDecomp0xFF(t *testing.T) {
	x := []uint64{0xFF}
	bits := FcallBinDecomp(x)
	expected := []uint64{1, 1, 1, 1, 1, 1, 1, 1}
	if !sliceEqual(bits, expected) {
		t.Errorf("got %v, want %v", bits, expected)
	}
}

func TestBinDecompMultiLimb(t *testing.T) {
	x := []uint64{0, 1}
	bits := FcallBinDecomp(x)
	if len(bits) != 65 {
		t.Errorf("expected 65 bits, got %d", len(bits))
		return
	}
	if bits[0] != 1 {
		t.Errorf("expected MSB bits[0]=1, got %d", bits[0])
	}
	for i := 1; i < 65; i++ {
		if bits[i] != 0 {
			t.Errorf("expected bits[%d]=0, got %d", i, bits[i])
		}
	}
}

func sliceEqual(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
