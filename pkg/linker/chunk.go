package linker

type IChunk interface {
	GetName() string
	GetShdr() *Shdr
	UpdateShdr(ctx *Context)
	GetShndx() int64
	CopyBuf(ctx *Context)
}

type Chunk struct {
	Name  string
	Shdr  Shdr
	Shndx int64
}

func NewChunk() Chunk {
	return Chunk{
		Shdr: Shdr{AddrAlign: 1},
	}
}
func (c *Chunk) GetName() string {
	return c.Name
}
func (c *Chunk) GetShdr() *Shdr {
	return &c.Shdr
}
func (c *Chunk) UpdateShdr(ctx *Context) {}
func (c *Chunk) GetShndx() int64 {
	return c.Shndx
}

func (c *Chunk) CopyBuf(ctx *Context) {}
