package linker

import (
	"bytes"
	"debug/elf"

	"github.com/fanyingfx/rvld/pkg/utils"
)

type FileType uint8

const (
	FileTypeUnknown FileType = iota
	FileTypeEmpty   FileType = iota
	FileTypeObject  FileType = iota
	FileTypeArchive FileType = iota
)

func GetFileType(contents []byte) FileType {
	if len(contents) == 0 {
		return FileTypeEmpty
	}
	if CheckMagic(contents) {
		elfType := elf.Type(utils.Read[uint16](contents[16:]))
		switch elfType {
		case elf.ET_REL:
			return FileTypeObject

		}

		return FileTypeUnknown
	}
	if bytes.HasPrefix(contents, []byte("!<arch>\n")) {
		return FileTypeArchive

	}
	return FileTypeUnknown
}
func CheckFileCompatibility(ctx *Context, file *File) {
	mt := GetMachineTypeFromContents(file.Contents)
	if mt != ctx.Args.Emulation {
		utils.Fatal("incompatible file type")
	}
}
