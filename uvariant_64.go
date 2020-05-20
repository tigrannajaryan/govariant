// +build !386

package main

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	ptr unsafe.Pointer
	lenOrType int
	capOrVal int
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

func IntVariant(v int) Variant {
	return Variant{lenOrType:VariantTypeInt, capOrVal: v}
}

func StringVariant(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType:hdr.Len, capOrVal:-1}
}

func BytesVariant(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType:hdr.Len, capOrVal: hdr.Cap}
}

func Float64Variant(v float64) Variant {
	return Variant{lenOrType:VariantTypeFloat64, capOrVal: *(*int)(unsafe.Pointer(&v))}
}

func (v* Variant) Int() int {
	return v.capOrVal
}

func (v* Variant) Float64() float64 {
	return *(*float64)(unsafe.Pointer(&v.capOrVal))
}

func (v* Variant) String() (s string) {
	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenOrType
	return s
}

func (v* Variant) Bytes() (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenOrType
	dest.Cap = v.capOrVal
	return b
}
