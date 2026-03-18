//go:build tamago && riscv64

package blake2br

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/blake2br.rs

type SyscallBlake2bRoundParams struct {
	Index uint64
	State *[16]uint64
	Input *[16]uint64
}

//go:noinline
func Blake2bRound(params *SyscallBlake2bRoundParams) {
	ptr := uintptr(unsafe.Pointer(params))
	syscall_blake2br(ptr)
}

func syscall_blake2br(ptr uintptr)
