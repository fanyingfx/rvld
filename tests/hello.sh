#!/bin/bash
mkdir -p out/tests/hello
OUTPUT_DIR=out/tests
riscv64-linux-gnu-gcc tests/hello.c -c  -o out/tests/hello/a.o
./rvld $OUTPUT_DIR/hello/a.o