//go:build tamago && riscv64

package main

import (
	"fmt"
	"tamagotest/tamaboards/zkvm"
)

func main() {
	fmt.Println("=== TamaGo Program Starting ===")
	fmt.Println("Program is running on ZisK VM...")
	fmt.Println("Calling zkvm.Shutdown()...")
	zkvm.Shutdown()
	fmt.Println("This line should never be reached!")
}
