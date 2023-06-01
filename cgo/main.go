package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
#cgo CFLAGS:-I /home/pxf/testgo/cgo
#cgo LDFLAGS:-L/home/pxf/testgo/cgo -lhello
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include "hello.h"

typedef struct person{
	int Num;
	int *Score;
}person;

void func1(int *arr)
{
	printf("%d\n",arr[0]);
	printf("%d\n",arr[1]);
}

void func2(person *p)
{
	p->Num = 5;
	p->Score = (int*)malloc(sizeof(int)*5);
	p->Score[0] = 10;
	p->Score[1] = 20;
	p->Score[2] = 30;
	p->Score[3] = 40;
	p->Score[4] = 50;
}
*/
import "C"

func ToGoFloat32Slice(p unsafe.Pointer, n int) []float32 {
	var slice []float32
	sliceptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceptr.Data = uintptr(p)
	sliceptr.Len = n
	sliceptr.Cap = n
	return slice
}

func ToGoFloat64Slice(p unsafe.Pointer, n int) []float64 {
	var slice []float64
	sliceptr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceptr.Data = uintptr(p)
	sliceptr.Len = n
	sliceptr.Cap = n
	return slice
}

type Person struct {
	Score []float64
}

func main() {
	// var p Person

	// var buf bytes.Buffer
	// encoder := gob.NewEncoder(&buf)
	// encoder.Encode(&p)

	var m = make(map[int]int)
	m[0] = 1
	m[1] += 2
	fmt.Println(m)
}

/* *****************cgo 准则*****************
1.c字符串:
s := C.CString("123")
defer C.free(unsafe.Pointer(s))

2.传递数组:
arr := (*C.int)(C.malloc(8))
defer C.free(unsafe.Pointer(arr))
*(*int)(unsafe.Pointer(arr)) = 10
*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(arr)) + 4)) = 10
C.func1(arr)

3.

4.

*/
