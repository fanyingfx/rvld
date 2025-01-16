package linker

import "github.com/fanyingfx/rvld/pkg/utils"

func ReadArchiveMembers(file *File) []*File {
	utils.Assert(GetFileType(file.Contents) == FileTypeArchive)
	var strTab []byte
	var files []*File
	// skip !<arch>\n
	pos := 8
	for len(file.Contents)-pos > 1 {
		// align with 2 bytes
		if pos%2 == 1 {
			pos += 1
		}
		hdr := utils.Read[ArHdr](file.Contents[pos:])
		dataStart := pos + ArHdrSize
		pos = dataStart + hdr.GetSize()
		dataEnd := pos
		contents := file.Contents[dataStart:dataEnd]
		if hdr.IsSymtab() {
			continue
		} else if hdr.IsStrtab() {
			strTab = contents
			continue
		}
		files = append(files, &File{
			Name:     hdr.ReadName(strTab),
			Contents: contents,
			Parent:   file,
		})

	}
	return files
}
