// +build 386

package variant

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	// Pointer to the slice start for slice-based types.
	ptr unsafe.Pointer

	// Len and Type fields.
	// Type uses `TypeFieldBitCount` least significant bits, Len uses the rest.
	// Len is used only for the slice-based types.
	lenAndType int

	// Capacity for slice-based types, or the value for other types. For Float64 type
	// contains the 64 bits of the floating point value.
	capOrVal int64
}

func NewInt(v int) Variant {
	return Variant{lenAndType: VTypeInt, capOrVal: int64(v)}
}

func NewFloat64(v float64) (r Variant) {
	r.lenAndType = VTypeFloat64
	*(*float64)(unsafe.Pointer(&r.capOrVal)) = v
	return r
}

func NewBytes(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeBytes, capOrVal: int64(hdr.Cap)}
}

func NewValueList(v []Variant) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeValueList, capOrVal: int64(hdr.Cap)}
}

func NewKeyValueList(cap int) Variant {
	v := make([]KeyValue, 0, cap)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeKeyValueList, capOrVal: int64(hdr.Cap)}
}
