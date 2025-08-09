//go:build tamago && riscv64

package main

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	_ "tamagotest/tamaboards/zkvm"
)

func main() {
	core.ExecuteStateless(nil, vm.Config{}, nil, nil)
}