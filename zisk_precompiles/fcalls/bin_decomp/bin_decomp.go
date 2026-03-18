//go:build tamago && riscv64

package bin_decomp

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/bin_decomp.rs

//go:noinline
func FcallBinDecomp(x []uint64) []uint64 {
	fcallParam1(uint64(len(x)))
	for i := range x {
		fcallParam1(x[i])
	}

	fcallTrigger()

	lenBits := fcallGet()
	bits := make([]uint64, lenBits)
	for i := uint64(0); i < lenBits; i++ {
		bits[i] = fcallGet()
	}

	return bits
}

func fcallParam1(val uint64)
func fcallGet() uint64
func fcallTrigger()
