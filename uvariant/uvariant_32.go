// +build 386

package uvariant

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	ptr       unsafe.Pointer
	lenOrType int
	capOrVal  int64 // cap of bytes, int or float value.
}

func IntVariant(v int) Variant {
	return Variant{lenOrType: VariantTypeInt, capOrVal: int64(v)}
}

func Float64Variant(v float64) (r Variant) {
	r.lenOrType = VariantTypeFloat64
	*(*float64)(unsafe.Pointer(&r.capOrVal)) = v
	return r
}

func BytesVariant(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType: hdr.Len, capOrVal: int64(hdr.Cap)}
}
