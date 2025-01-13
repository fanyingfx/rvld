package linker

import "unsafe"

const EhdrSize = int(unsafe.Sizeof(Ehdr{}))
const ShdrSize = int(unsafe.Sizeof(Shdr{}))

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
	ShStrndr  uint16
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
