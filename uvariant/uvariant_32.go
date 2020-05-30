// +build 386

package uvariant

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	ptr        unsafe.Pointer
	lenAndType int
	capOrVal   int64 // cap of bytes, int or float value.
}

func NewInt(v int) Variant {
	return Variant{lenAndType: VariantTypeInt, capOrVal: int64(v)}
}

func NewFloat64(v float64) (r Variant) {
	r.lenAndType = VariantTypeFloat64
	*(*float64)(unsafe.Pointer(&r.capOrVal)) = v
	return r
}

func NewBytes(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VariantTypeBytes, capOrVal: int64(hdr.Cap)}
}

func NewSlice(v []Variant) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VariantTypeSlice, capOrVal: int64(hdr.Cap)}
}
