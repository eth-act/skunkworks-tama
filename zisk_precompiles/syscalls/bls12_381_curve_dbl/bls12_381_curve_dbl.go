//go:build tamago && riscv64

package bls12_381_curve_dbl

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bls12_381_curve_dbl.rs

type Point384 struct {
	X [6]uint64
	Y [6]uint64
}

//go:noinline
func Bls12381Dbl(arr *Point384) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_dbl(ptr)
}

func syscall_bls12_381_dbl(ptr uintptr)
