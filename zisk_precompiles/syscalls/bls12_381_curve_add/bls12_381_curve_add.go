//go:build tamago && riscv64

package bls12_381_curve_add

import (
	"unsafe"
)

// Arguments are described here: https://github.com/0xPolygonHermez/zisk/tree/v0.16.0/ziskos/entrypoint/src/syscalls/bls12_381_curve_add.rs

type Point384 struct {
	X [6]uint64
	Y [6]uint64
}

type SyscallBls12381AddParams struct {
	P1 *Point384
	P2 *Point384
}

var Generator = Point384{
	X: [6]uint64{0xfb3af00adb22c6bb, 0x6c55e83ff97a1aef, 0xa14e3a3f171bac58, 0xc3688c4f9774b905, 0x2695638c4fa9ac0f, 0x17f1d3a73197d794},
	Y: [6]uint64{0x0caa232946c5e7e1, 0xd03cc744a2888ae4, 0x00db18cb2c04b3ed, 0xfcf5e095d5d00af6, 0xa09e30ed741d8ae4, 0x08b3f481e3aaa0f1},
}

var Neutral = Point384{
	X: [6]uint64{0, 0, 0, 0, 0, 0},
	Y: [6]uint64{0, 0, 0, 0, 0, 0},
}

//go:noinline
func Bls12381Add(arr *SyscallBls12381AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bls12_381_add(ptr)
}

func syscall_bls12_381_add(ptr uintptr)
