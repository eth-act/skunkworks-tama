//go:build tamago && riscv64

package bls12_381_complex_mul

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bls12_381_complex_mul.rs

type Complex struct {
	X [6]uint64
	Y [6]uint64
}

type SyscallBls12381ComplexMulParams struct {
	F1 *Complex
	F2 *Complex
}

//go:noinline
func Bls12381ComplexMul(arr *SyscallBls12381ComplexMulParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_complex_mul(ptr)
}

func syscall_bls12_381_complex_mul(ptr uintptr)
