//go:build tamago && riscv64

package fp2

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/bls12_381/fp2.rs

type Fp2Element struct {
	C0 [6]uint64
	C1 [6]uint64
}

type Fp2SqrtResult struct {
	Exists bool
	Value  Fp2Element
}

//go:noinline
func FcallInv(input *[12]uint64) [12]uint64 {
	var result [12]uint64
	ptr := uintptr(unsafe.Pointer(input))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_inv(ptr, resultPtr)
	return result
}

//go:noinline
func FcallSqrt(element *Fp2Element) Fp2SqrtResult {
	var padded [16]uint64
	copy(padded[:6], element.C0[:])
	copy(padded[6:12], element.C1[:])
	var rawResult [13]uint64
	ptr := uintptr(unsafe.Pointer(&padded[0]))
	resultPtr := uintptr(unsafe.Pointer(&rawResult[0]))
	fcall_sqrt(ptr, resultPtr)

	var result Fp2SqrtResult
	result.Exists = rawResult[0] != 0
	copy(result.Value.C0[:], rawResult[1:7])
	copy(result.Value.C1[:], rawResult[7:13])
	return result
}

func fcall_inv(input, result uintptr)
func fcall_sqrt(element uintptr, result uintptr)
