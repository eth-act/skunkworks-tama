.PHONY: all clean build-tamago build-zisk init-submodules compile-empty compile-addition compile-geth-stateless run-empty-emu run-empty-emu-quiet run-addition run-addition-verbose run-empty-rom setup-empty-rom

TAMAGO_DIR = tamago-go-latest
TAMAGO_SRC = $(TAMAGO_DIR)/src
TAMAGO_BIN = $(TAMAGO_DIR)/bin
TAMAGO = GOROOT=$(PWD)/$(TAMAGO_DIR) $(TAMAGO_BIN)/go

ZISK_DIR = zisk
ZISKEMU = $(realpath $(ZISK_DIR)/target/debug/ziskemu)
CARGO_ZISK = $(ZISK_DIR)/target/debug/cargo-zisk

# Compilation flags for TamaGo
GCFLAGS = -gcflags=""

LDFLAGS_INTERNAL = -ldflags="\
       -T 0x80001000 -D 0xa0020000"

# Default to external linking
LDFLAGS = $(LDFLAGS_EXTERNAL)
TAGS = -tags tamago,linkcpuinit,linkramstart,linkramsize,linkprintk,tinygo.wasm,tinygo,riscv64

all: init-submodules build-tamago build-zisk compile-empty run-empty-emu-quiet

init-submodules:
	@git submodule update --init --depth 1

build-tamago:
	cd $(TAMAGO_SRC) && ./make.bash
	@cd $(TAMAGO_BIN) && export TAMAGO=$$(pwd)/go && echo "TAMAGO set to: $$TAMAGO"

build-zisk:
	cd $(ZISK_DIR) && cargo build -p ziskemu -p cargo-zisk
	@cd $(ZISK_DIR) && export ZISKEMU=$$(pwd)/target/debug/ziskemu && echo "ZISKEMU set to: $$ZISKEMU"

build: build-tamago build-zisk

clean:
	cd $(ZISK_DIR) && cargo clean

compile-zisk-precompiles:
	cd zisk_precompiles && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) .

test-zisk-precompiles:
	cd zisk_precompiles && ZISKEMU=$(ZISKEMU) CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../$(TAMAGO_BIN)/go test $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -p 1 -v -exec="$(CURDIR)/test-runner.sh" ./...

compile-calculator-int:
	cd tama-programs/calculator-int && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -o calculator-int.elf .
	@echo "=== Generating witness for calculator-int program ==="
	go run ./tama-witgen/two-int32s/main.go > ./tama-programs/calculator-int/witness.bin

compile-calculator-float:
	cd tama-programs/calculator-float && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -o calculator-float.elf .
	@echo "=== Generating witness for calculator-float program ==="
	go run ./tama-witgen/two-int32s/main.go > ./tama-programs/calculator-float/witness.bin

compile-empty:
	cd tama-programs/empty && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -o empty.elf .

compile-addition:
	cd tama-programs/addition && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -o addition.elf .
	@echo "=== Generating witness for addition program ==="
	go run ./tama-witgen/two-int32s/main.go > ./tama-programs/addition/witness.bin

compile-stateless:
	cd tama-programs/stateless && CGO_ENABLED=0 GOROOT=$(PWD)/$(TAMAGO_DIR) GOOS=tamago GOARCH=riscv64 ../../$(TAMAGO_BIN)/go build $(GCFLAGS) $(LDFLAGS_INTERNAL) $(TAGS) -o stateless.elf .
	@echo "=== Generating witness for stateless program ==="
# 	-mod=read-only can be removed later. Mainly here because geth uses tablewriter 0.0.5 and this repo keeps updating to v1 	
	cd tama-witgen/stateless && GO111MODULE=on go run -mod=readonly main.go

run-calculator-int-emu-quiet:
	$(ZISKEMU) --elf tama-programs/calculator-int/calculator-int.elf -i tama-programs/calculator-int/witness.bin -c

run-calculator-float-emu-quiet:
	$(ZISKEMU) --elf tama-programs/calculator-float/calculator-float.elf -i tama-programs/calculator-float/witness.bin -c

run-empty-emu:
	$(ZISKEMU) --elf tama-programs/empty/empty.elf -v -c

run-empty-emu-quiet:
	$(ZISKEMU) --elf tama-programs/empty/empty.elf -c

run-addition:
	$(ZISKEMU) --elf tama-programs/addition/addition.elf -i tama-programs/addition/witness.bin -c

run-addition-verbose:
	$(ZISKEMU) --elf tama-programs/addition/addition.elf -i tama-programs/addition/witness.bin -v -c

run-empty-rom:
	$(CARGO_ZISK) rom-setup --elf tama-programs/empty/empty.elf

run-stateless-rom:
	$(CARGO_ZISK) rom-setup --elf tama-programs/stateless/stateless.elf -v

run-stateless:
	$(ZISKEMU) --elf tama-programs/stateless/stateless.elf -i tama-programs/stateless/witness.bin -c

run-stateless-stats:
	$(ZISKEMU) --elf tama-programs/stateless/stateless.elf -i tama-programs/stateless/witness.bin -c --stats

# Setup ROM for cargo-zisk run
.PHONY: setup-empty-rom
setup-empty-rom: compile-empty
	@echo "=== Setting up ROM with cargo-zisk ==="
	$(CARGO_ZISK) rom-setup --elf tama-programs/empty/empty.elf
