//go:build tamago && riscv64

package bn254

import (
	"unsafe"
)

type Point256 struct {
	X [4]uint64
	Y [4]uint64
}

type Complex struct {
	X [4]uint64
	Y [4]uint64
}

type SyscallBn254ComplexAddParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBn254ComplexMulParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBn254ComplexSubParams struct {
	F1 *Complex
	F2 *Complex
}

type SyscallBn254AddParams struct {
	P1 *Point256
	P2 *Point256
}

var Generator = Point256{
	X: [4]uint64{0x0000000000000001, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
	Y: [4]uint64{0x0000000000000002, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000},
}

var Neutral = Point256{
	X: [4]uint64{0, 0, 0, 0},
	Y: [4]uint64{0, 0, 0, 0},
}

//go:noinline
func Bn254Add(arr *SyscallBn254AddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_add(ptr)
}

//go:noinline
func Bn254Dbl(arr *Point256) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_dbl(ptr)
}

//go:noinline
func Bn254ComplexAdd(arr *SyscallBn254ComplexAddParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_complex_add(ptr)
}

//go:noinline
func Bn254ComplexMul(arr *SyscallBn254ComplexMulParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_complex_mul(ptr)
}

//go:noinline
func Bn254ComplexSub(arr *SyscallBn254ComplexSubParams) {
	ptr := uintptr(unsafe.Pointer(arr))
	syscall_bn254_complex_sub(ptr)
}

func syscall_bn254_add(ptr uintptr)
func syscall_bn254_dbl(ptr uintptr)
func syscall_bn254_complex_add(ptr uintptr)
func syscall_bn254_complex_mul(ptr uintptr)
func syscall_bn254_complex_sub(ptr uintptr)
