package dma_memcmp

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
	"unsafe"
)

// Expected values captured from ziskemu v0.16.0

func TestDmaMemcmpEqual(t *testing.T) {
	a := [4]uint64{1, 2, 3, 4}
	b := [4]uint64{1, 2, 3, 4}

	result := DmaMemcmp(
		uintptr(unsafe.Pointer(&a[0])),
		uintptr(unsafe.Pointer(&b[0])),
		32,
	)

	if result != 0 {
		t.Errorf("got %d, want 0", result)
	}
}

func TestDmaMemcmpNotEqual(t *testing.T) {
	a := [4]uint64{1, 2, 3, 4}
	b := [4]uint64{1, 2, 3, 5}

	result := DmaMemcmp(
		uintptr(unsafe.Pointer(&a[0])),
		uintptr(unsafe.Pointer(&b[0])),
		32,
	)

	if result == 0 {
		t.Errorf("got %d, want non-zero", result)
	}
}
