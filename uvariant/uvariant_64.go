// +build !386

package uvariant

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	ptr       unsafe.Pointer
	lenOrType int
	capOrVal  int
}

func IntVariant(v int) Variant {
	return Variant{lenOrType: VariantTypeInt, capOrVal: int(v)}
}

func Float64Variant(v float64) Variant {
	return Variant{lenOrType: VariantTypeFloat64, capOrVal: *(*int)(unsafe.Pointer(&v))}
}

func BytesVariant(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType: hdr.Len, capOrVal: hdr.Cap}
}
