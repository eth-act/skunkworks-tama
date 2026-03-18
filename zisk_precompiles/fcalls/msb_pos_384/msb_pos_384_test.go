package msb_pos_384

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestMsbPos384HighLimb(t *testing.T) {
	x := [6]uint64{0, 0, 0, 0, 0, 1}
	y := [6]uint64{1, 0, 0, 0, 0, 0}
	idx, pos := FcallMsbPos384(&x, &y)
	if idx != 5 || pos != 0 {
		t.Errorf("got idx=%d pos=%d, want idx=5 pos=0", idx, pos)
	}
}

func TestMsbPos384MidLimb(t *testing.T) {
	x := [6]uint64{0, 0, 0x8000000000000000, 0, 0, 0}
	y := [6]uint64{0, 0x100, 0, 0, 0, 0}
	idx, pos := FcallMsbPos384(&x, &y)
	if idx != 2 || pos != 63 {
		t.Errorf("got idx=%d pos=%d, want idx=2 pos=63", idx, pos)
	}
}
