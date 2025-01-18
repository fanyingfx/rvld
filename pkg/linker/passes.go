package linker

import (
	"debug/elf"
	"math"
	"sort"

	"github.com/fanyingfx/rvld/pkg/utils"
)

func CreateInternalFile(ctx *Context) {
	obj := &ObjectFile{}
	// ctx.InternalObj = obj
	ctx.Objs = append(ctx.Objs, obj)

	// ctx.InternalEsyms = make([]Sym, 1)
	obj.Symbols = append(obj.Symbols, NewSymbol(""))

	obj.FirstGlobal = 1
	obj.IsAlive = true

	// obj.ElfSyms = ctx.InternalEsyms

}

func ResolveSymbols(ctx *Context) {
	for _, file := range ctx.Objs {
		file.ResolveSymbols()
	}
	MarkLiveObjects(ctx)
	for _, file := range ctx.Objs {
		if !file.IsAlive {
			file.ClearSymbols()
		}
	}
	ctx.Objs = utils.RemoveIf(ctx.Objs, func(file *ObjectFile) bool {
		return !file.IsAlive
	})

}

func MarkLiveObjects(ctx *Context) {
	roots := make([]*ObjectFile, 0)
	for _, file := range ctx.Objs {
		if file.IsAlive {
			roots = append(roots, file)
		}
	}
	utils.Assert(len(roots) > 0)
	for len(roots) > 0 {
		file := roots[0]
		if !file.IsAlive {
			continue
		}
		file.MarkLiveObjects(func(file *ObjectFile) {
			roots = append(roots, file)
		})
		roots = roots[1:]

	}

}
func RegisterSectionPieces(ctx *Context) {
	for _, file := range ctx.Objs {
		file.RegisterSectionPieces()
	}
}
func CreateSyntheticSections(ctx *Context) {
	push := func(chunk IChunk) IChunk {
		ctx.Chunks = append(ctx.Chunks, chunk)
		return chunk
	}
	ctx.Ehdr = push(NewOutputEhdr()).(*OutputEhdr)
	ctx.Shdr = push(NewOutputShdr()).(*OutputShdr)
	ctx.Phdr = push(NewOutputPhdr()).(*OutputPhdr)
	ctx.Got = push(NewGotSection()).(*GotSection)
}

func SetOutputSectionOffsets(ctx *Context) uint64 {
	addr := IMAGE_BASE
	for _, chunk := range ctx.Chunks {
		if chunk.GetShdr().Flags&uint64(elf.SHF_ALLOC) == 0 {
			continue
		}
		addr = utils.AlignTo(addr, chunk.GetShdr().AddrAlign)
		chunk.GetShdr().Addr = addr
		if !isTbss(chunk) {
			addr += chunk.GetShdr().Size
		}
	}
	i := 0
	first := ctx.Chunks[0]
	for {
		shdr := ctx.Chunks[i].GetShdr()
		shdr.Offset = shdr.Addr - first.GetShdr().Addr
		i++
		if i >= len(ctx.Chunks) || ctx.Chunks[i].GetShdr().Flags&uint64(elf.SHF_ALLOC) == 0 {
			break
		}
	}
	lastShdr := ctx.Chunks[i-1].GetShdr()
	fileoff := lastShdr.Offset + lastShdr.Size

	for ; i < len(ctx.Chunks); i++ {
		shdr := ctx.Chunks[i].GetShdr()
		fileoff = utils.AlignTo(fileoff, shdr.AddrAlign)
		shdr.Offset = fileoff
		fileoff += shdr.Size
	}
	ctx.Phdr.UpdateShdr(ctx)

	return fileoff
}

func BinSections(ctx *Context) {
	group := make([][]*InputSection, len(ctx.OutputSections))
	for _, file := range ctx.Objs {
		for _, isec := range file.Sections {
			if isec == nil || !isec.IsAlive {
				continue
			}
			idx := isec.OutputSection.Idx
			group[idx] = append(group[idx], isec)
		}
	}
	for idx, osec := range ctx.OutputSections {
		osec.Members = group[idx]
	}
}

func CollectOutputSections(ctx *Context) []IChunk {
	osecs := make([]IChunk, 0)
	for _, osec := range ctx.OutputSections {
		if len(osec.Members) > 0 {
			osecs = append(osecs, osec)
		}

	}
	for _, osec := range ctx.MergedSections {
		if osec.Shdr.Size > 0 {
			osecs = append(osecs, osec)
		}
	}
	return osecs
}
func ComputeSectionSizes(ctx *Context) {
	for _, osec := range ctx.OutputSections {
		offset := uint64(0)
		p2align := int64(0)
		for _, isec := range osec.Members {
			offset = utils.AlignTo(offset, 1<<isec.P2Align)
			isec.Offset = uint32(offset)
			offset += uint64(isec.ShSize)
			if p2align < int64(isec.P2Align) {
				p2align = int64(isec.P2Align)
			}
		}
		osec.Shdr.Size = offset
		osec.Shdr.AddrAlign = 1 << p2align
	}
}
func SortOutputSections(ctx *Context) {
	rank := func(chunk IChunk) int32 {
		typ := chunk.GetShdr().Type
		flags := chunk.GetShdr().Flags
		if chunk == ctx.Ehdr {
			return 0
		}
		if chunk == ctx.Phdr {
			return 1
		}
		if typ == uint32(elf.SHT_NOTE) {
			return 2
		}
		if flags&uint64(elf.SHF_ALLOC) == 0 {
			return math.MaxInt32 - 1
		}
		if chunk == ctx.Shdr {
			return math.MaxInt32
		}
		b2i := func(b bool) int {
			if b {
				return 1
			}
			return 0
		}
		writeable := b2i(flags&uint64(elf.SHF_WRITE) != 0)
		notExec := b2i(flags&uint64(elf.SHF_EXECINSTR) == 0)
		notTls := b2i(flags&uint64(elf.SHF_TLS) == 0)
		isBss := b2i(typ == uint32(elf.SHT_NOBITS))
		return int32(writeable<<7 | notExec<<6 | notTls<<5 | isBss<<4)

	}
	sort.SliceStable(ctx.Chunks, func(i, j int) bool {
		return rank(ctx.Chunks[i]) < rank(ctx.Chunks[j])
	})
}
func ComputeMergedSectionSizes(ctx *Context) {
	for _, osec := range ctx.MergedSections {
		osec.AssignOffsets()
	}
}
func ScanRelocations(ctx *Context) {
	for _, file := range ctx.Objs {
		file.ScanRelocations()
	}
	syms := make([]*Symbol, 0)
	for _, file := range ctx.Objs {
		for _, sym := range file.Symbols {
			if sym.File == file && sym.Flags != 0 {
				syms = append(syms, sym)
			}
		}
	}
	for _, sym := range syms {
		if sym.Flags&NeedsGotTp != 0 {
			ctx.Got.AddGotTpSymbol(sym)
		}
		sym.Flags = 0
	}
}

func isTbss(chunk IChunk) bool {
	shdr := chunk.GetShdr()
	return shdr.Type == uint32(elf.SHT_NOBITS) && shdr.Flags&uint64(elf.SHF_TLS) != 0
}
func (m *MergedSection) CopyBuf(ctx *Context) {
	buf := ctx.Buf[m.Shdr.Offset:]
	for key := range m.Map {
		if frag, ok := m.Map[key]; ok {
			copy(buf[frag.Offset:], key)
		}
	}
}
