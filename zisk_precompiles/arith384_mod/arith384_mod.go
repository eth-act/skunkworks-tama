//go:build tamago && riscv64

package arith384_mod

import (
	"unsafe"
)

type SyscallArith384ModParams struct {
	A      *[6]uint64
	B      *[6]uint64
	C      *[6]uint64
	Module *[6]uint64
	D      *[6]uint64
}

//go:noinline
func Arith384_mod(arr *SyscallArith384ModParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_arith384_mod(ptr)
}

func syscall_arith384_mod(ptr uintptr)
