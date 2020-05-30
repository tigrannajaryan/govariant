package uvariant

import (
	"reflect"
	"unsafe"
)

type VariantType int

const (
	VariantTypeEmpty = iota
	VariantTypeInt
	VariantTypeFloat64
	VariantTypeString
	VariantTypeMap
	VariantTypeBytes
	VariantTypeKVList
)

type KeyValue struct {
	key   string
	value Variant
}

func (v *Variant) Type() VariantType {
	if v.ptr == nil {
		// Primitive type, no pointer.
		return VariantType(v.lenOrType)
	}

	// Pointer type.

	if v.lenOrType < 0 {
		return VariantType(-v.lenOrType)
	}

	if v.capOrVal == -1 {
		return VariantTypeString
	}

	return VariantTypeBytes
}

func EmptyVariant() Variant {
	return Variant{}
}

func StringVariant(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenOrType: hdr.Len, capOrVal: -1}
}

func MapVariant(cap int) Variant {
	m := make(map[string]Variant, cap)
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&m))
	return Variant{ptr: ptr, lenOrType: -VariantTypeMap}
}

func (v *Variant) Int() int {
	return int(v.capOrVal)
}

func (v *Variant) Float64() float64 {
	return *(*float64)(unsafe.Pointer(&v.capOrVal))
}

func (v *Variant) String() (s string) {
	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenOrType
	return s
}

func (v *Variant) Map() map[string]Variant {
	return *(*map[string]Variant)(unsafe.Pointer(&v.ptr))
}

func (v *Variant) Bytes() (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenOrType
	dest.Cap = int(v.capOrVal)
	return b
}
