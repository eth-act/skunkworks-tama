//go:build tamago && riscv64

package ecdsa

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/zisklib/fcalls/secp256k1/ecdsa.rs

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

type EcdsaVerifyParams struct {
	Pk [8]uint64
	Z  [4]uint64
	R  [4]uint64
	S  [4]uint64
}

//go:noinline
func FcallVerify(params *EcdsaVerifyParams) Point256 {
	var result [8]uint64
	pkPtr := uintptr(unsafe.Pointer(&params.Pk[0]))
	zPtr := uintptr(unsafe.Pointer(&params.Z[0]))
	rPtr := uintptr(unsafe.Pointer(&params.R[0]))
	sPtr := uintptr(unsafe.Pointer(&params.S[0]))
	resultPtr := uintptr(unsafe.Pointer(&result[0]))
	fcall_verify(pkPtr, zPtr, rPtr, sPtr, resultPtr)

	var p Point256
	copy(p.X[:], result[0:4])
	copy(p.Y[:], result[4:8])
	return p
}

func fcall_verify(pk, z, r, s, result uintptr)
