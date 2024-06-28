package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//x := uint32(4294967295)
	y := int32(22222)
	s := *(*[4]byte)(unsafe.Pointer(&y))
	fmt.Println(s)
}
