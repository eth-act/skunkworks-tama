//go:build tamago && riscv64

package arith384_mod

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/arith384_mod.rs

type SyscallArith384ModParams struct {
	A       *[6]uint64
	B       *[6]uint64
	C       *[6]uint64
	Modulus *[6]uint64
	D       *[6]uint64
}

//go:noinline
func Arith384Mod(arr *SyscallArith384ModParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_arith384_mod(ptr)
}

func syscall_arith384_mod(ptr uintptr)
