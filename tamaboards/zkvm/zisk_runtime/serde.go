//go:build tamago && riscv64

package zisk_runtime

import (
	"crypto/sha256"
	"hash"
	"github.com/eth-act/skunkworks-tama/tamaboards/serialization"
)

// publicInputs will store a hash of all public inputs
// This is an optimization, where instead of making all public inputs
// public, we make a hash of the public inputs publish instead.
// TODO: This is curently unused. See `Commit`
var publicInputsHash hash.Hash = sha256.New()

// Read deserializes any supported type from the input buffer
// The first 64 bits specify the length of the serialized data
//
// Usage examples:
//
//	x := Read[int64]()           // Read an int64
//	s := Read[string]()          // Read a string
//	p := Read[Person]()          // Read a struct
//	data := Read[[]byte]()       // Read a byte slice
func Read[T any]() T {
	// Read the length first (8 bytes)
	// TODO: Note that this wastes bytes if we are serializing small structures like u8
	// TODO: since that will take up 72 bits. We can optimize this later.
	lengthBytes := UnsafeReadBytes(8)
	length := uint64(0)
	for i := 0; i < 8; i++ {
		length |= uint64(lengthBytes[i]) << (i * 8)
	}

	// Read the serialized data of the specified length
	data := UnsafeReadBytes(int(length))

	// Deserialize it into the result
	var result T
	serialization.DeserializeData(data, &result)

	return result
}

// Commit serializes any supported type and commits it as a public input/output
//
// Note: The public inputs optimization would modify the below to only call
// CommitBytes on the hash.
//
// Usage examples:
//
//	Commit(result)               // Commit a computation result
//	Commit(hash)                 // Commit a hash value
//	Commit(publicKey)            // Commit a public key
func Commit[T any](value T) {
	// Serialize the value
	data := serialization.MustSerializeData(value)

	// TODO: Hash all public inputs and make this commit a hash instead
	// TODO: This can be hidden from the verifier by making a PublicInputs struct
	// TODO: That takes all of the public inputs in and does the hashing or something to that nature
	// TODO: Eg: VerifyProof(proof : []byte, publicInput : PublicInput) -> bool

	// Commit the serialized bytes
	CommitBytes(data)
}
