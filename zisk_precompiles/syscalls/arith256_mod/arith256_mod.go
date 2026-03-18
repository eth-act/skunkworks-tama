//go:build tamago && riscv64

package arith256_mod

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/arith256_mod.rs

type SyscallArith256ModParams struct {
	A       *[4]uint64
	B       *[4]uint64
	C       *[4]uint64
	Modulus *[4]uint64
	D       *[4]uint64
}

//go:noinline
func Arith256Mod(arr *SyscallArith256ModParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_arith256_mod(ptr)
}

func syscall_arith256_mod(ptr uintptr)
