//go:build tamago && riscv64

package arith256

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/arith256.rs

type SyscallArith256Params struct {
	A  *[4]uint64
	B  *[4]uint64
	C  *[4]uint64
	Dl *[4]uint64
	Dh *[4]uint64
}

//go:noinline
func Arith256(arr *SyscallArith256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_arith256(ptr)
}

func syscall_arith256(ptr uintptr)
