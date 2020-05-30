// +build !386

package uvariant

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	ptr        unsafe.Pointer
	lenAndType int
	capOrVal   int
}

func NewInt(v int) Variant {
	return Variant{lenAndType: VariantTypeInt, capOrVal: v}
}

func NewFloat64(v float64) Variant {
	return Variant{lenAndType: VariantTypeFloat64, capOrVal: *(*int)(unsafe.Pointer(&v))}
}

func NewBytes(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VariantTypeBytes, capOrVal: hdr.Cap}
}

func NewSlice(v []Variant) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VariantTypeSlice, capOrVal: hdr.Cap}
}
