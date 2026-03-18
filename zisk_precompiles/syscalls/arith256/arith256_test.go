package arith256

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"github.com/eth-act/skunkworks-tama/zisk_precompiles/internal"
	"math/big"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/precompiles/arith_eq/src/test_data/arith256_test_data.rs
func TestArith256(t *testing.T) {
	tests := []struct {
		a, b, c, dh, dl string
	}{
		{"3", "2", "5", "0", "11"},
		{"256", "256", "1", "0", "65537"},
		{"3000", "2000", "5000", "0", "6005000"},
		{"3000000", "2000000", "5000000", "0", "6000005000000"},
		{"3000", "0", "5000", "0", "5000"},
	}
	for i, tt := range tests {
		a := internal.DecStringToUint64x4(tt.a)
		b := internal.DecStringToUint64x4(tt.b)
		c := internal.DecStringToUint64x4(tt.c)
		expectedDh := internal.DecStringToUint64x4(tt.dh)
		expectedDl := internal.DecStringToUint64x4(tt.dl)
		dh := [4]uint64{}
		dl := [4]uint64{}
		params := SyscallArith256Params{A: &a, B: &b, C: &c, Dh: &dh, Dl: &dl}
		Arith256(&params)
		if dl != expectedDl || dh != expectedDh {
			t.Errorf("vector %d: got dh=%v dl=%v, want dh=%v dl=%v", i, dh, dl, expectedDh, expectedDl)
		}
	}
}

func verifyComputation(params *SyscallArith256Params, t *testing.T) {
	a := internal.Uint64ArrayToBigInt(params.A)
	b := internal.Uint64ArrayToBigInt(params.B)
	c := internal.Uint64ArrayToBigInt(params.C)
	dl := internal.Uint64ArrayToBigInt(params.Dl)
	dh := internal.Uint64ArrayToBigInt(params.Dh)
	mask := internal.Uint64ArrayToBigInt(&[4]uint64{0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff, 0xffffffffffffffff})

	expected_dhdl := new(big.Int)
	expected_dhdl.Mul(a, b)
	expected_dhdl.Add(expected_dhdl, c)

	expected_dl := new(big.Int)
	expected_dl.And(expected_dhdl, mask)

	expected_dh := new(big.Int)
	expected_dh.Rsh(expected_dhdl, 256)
	expected_dh.And(expected_dh, mask)

	match := expected_dl.Cmp(dl) == 0 && expected_dh.Cmp(dh) == 0
	if !match {
		t.Errorf("got dl=%v dh=%v, want dl=%v dh=%v", dl, dh, expected_dl, expected_dh)
	}
}

func TestArith256Verify(t *testing.T) {
	a := [4]uint64{0x112210F4B2D230D2, 0x74172ed0ab34612c, 0x5d46fa993452fef9, 0x657aaf58cb17a53c}
	b := [4]uint64{0x84B2C7C5BC5C597A, 0x2e6371d577af32ab, 0x3c8f0bc543b1286c, 0x037f0aa5d1c48341}
	c := [4]uint64{0xC19E6C2252E377C7, 0xc864ca220f556c33, 0xae15c7d0cd7774ba, 0x112ddde737d94162}
	dh := [4]uint64{0, 0, 0, 0}
	dl := [4]uint64{0, 0, 0, 0}

	params := SyscallArith256Params{A: &a, B: &b, C: &c, Dh: &dh, Dl: &dl}
	Arith256(&params)
	verifyComputation(&params, t)
}
