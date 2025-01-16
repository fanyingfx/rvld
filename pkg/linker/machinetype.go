package linker

import (
	"debug/elf"

	"github.com/fanyingfx/rvld/pkg/utils"
)

type MachineType uint8

const (
	MachineTypeNone    MachineType = iota
	MachineTypeRISCV64 MachineType = iota
)

func GetMachineTypeFromContents(contents []byte) MachineType {
	ftype := GetFileType(contents)
	switch ftype {

	case FileTypeObject:
		machine := utils.Read[uint16](contents[18:])
		if machine == uint16(elf.EM_RISCV) {
			class := elf.Class(contents[4])
			switch class {
			case elf.ELFCLASS64:
				return MachineTypeRISCV64
			}
		}
		// case FileTypeEmpty:
		// case FileTypeUnknown:
		// default:
		// return MachineTypeNone
	}
	return MachineTypeNone
}

func (m MachineType) String() string {
	switch m {
	case MachineTypeRISCV64:
		return "riscv64"

	}
	utils.Assert(m == MachineTypeNone)
	return "none"

}
