package linker

import "bytes"

const MAGIC_HEADER = "\177ELF"

func CheckMagic(contents []byte) bool {

	return bytes.HasPrefix(contents, []byte(MAGIC_HEADER))

}
func WriteMagic(contents []byte) {
	copy(contents, []byte(MAGIC_HEADER))
}
