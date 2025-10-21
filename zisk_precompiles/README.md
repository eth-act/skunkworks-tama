# Rust and Go syscall parameters types compatibility

The syscall parameter types are a direct translation from Rust to Go.
As it happens their memory representation is the same. Go doesn't provide
any guarantees for the memory layout of Go struct thus the current implementation
is depedent on the compiler version. For the sake of PoC that's sufficient.
Ultimately parameters to syscalls could be represented just as arrays of bytes
that are serialized and deserialized to Go types or CGO could be used if
supported by TamaGo.

# Testing
Use `make test-zisk-precompiles` target to execute tests.

Basic tests are provided for all syscalls except for sha256.

