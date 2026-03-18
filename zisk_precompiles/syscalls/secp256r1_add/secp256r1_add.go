//go:build tamago && riscv64

package secp256r1_add

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/secp256r1_add.rs

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallSecp256r1AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [4]uint64{0xF4A13945D898C296, 0x77037D812DEB33A0, 0xF8BCE6E563A440F2, 0x6B17D1F2E12C4247},
	Y: [4]uint64{0xCBB6406837BF51F5, 0x2BCE33576B315ECE, 0x8EE7EB4A7C0F9E16, 0x4FE342E2FE1A7F9B},
}

var Neutral = Point256{
	X: [4]uint64{0, 0, 0, 0},
	Y: [4]uint64{0, 0, 0, 0},
}

//go:noinline
func Secp256r1Add(arr *SyscallSecp256r1AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_secp256r1_add(ptr)
}

func syscall_secp256r1_add(ptr uintptr)
