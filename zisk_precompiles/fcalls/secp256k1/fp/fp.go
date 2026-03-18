//go:build tamago && riscv64

package fp

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/secp256k1/fp.rs

type FpSqrtResult struct {
	Exists bool
	Value  [4]uint64
}

//go:noinline
func FcallInv(input *[4]uint64) [4]uint64 {
	var result [4]uint64
	ptr := uintptr(unsafe.Pointer(input))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_inv(ptr, resultPtr)
	return result
}

//go:noinline
func FcallSqrt(input *[4]uint64, parity uint64) FpSqrtResult {
	var rawResult [5]uint64
	ptr := uintptr(unsafe.Pointer(input))
	resultPtr := uintptr(unsafe.Pointer(&rawResult[0]))
	fcall_sqrt(ptr, uintptr(parity), resultPtr)
	var result FpSqrtResult
	result.Exists = rawResult[0] != 0
	copy(result.Value[:], rawResult[1:5])
	return result
}

func fcall_inv(input, result uintptr)
func fcall_sqrt(input, parity, result uintptr)
