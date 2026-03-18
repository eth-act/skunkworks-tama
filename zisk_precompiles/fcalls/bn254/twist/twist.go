//go:build tamago && riscv64

package twist

import "unsafe"

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/bn254/twist.rs

type TwistLineCoeffs struct {
	Lambda [8]uint64
	Mu     [8]uint64
}

//go:noinline
func FcallAddLineCoeffs(p1, p2 *[16]uint64) TwistLineCoeffs {
	var result [16]uint64
	p1Ptr := uintptr(unsafe.Pointer(p1))
	p2Ptr := uintptr(unsafe.Pointer(p2))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_add(p1Ptr, p2Ptr, resultPtr)
	var coeffs TwistLineCoeffs
	copy(coeffs.Lambda[:], result[0:8])
	copy(coeffs.Mu[:], result[8:16])
	return coeffs
}

//go:noinline
func FcallDblLineCoeffs(p *[16]uint64) TwistLineCoeffs {
	var result [16]uint64
	ptr := uintptr(unsafe.Pointer(p))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_dbl(ptr, resultPtr)
	var coeffs TwistLineCoeffs
	copy(coeffs.Lambda[:], result[0:8])
	copy(coeffs.Mu[:], result[8:16])
	return coeffs
}

func fcall_add(p1, p2, result uintptr)
func fcall_dbl(p, result uintptr)
