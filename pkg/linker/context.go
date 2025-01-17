package linker

type ContextArgs struct {
	Output       string
	Emulation    MachineType
	LibraryPaths []string
}
type Context struct {
	Args           ContextArgs
	Chunks         []IChunk
	Ehdr           *OutputEhdr
	Shdr           *OutputShdr
	OutputSections []*OutputSection
	Buf            []byte
	Objs           []*ObjectFile
	SymbolMap      map[string]*Symbol
	MergedSections []*MergedSection
	InternalObj    *ObjectFile
	InternalEsyms  []Sym
}

func NewContext() *Context {
	return &Context{
		Args:      ContextArgs{Output: "a.out", Emulation: MachineTypeNone},
		SymbolMap: make(map[string]*Symbol),
	}
}
