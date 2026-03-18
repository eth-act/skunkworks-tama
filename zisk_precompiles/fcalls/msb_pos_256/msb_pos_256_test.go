package msb_pos_256

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
)

// Expected values captured from ziskemu v0.16.0

func TestMsbPos256HighLimb(t *testing.T) {
	x := [4]uint64{0, 0, 0, 1}
	y := [4]uint64{1, 0, 0, 0}
	idx, pos := FcallMsbPos256(&x, &y)
	if idx != 3 || pos != 0 {
		t.Errorf("got idx=%d pos=%d, want idx=3 pos=0", idx, pos)
	}
}

func TestMsbPos256MidLimb(t *testing.T) {
	x := [4]uint64{0, 0, 0x8000000000000000, 0}
	y := [4]uint64{0, 0x100, 0, 0}
	idx, pos := FcallMsbPos256(&x, &y)
	if idx != 2 || pos != 63 {
		t.Errorf("got idx=%d pos=%d, want idx=2 pos=63", idx, pos)
	}
}

func TestMsbPos256YLarger(t *testing.T) {
	x := [4]uint64{0x80, 0, 0, 0}
	y := [4]uint64{0, 0, 0, 0x8000000000000000}
	idx, pos := FcallMsbPos256(&x, &y)
	if idx != 3 || pos != 63 {
		t.Errorf("got idx=%d pos=%d, want idx=3 pos=63", idx, pos)
	}
}
