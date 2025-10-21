//go:build tamago && riscv64

package arith256_mod

import (
	"unsafe"
)

type SyscallArith256ModParams struct {
	A      *[4]uint64
	B      *[4]uint64
	C      *[4]uint64
	Module *[4]uint64
	D      *[4]uint64
}

//go:noinline
func Arith256_mod(arr *SyscallArith256ModParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_arith256_mod(ptr)
}

func syscall_arith256_mod(ptr uintptr)
