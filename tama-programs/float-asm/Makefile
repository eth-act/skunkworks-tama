# Simple RISC-V assembly test for floating-point operations

CC = riscv64-linux-gnu-gcc
OBJCOPY = riscv64-linux-gnu-objcopy
CFLAGS = -march=rv64ifd -mabi=lp64d -nostdlib -nostartfiles -static -T linker.ld

test.elf: test.s linker.ld
	$(CC) $(CFLAGS) -o test.elf test.s

clean:
	rm -f test.elf

.PHONY: clean