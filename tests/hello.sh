#!/bin/bash
# set -xe
CC=riscv64-linux-gnu-gcc
mkdir -p out/tests/hello
OUTPUT_DIR=out/tests/hello
$CC tests/hello.c -c  -o $OUTPUT_DIR/a.o
$CC -B. -static $OUTPUT_DIR/a.o -o $OUTPUT_DIR/out