package linker

import (
	"github.com/fanyingfx/rvld/pkg/utils"
)

type InputFile struct {
	File        *File
	ElfSections []Shdr
}

func NewInputFile(file *File) InputFile {
	if len(file.Contents) < EhdrSize {
		utils.Fatal("file too small")
	}
	if !CheckMagic(file.Contents) {
		utils.Fatal("not an ELF file")
	}

	ehdr := utils.Read[Ehdr](file.Contents)
	contents := file.Contents[ehdr.ShOff:]
	shdr := utils.Read[Shdr](contents)
	numSections := uint64(ehdr.ShNum)
	if ehdr.ShNum == 0 {
		numSections = shdr.Size
	}
	f := InputFile{File: file}
	f.ElfSections = []Shdr{shdr}
	for numSections > 1 {
		contents = contents[ShdrSize:]
		f.ElfSections = append(f.ElfSections, utils.Read[Shdr](contents))
		numSections--
	}

	return f
}
