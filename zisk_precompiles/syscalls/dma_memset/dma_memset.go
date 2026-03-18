//go:build tamago && riscv64

package dma_memset

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/precompiles/dma/src/dma/dma_memset.rs

//go:noinline
func DmaMemsetZero(dst, size uintptr) {
	dma_memset_zero(dst, size)
}

func dma_memset_zero(dst, size uintptr)
