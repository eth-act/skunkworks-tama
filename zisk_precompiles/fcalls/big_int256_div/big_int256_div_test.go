package big_int256_div

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Test vectors from: https://github.com/0xPolygonHermez/zisk/blob/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls_impl/big_int256_div.rs

func TestDiv(t *testing.T) {
	a := [4]uint64{0x16b12176aedd308e, 0x9d331c2b34766fc9, 0xb7f85b22001249e, 0x3b4e3fc5e0d8b014}
	b := [4]uint64{0x16b12176aedd308e, 0x9d331c2b34766fc9, 0xb7f85b22001249e, 0x0}
	expectedQuo := [4]uint64{0x2868ebf5edfaecd5, 0x5, 0x0, 0x0}
	expectedRem := [4]uint64{0xdbb84a86764e268, 0xfd48d6ec2b636246, 0xadbb6db4207ffb8, 0x0}

	result := FcallBigInt256Div(&a, &b)
	if result.Quotient != expectedQuo {
		t.Errorf("quotient: got %v, want %v", result.Quotient, expectedQuo)
	}
	if result.Remainder != expectedRem {
		t.Errorf("remainder: got %v, want %v", result.Remainder, expectedRem)
	}
}
