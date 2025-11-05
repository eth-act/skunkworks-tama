//go:build tamago && riscv64

package main

import (
	"fmt"
	"github.com/eth-act/skunkworks-tama/tamaboards/zkvm/zisk_runtime"
)

func main() {	
	// Read two integers using the generic Read function
	x := zisk_runtime.Read[int64]()
	y := zisk_runtime.Read[int64]()
	
	// Calculate sum
	sum := x + y
	
	// Display result
	fmt.Printf("Input: %d + %d = %d\n", x, y, sum)
	
	// Commit result as public output (no length prefix)
	zisk_runtime.Commit(sum)
}