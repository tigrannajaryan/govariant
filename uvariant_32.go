// +build 386

package main

import (
	"unsafe"
)
import "reflect"

type Variant struct {
	ptr unsafe.Pointer
	lenOrType int
	//last32bit int // used for second half of float64.
	//capOrVal [8]byte
	capOrVal int64
	//capOrVal int
}

func (v* Variant) Type() VariantType {
	if v.ptr == nil {
		return VariantType(v.lenOrType)
	}

	if v.capOrVal == -1 {
		return VariantTypeString
	}

	return VariantTypeBytes
}

func IntVariant(v int) (r Variant) {
	r.lenOrType = VariantTypeInt
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&r.capOrVal)))) = v
	return r
}

func Float64Variant(v float64) (r Variant) {
	r.lenOrType = VariantTypeFloat64
	*(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(&r.capOrVal)))) = v
	//r.capOrVal = v
	return r
}

func StringVariant(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType:hdr.Len, capOrVal:-1}
}

func BytesVariant(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	vr := Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType:hdr.Len}
	*(*int)(unsafe.Pointer(&vr.capOrVal)) = hdr.Cap
	return vr
}

func (v* Variant) Int() int {
	return *(*int)(unsafe.Pointer(&v.capOrVal))
}

func (v* Variant) Float64() float64 {
	return *(*float64)(unsafe.Pointer(&v.capOrVal))
	//return v.capOrVal
}

func (v* Variant) String() (s string) {
	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	src := (*reflect.StringHeader)(unsafe.Pointer(v))
	*dest = *src
	return s
}

func (v* Variant) Bytes() (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	src := (*reflect.SliceHeader)(unsafe.Pointer(v))
	*dest = *src
	return b
}
