package linker

type ContextArgs struct {
	Output       string
	Emulation    MachineType
	LibraryPaths []string
}
type Context struct {
	Args ContextArgs
	Buf  []byte

	Chunks         []IChunk
	Ehdr           *OutputEhdr
	Shdr           *OutputShdr
	Phdr           *OutputPhdr
	Got            *GotSection
	TpAddr         uint64 // thread local pointer
	OutputSections []*OutputSection
	Objs           []*ObjectFile
	SymbolMap      map[string]*Symbol
	MergedSections []*MergedSection
}

func NewContext() *Context {
	return &Context{
		Args:      ContextArgs{Output: "a.out", Emulation: MachineTypeNone},
		SymbolMap: make(map[string]*Symbol),
	}
}
