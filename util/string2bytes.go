package util

import (
	"reflect"
	"unsafe"
)

//go:noinline
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(sh))
}

//go:noinline
func Bytes2String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
