// +build 386

package main

import (
	"unsafe"
)

type VariantF64 struct {
	ptr unsafe.Pointer
	lenOrType int
	//capOrVal int
	//last32bit int // used for second half of float64.
	f float64
}

func Float64VariantF64(v float64) (r VariantF64) {
	return VariantF64{
		lenOrType:1,
		// First half of float64.
		//capOrVal: *(*int)(unsafe.Pointer(&v)),
		//// Second half of float64.
		//last32bit: *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&v))+unsafe.Sizeof(int(0)))),
		//f: v,
	}
}

func (v* VariantF64) Float64() float64 {
	//return *(*float64)(unsafe.Pointer(&v.capOrVal))
	// return v.f
	return *(*float64)(unsafe.Pointer(&v.f))
	//return math.Float64frombits(v.f)
}

func main() {
	v:=VariantF64{}
	v.Float64()
}