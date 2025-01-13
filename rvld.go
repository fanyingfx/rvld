package main

import (
	"os"

	"github.com/fanyingfx/rvld/pkg/linker"
	"github.com/fanyingfx/rvld/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.Fatal("wrong args")
	}

	filename := os.Args[1]
	file := linker.MustNewFile(filename)
	inputfile := linker.NewInputFile(file)
	utils.Assert(len(inputfile.ElfSections) == 11)

}
