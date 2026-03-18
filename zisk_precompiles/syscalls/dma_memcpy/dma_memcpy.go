//go:build tamago && riscv64

package dma_memcpy

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/precompiles/dma/src/dma/dma_memcpy.rs

//go:noinline
func DmaMemcpy(dst, src, size uintptr) {
	dma_memcpy(dst, src, size)
}

func dma_memcpy(dst, src, size uintptr)
