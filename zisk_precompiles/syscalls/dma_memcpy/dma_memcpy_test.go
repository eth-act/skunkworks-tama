package dma_memcpy

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
	"unsafe"
)

// Expected values captured from ziskemu v0.16.0

func TestDmaMemcpy(t *testing.T) {
	src := [4]uint64{0x1111111111111111, 0x2222222222222222, 0x3333333333333333, 0x4444444444444444}
	dst := [4]uint64{0, 0, 0, 0}

	DmaMemcpy(
		uintptr(unsafe.Pointer(&dst[0])),
		uintptr(unsafe.Pointer(&src[0])),
		32,
	)

	if dst != src {
		t.Errorf("got %v, want %v", dst, src)
	}
}
