//go:build tamago && riscv64

package big_int256_div

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/big_int256_div.rs

type DivResult struct {
	Quotient  [4]uint64
	Remainder [4]uint64
}

//go:noinline
func FcallBigInt256Div(a, b *[4]uint64) DivResult {
	var result [8]uint64
	aPtr := uintptr(unsafe.Pointer(a))
	bPtr := uintptr(unsafe.Pointer(b))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_big_int256_div(aPtr, bPtr, resultPtr)
	var divResult DivResult
	copy(divResult.Quotient[:], result[0:4])
	copy(divResult.Remainder[:], result[4:8])
	return divResult
}

func fcall_big_int256_div(a, b, result uintptr)
