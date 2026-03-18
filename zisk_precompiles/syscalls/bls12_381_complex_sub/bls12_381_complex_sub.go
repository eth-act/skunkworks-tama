//go:build tamago && riscv64

package bls12_381_complex_sub

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bls12_381_complex_sub.rs

type Complex struct {
	X [6]uint64
	Y [6]uint64
}

type SyscallBls12381ComplexSubParams struct {
	F1 *Complex
	F2 *Complex
}

//go:noinline
func Bls12381ComplexSub(arr *SyscallBls12381ComplexSubParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_complex_sub(ptr)
}

func syscall_bls12_381_complex_sub(ptr uintptr)
