//go:build tamago && riscv64

package fp

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/bls12_381/fp.rs

type FpSqrtResult struct {
	Exists bool
	Value  [6]uint64
}

//go:noinline
func FcallInv(input *[6]uint64) [6]uint64 {
	var padded [8]uint64
	copy(padded[:6], input[:])
	var result [6]uint64
	ptr := uintptr(unsafe.Pointer(&padded[0]))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_inv(ptr, resultPtr)
	return result
}

//go:noinline
func FcallSqrt(input *[6]uint64) FpSqrtResult {
	var padded [8]uint64
	copy(padded[:6], input[:])
	var rawResult [7]uint64
	ptr := uintptr(unsafe.Pointer(&padded[0]))
	resultPtr := uintptr(unsafe.Pointer(&rawResult[0]))
	fcall_sqrt(ptr, resultPtr)
	var result FpSqrtResult
	result.Exists = rawResult[0] != 0
	copy(result.Value[:], rawResult[1:7])
	return result
}

func fcall_inv(input, result uintptr)
func fcall_sqrt(input, result uintptr)
