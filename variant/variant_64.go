// +build amd64

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
	capOrVal int
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
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeBytes, capOrVal: hdr.Cap}
}

func NewValueList(v []Variant) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeValueList, capOrVal: hdr.Cap}
}

func NewKeyValueList(cap int) Variant {
	v := make([]KeyValue, 0, cap)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeKeyValueList, capOrVal: hdr.Cap}
}
