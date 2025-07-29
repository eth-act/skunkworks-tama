.PHONY: all clean

TAMAGO_DIR = tamago-go-latest
TAMAGO_SRC = $(TAMAGO_DIR)/src
TAMAGO_BIN = $(TAMAGO_DIR)/bin

all: build

build:
	cd $(TAMAGO_SRC) && ./all.bash
	cd $(TAMAGO_BIN) && export TAMAGO=$$(pwd)/go && echo "TAMAGO set to: $$TAMAGO"