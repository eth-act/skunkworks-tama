//go:build tamago && riscv64

package dma_inputcpy

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/precompiles/dma/src/dma/dma_inputcpy.rs

//go:noinline
func DmaInputcpy(dst, size uintptr) {
	dma_inputcpy(dst, size)
}

func dma_inputcpy(dst, size uintptr)
