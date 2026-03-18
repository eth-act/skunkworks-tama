package big_int_div

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestBigIntDivSimple(t *testing.T) {
	a := []uint64{10}
	b := []uint64{3}
	result := FcallBigIntDiv(a, b)

	expectedQuo := []uint64{3, 0, 0, 0}
	expectedRem := []uint64{1, 0, 0, 0}

	if !sliceEqual(result.Quotient, expectedQuo) {
		t.Errorf("quotient: got %v, want %v", result.Quotient, expectedQuo)
	}
	if !sliceEqual(result.Remainder, expectedRem) {
		t.Errorf("remainder: got %v, want %v", result.Remainder, expectedRem)
	}
}

func TestBigIntDivMultiLimb(t *testing.T) {
	a := []uint64{0x16b12176aedd308e, 0x9d331c2b34766fc9, 0xb7f85b22001249e, 0x3b4e3fc5e0d8b014}
	b := []uint64{0x16b12176aedd308e, 0x9d331c2b34766fc9}
	result := FcallBigIntDiv(a, b)

	expectedQuo := []uint64{0xcc3c538a93632e1a, 0x6094542194c94390, 0x0, 0x0}
	expectedRem := []uint64{0x5a1f8f05dcc2be22, 0x309a5adb5970397d, 0x0, 0x0}

	if !sliceEqual(result.Quotient, expectedQuo) {
		t.Errorf("quotient: got %v, want %v", result.Quotient, expectedQuo)
	}
	if !sliceEqual(result.Remainder, expectedRem) {
		t.Errorf("remainder: got %v, want %v", result.Remainder, expectedRem)
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
