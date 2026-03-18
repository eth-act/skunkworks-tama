//go:build tamago && riscv64

package bls12_381_complex_add

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bls12_381_complex_add.rs

type Complex struct {
	X [6]uint64
	Y [6]uint64
}

type SyscallBls12381ComplexAddParams struct {
	F1 *Complex
	F2 *Complex
}

//go:noinline
func Bls12381ComplexAdd(arr *SyscallBls12381ComplexAddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_complex_add(ptr)
}

func syscall_bls12_381_complex_add(ptr uintptr)
