//go:build tamago && riscv64

package bn254_complex_add

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bn254_complex_add.rs

type Complex struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallBn254ComplexAddParams struct {
	F1 *Complex
	F2 *Complex
}

//go:noinline
func Bn254ComplexAdd(arr *SyscallBn254ComplexAddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_complex_add(ptr)
}

func syscall_bn254_complex_add(ptr uintptr)
