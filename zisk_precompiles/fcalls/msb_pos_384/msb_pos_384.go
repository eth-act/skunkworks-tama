//go:build tamago && riscv64

package msb_pos_384

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/msb_pos_384.rs

//go:noinline
func FcallMsbPos384(x, y *[6]uint64) (uint64, uint64) {
	var xPadded, yPadded [8]uint64
	copy(xPadded[:6], x[:])
	copy(yPadded[:6], y[:])
	var result [2]uint64
	xPtr := uintptr(unsafe.Pointer(&xPadded[0]))
	yPtr := uintptr(unsafe.Pointer(&yPadded[0]))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_msb_pos_384(xPtr, yPtr, resultPtr)
	return result[0], result[1]
}

func fcall_msb_pos_384(x, y, result uintptr)
