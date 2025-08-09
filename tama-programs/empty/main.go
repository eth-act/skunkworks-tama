//go:build tamago && riscv64

package main

import (
	"fmt"
	"unsafe"
	_ "tamagotest/tamaboards/zkvm"
)

func main() {	
	fmt.Println("Hello World")
}