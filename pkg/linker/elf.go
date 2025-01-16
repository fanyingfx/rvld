package linker

import (
	"bytes"
	"debug/elf"
	"strconv"
	"strings"
	"unsafe"

	"github.com/fanyingfx/rvld/pkg/utils"
)

const EhdrSize = int(unsafe.Sizeof(Ehdr{}))
const ShdrSize = int(unsafe.Sizeof(Shdr{}))
const SymSize = int(unsafe.Sizeof(Sym{}))
const ArHdrSize = int(unsafe.Sizeof(ArHdr{}))

type Ehdr struct {
	Ident     [16]uint8
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint64
	PhOff     uint64
	ShOff     uint64
	Flags     uint32
	EhSize    uint16
	PhEntSize uint16
	PhNum     uint16
	ShEnSize  uint16
	ShNum     uint16
	ShStrndx  uint16
}
type Shdr struct {
	Name      uint32
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Link      uint32
	Info      uint32
	AddrAlgin uint64
	EntSize   uint64
}
type Sym struct {
	Name  uint32
	Info  uint8
	Other uint8
	Shndx uint16
	Val   uint64
	Size  uint64
}

func (s *Sym) IsAbs() bool {
	return s.Shndx == uint16(elf.SHN_ABS)
}
func (s *Sym) IsUndef() bool {
	return s.Shndx == uint16(elf.SHN_UNDEF)
}

type ArHdr struct {
	Name [16]byte
	Date [12]byte
	Uid  [6]byte
	Gid  [6]byte
	Mode [8]byte
	Size [10]byte
	Fmag [2]byte
}

func (arhdr *ArHdr) HasPrefix(s string) bool {
	return strings.HasPrefix(string(arhdr.Name[:]), s)
}
func (arhdr *ArHdr) IsStrtab() bool {
	return arhdr.HasPrefix("// ")
}
func (arhdr *ArHdr) IsSymtab() bool {
	return arhdr.HasPrefix("/ ") || arhdr.HasPrefix("/SYM64/ ")
}

func (arhdr *ArHdr) GetSize() int {
	size, err := strconv.Atoi(strings.TrimSpace(string(arhdr.Size[:])))
	utils.MustNo(err)
	return size

}
func (arhdr *ArHdr) ReadName(strTab []byte) string {
	// long filename
	if arhdr.HasPrefix("/") {
		start, err := strconv.Atoi(strings.TrimSpace(string(arhdr.Name[1:])))
		utils.MustNo(err)

		end := start + bytes.Index(strTab[start:], []byte("/\n"))
		return string(strTab[start:end])

	}
	// short filename
	end := bytes.Index(arhdr.Name[:], []byte("/"))
	utils.Assert(end != -1)
	return string(arhdr.Name[:end])

}

func ElfGetName(strTab []byte, offset uint32) string {
	length := uint32(bytes.Index(strTab[offset:], []byte{0}))
	return string(strTab[offset : offset+length])
}
