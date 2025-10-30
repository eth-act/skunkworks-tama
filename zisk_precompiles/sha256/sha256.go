//go:build tamago && riscv64

package sha256

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/blob/0207f975d0f03724f6f5784692ee2725e0d340d3/ziskos/entrypoint/src/syscalls/sha256f.rs

type SyscallSha256Params struct {
	State *[4]uint64
	Input *[8]uint64
}

//go:noinline
func Sha256(arr *SyscallSha256Params) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_sha256(ptr)
}

func syscall_sha256(ptr uintptr)
