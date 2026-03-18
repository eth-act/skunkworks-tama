package dma_memset

import (
	_ "github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"testing"
	"unsafe"
)

// Expected values captured from ziskemu v0.16.0

func TestDmaMemsetZero(t *testing.T) {
	data := [4]uint64{0xDEADBEEF, 0xCAFEBABE, 0x12345678, 0x9ABCDEF0}

	DmaMemsetZero(
		uintptr(unsafe.Pointer(&data[0])),
		32,
	)

	expected := [4]uint64{0, 0, 0, 0}
	if data != expected {
		t.Errorf("got %v, want %v", data, expected)
	}
}
