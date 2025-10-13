#/bin/bash

riscv64-unknown-elf-gcc -fPIC -nostdlib -c func.c -o func.o
ar rcs libfunc.a func.o

