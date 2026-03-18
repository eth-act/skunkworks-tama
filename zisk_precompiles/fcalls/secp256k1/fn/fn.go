//go:build tamago && riscv64

package fn

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/secp256k1/fn.rs

//go:noinline
func FcallInv(input *[4]uint64) [4]uint64 {
	var result [4]uint64
	ptr := uintptr(unsafe.Pointer(input))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_inv(ptr, resultPtr)
	return result
}

func fcall_inv(input, result uintptr)
