//go:build tamago && riscv64

package arith256

import (
	"unsafe"
)

type SyscallArith256Params struct {
	A  *[4]uint64
	B  *[4]uint64
	C  *[4]uint64
	Dl *[4]uint64
	Dh *[4]uint64
}

//go:noinline
func Arith256_mod(arr *SyscallArith256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_arith256_mod(ptr)
}

func syscall_arith256_mod(ptr uintptr)
