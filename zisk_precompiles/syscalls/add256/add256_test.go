package add256

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestAdd256Basic(t *testing.T) {
	a := [4]uint64{1, 0, 0, 0}
	b := [4]uint64{2, 0, 0, 0}
	var c [4]uint64

	params := SyscallAdd256Params{
		A:   &a,
		B:   &b,
		Cin: 0,
		C:   &c,
	}

	cout := Add256(&params)

	expectedC := [4]uint64{3, 0, 0, 0}
	if c != expectedC {
		t.Errorf("c: got %v, want %v", c, expectedC)
	}
	if cout != 0 {
		t.Errorf("cout: got %d, want 0", cout)
	}
}

func TestAdd256WithCarry(t *testing.T) {
	a := [4]uint64{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}
	b := [4]uint64{1, 0, 0, 0}
	var c [4]uint64

	params := SyscallAdd256Params{
		A:   &a,
		B:   &b,
		Cin: 0,
		C:   &c,
	}

	cout := Add256(&params)

	expectedC := [4]uint64{0, 0, 0, 0}
	if c != expectedC {
		t.Errorf("c: got %v, want %v", c, expectedC)
	}
	if cout != 1 {
		t.Errorf("cout: got %d, want 1", cout)
	}
}
