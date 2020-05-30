// +build !386

package variant

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
	return Variant{lenAndType: VTypeInt, capOrVal: v}
}

func NewFloat64(v float64) Variant {
	return Variant{lenAndType: VTypeFloat64, capOrVal: *(*int)(unsafe.Pointer(&v))}
}

func NewBytes(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VTypeBytes, capOrVal: hdr.Cap}
}

func NewValueList(v []Variant) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VTypeValueList, capOrVal: hdr.Cap}
}

func NewKeyValueList(cap int) Variant {
	v := make([]KeyValue, 0, cap)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VTypeKeyValueList, capOrVal: hdr.Cap}
}
