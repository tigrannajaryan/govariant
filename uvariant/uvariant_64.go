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

func IntVariant(v int) Variant {
	return Variant{lenAndType: VariantTypeInt, capOrVal: v}
}

func Float64Variant(v float64) Variant {
	return Variant{lenAndType: VariantTypeFloat64, capOrVal: *(*int)(unsafe.Pointer(&v))}
}

func BytesVariant(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VariantTypeBytes, capOrVal: hdr.Cap}
}
