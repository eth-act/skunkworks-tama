//go:build tamago && riscv64

package zisk_runtime

import (
	"github.com/eth-act/skunkworks-tama/tamaboards/zkvm"
	"unsafe"
)

/*

Note: This file only defines guest reads and writes(public inputs). It does not define host methods.

*/

// Input/Output pointers to track current read/write positions
var (
	// INPUT_PTR tracks the current read position in the input buffer
	// Starts at 0, which maps to INPUT_ADDR + 16 (where actual data begins)
	//
	// Zisk stores input at INPUT_ADDR with this format:
	// INPUT_ADDR + 0: 8 bytes - free input (reserved, currently 0)
	// INPUT_ADDR + 8: 8 bytes - input length (of all inputs) -- We actually do not want to use this. We want the input length per item
	// TODO: To be safe, we can take the total input length and ensure that we never read past it by comparing to INPUT_PTR
	// INPUT_ADDR + 16: actual input data
	INPUT_PTR uint64 = 0

	// OUTPUT_PTR tracks the current write position in the output buffer
	OUTPUT_PTR uint64 = 0
)

// // ReadUint64 reads the next uint64 from the input buffer and advances the pointer
// func ReadUint64() uint64 {
// 	ptr := (*uint64)(unsafe.Pointer(uintptr(INPUT_ADDR + 16 + INPUT_PTR)))
// 	INPUT_PTR += 8
// 	return *ptr
// }

// ReadBytes reads n bytes from the input buffer and advances the pointer
// This method reads the bytes from the RAM, so its the most low level
// function we have for reading input.
//
// We recommend users use the Read methods. This is "unsafe" because
// it is possible to read halfway into a struct and cause subsequent reads
// to be misaligned
func UnsafeReadBytes(n int) []byte {
	result := make([]byte, n)
	src := unsafe.Pointer(uintptr(zkvm.INPUT_ADDR + 16 + INPUT_PTR))
	dst := unsafe.Pointer(&result[0])
	// Copy the bytes
	for i := 0; i < n; i++ {
		*(*byte)(unsafe.Add(dst, i)) = *(*byte)(unsafe.Add(src, i))
	}
	INPUT_PTR += uint64(n)
	return result
	/*
		TODO: Try the following:
		func UnsafeReadBytes(n int) []byte {
		    result := make([]byte, n)
		    src := unsafe.Pointer(uintptr(zkvm.INPUT_ADDR + 16 + INPUT_PTR))
		    srcSlice := unsafe.Slice((*byte)(src), n)
		    copy(result, srcSlice)
		    INPUT_PTR += uint64(n)
		    return result
		}
	*/
}

// // WriteUint64 writes a uint64 to the output buffer and advances the pointer
// func WriteUint64(value uint64) {
// 	ptr := (*uint64)(unsafe.Pointer(uintptr(OUTPUT_ADDR + OUTPUT_PTR)))
// 	*ptr = value
// 	OUTPUT_PTR += 8
// }

// CommitBytes writes bytes to the output buffer and advances the pointer
//
// Note: This method is used for writing public inputs.
// It is not used for sending data from the host to guest.
//
// The emulator expects output in this format:
// - First 4 bytes at OUTPUT_ADDR: count of u32 values (not bytes!)
// - Then the actual data as u32 values starting at OUTPUT_ADDR + 4
// TODO: Think about function naming
func CommitBytes(data []byte) {
	if len(data) == 0 {
		return
	}

	// Write the data
	dst := unsafe.Pointer(uintptr(zkvm.OUTPUT_ADDR + 4 + OUTPUT_PTR))
	src := unsafe.Pointer(&data[0])
	// Copy the bytes
	for i := 0; i < len(data); i++ {
		*(*byte)(unsafe.Add(dst, i)) = *(*byte)(unsafe.Add(src, i))
	}
	OUTPUT_PTR += uint64(len(data))

	// Update the count header with total u32 values written so far
	// OUTPUT_PTR now contains total bytes written (since we write at OUTPUT_ADDR + 4 + OUTPUT_PTR)
	totalBytes := OUTPUT_PTR
	numU32s := (totalBytes + 3) / 4 // Round up to nearest u32
	countPtr := (*uint32)(unsafe.Pointer(uintptr(zkvm.OUTPUT_ADDR)))
	*countPtr = uint32(numU32s)

	/*
		TODO: Try the following:
		func CommitBytes(data []byte) {
		    if len(data) == 0 {
		        return // Handle empty slice case
		    }
		    dst := unsafe.Pointer(uintptr(zkvm.OUTPUT_ADDR + OUTPUT_PTR))
		    dstSlice := unsafe.Slice((*byte)(dst), len(data))
		    copy(dstSlice, data)
		    OUTPUT_PTR += uint64(len(data))
		}
	*/

}
