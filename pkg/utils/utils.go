package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

func Fatal(v any) {
	fmt.Printf("rvld: \033[0;1;31mfatal:\033[0m %v\n", v)
	debug.PrintStack()
	os.Exit(1)
}
func Assert(condition bool) {
	if !condition {
		Fatal("assertion error")
	}
}
func MustNo(err error) {
	if err != nil {
		Fatal(err)
	}
}
func Read[T any](data []byte) (val T) {

	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &val)
	MustNo(err)
	return val
}
func ReadSlice[T any](data []byte, sz int) []T {
	nums := len(data) / sz
	res := make([]T, 0, nums)
	for nums > 0 {
		res = append(res, Read[T](data))
		data = data[sz:]
		nums -= 1
	}
	return res
}

func RemovePrefix(s string, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
		return s, true
	}
	return s, false

}
func RemoveIf[T any](elems []T, predicate func(T) bool) []T {
	i := 0
	for _, elem := range elems {
		if predicate(elem) {
			continue
		}
		elems[i] = elem
		i += 1
	}
	return elems[:i]
}
