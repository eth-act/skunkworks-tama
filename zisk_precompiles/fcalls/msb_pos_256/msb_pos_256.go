//go:build tamago && riscv64

package msb_pos_256

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/msb_pos_256.rs

//go:noinline
func FcallMsbPos256(x, y *[4]uint64) (uint64, uint64) {
	var result [2]uint64
	xPtr := uintptr(unsafe.Pointer(x))
	yPtr := uintptr(unsafe.Pointer(y))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_msb_pos_256(2, xPtr, yPtr, resultPtr)
	return result[0], result[1]
}

func fcall_msb_pos_256(count uint64, x, y, result uintptr)
