//go:build tamago && riscv64

package big_int_div

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/big_int_div.rs

type BigIntDivResult struct {
	Quotient  []uint64
	Remainder []uint64
}

//go:noinline
func FcallBigIntDiv(a, b []uint64) BigIntDivResult {
	fcallParam1(uint64(len(a)))
	for i := range a {
		fcallParam1(a[i])
	}

	fcallParam1(uint64(len(b)))
	for i := range b {
		fcallParam1(b[i])
	}

	fcallTrigger()

	lenQuo := fcallGet()
	quotient := make([]uint64, lenQuo)
	for i := uint64(0); i < lenQuo; i++ {
		quotient[i] = fcallGet()
	}

	lenRem := fcallGet()
	remainder := make([]uint64, lenRem)
	for i := uint64(0); i < lenRem; i++ {
		remainder[i] = fcallGet()
	}

	return BigIntDivResult{
		Quotient:  quotient,
		Remainder: remainder,
	}
}

func fcallParam1(val uint64)
func fcallGet() uint64
func fcallTrigger()
