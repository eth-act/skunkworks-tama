//go:build tamago && riscv64

package main

import (
	"fmt"
	"math/big"

	"github.com/eth-act/skunkworks-tama/tamaboards/zkvm/zisk_runtime"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/stateless"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	// Read witness bytes from zkVM input buffer
	// The witness generator will serialize this as []byte
	witnessBytes := zisk_runtime.Read[[]byte]()
	fmt.Printf("Read witness (%d bytes)\n", len(witnessBytes))

	// Deserialize the witness using RLP
	witness := new(stateless.Witness)
	if err := rlp.DecodeBytes(witnessBytes, witness); err != nil {
		panic(fmt.Sprintf("Failed to decode witness: %v", err))
	}
	fmt.Printf("Witness decoded - contains %d headers, %d state nodes, %d code entries\n",
		len(witness.Headers), len(witness.State), len(witness.Codes))

	// Read block bytes from zkVM input buffer
	blockBytes := zisk_runtime.Read[[]byte]()
	fmt.Printf("Read block (%d bytes)\n", len(blockBytes))

	// Deserialize the block using RLP
	block := new(types.Block)
	if err := rlp.DecodeBytes(blockBytes, block); err != nil {
		panic(fmt.Sprintf("Failed to decode block: %v", err))
	}
	fmt.Printf("Block decoded - #%d with %d transactions\n", 
		block.Number().Uint64(), len(block.Transactions()))

	// Use the appropriate chain config
	config := *params.AllEthashProtocolChanges
	config.TerminalTotalDifficulty = common.Big0
	config.MergeNetsplitBlock = big.NewInt(11)

	// Set times for Shanghai and Cancun based on block timestamp
	timestamp := block.Time()
	shanghaiTime := timestamp - 10
	cancunTime := timestamp - 5
	config.ShanghaiTime = &shanghaiTime
	config.CancunTime = &cancunTime
	config.BlobScheduleConfig = params.DefaultBlobSchedule

	// Create VM config
	vmConfig := vm.Config{}

	// Execute stateless with the witness and block
	stateRoot, receiptRoot, err := core.ExecuteStateless(&config, vmConfig, block, witness)
	if err != nil {
		panic(fmt.Sprintf("ExecuteStateless failed: %v", err))
	}

	// Log the results
	fmt.Printf("ExecuteStateless succeeded!\n")
	fmt.Printf("State root: %x\n", stateRoot)
	fmt.Printf("Receipt root: %x\n", receiptRoot)
	fmt.Printf("Gas used: %d\n", block.GasUsed())

	// Commit the results as public outputs
	// We'll commit both state root and receipt root as byte arrays
	zisk_runtime.Commit(stateRoot[:])
	zisk_runtime.Commit(receiptRoot[:])

	// Verify non-empty roots
	if stateRoot == (common.Hash{}) {
		panic("State root is empty")
	}
	if receiptRoot == (common.Hash{}) {
		panic("Receipt root is empty")
	}
}