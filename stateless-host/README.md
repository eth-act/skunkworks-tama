That shares 99% of code with "../tama-programs/stateless/" but can be compiled
with a normal GO compiler. The data is read from the standard input.

Compile it:
- for host:
```
CGO_ENABLED=0 go build -ldflags="-s -w" -o stateless-host main.go
```

- for RISCV:
```
GOOS=linux GOARCH=riscv64 CGO_ENABLED=0 go build -ldflags="-s -w" -o stateless-riscv64 main.go
```

Then strace:
- for host:

Use "-f" option to follow threads.

```
cat witness.bin | strace -f ./stateless-host 2> strace_host
```
- for RISCV:

Following threads is not possible with qemu.

```
cat witness.bin | qemu-riscv64 -strace ./stateless-riscv64 2> strace_riscv
```

Add "-c" switch to strace to produce only a summary of syscalls instead of a detailed log.

Prior to stracing "witness.bin" has to be generated with "tama-witgen".

