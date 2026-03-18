//go:build tamago && riscv64

package dma_memcmp

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/precompiles/dma/src/dma/dma_memcmp.rs

//go:noinline
func DmaMemcmp(a, b, size uintptr) int64 {
	return dma_memcmp(a, b, size)
}

func dma_memcmp(a, b, size uintptr) int64
