.PHONY: all clean build-tamago build-zisk compile-empty run-empty-emu run-empty-emu-quiet run-empty-rom setup-empty-rom

TAMAGO_DIR = tamago-go-latest
TAMAGO_SRC = $(TAMAGO_DIR)/src
TAMAGO_BIN = $(TAMAGO_DIR)/bin
TAMAGO = GOROOT=$(PWD)/$(TAMAGO_DIR) $(TAMAGO_BIN)/go

ZISK_DIR = zisk
ZISKEMU = $(ZISK_DIR)/target/debug/ziskemu
CARGO_ZISK = $(ZISK_DIR)/target/debug/cargo-zisk

# Compilation flags for TamaGo
GCFLAGS = -gcflags="all=-d=softfloat"

LDFLAGS_INTERNAL = -ldflags="\
       -T 0x80001000 -D 0xa0020000"

# Default to external linking
LDFLAGS = $(LDFLAGS_EXTERNAL)
TAGS = -tags tamago,linkcpuinit,linkramstart,linkramsize,linkprintk

all: build-tamago build-zisk compile-empty run-empty-emu-quiet

build-tamago:
	cd $(TAMAGO_SRC) && ./make.bash
	@cd $(TAMAGO_BIN) && export TAMAGO=$$(pwd)/go && echo "TAMAGO set to: $$TAMAGO"

build-zisk:
	cd $(ZISK_DIR) && cargo build -p ziskemu -p cargo-zisk
	@cd $(ZISK_DIR) && export ZISKEMU=$$(pwd)/target/debug/ziskemu && echo "ZISKEMU set to: $$ZISKEMU"

build: build-tamago build-zisk

clean:
	cd $(ZISK_DIR) && cargo clean

compile-empty:
	cd tama-programs/empty && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -o empty.elf .

run-empty-emu:
	cd tama-programs/empty && ../../$(ZISKEMU) --elf empty.elf -v -c

run-empty-emu-quiet:
	cd tama-programs/empty && ../../$(ZISKEMU) --elf empty.elf -c

run-empty-rom:
	$(CARGO_ZISK) rom-setup --elf tama-programs/empty/empty.elf

# Setup ROM for cargo-zisk run
.PHONY: setup-empty-rom
setup-empty-rom: compile-empty
	@echo "=== Setting up ROM with cargo-zisk ==="
	$(CARGO_ZISK) rom-setup --elf tama-programs/empty/empty.elf