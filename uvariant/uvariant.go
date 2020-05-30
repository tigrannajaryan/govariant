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

const TypeFieldMask = 0x07
const LenFieldBitShiftCount = 3

type KeyValue struct {
	key   string
	value Variant
}

func (v *Variant) Type() VariantType {
	return VariantType(v.lenAndType & TypeFieldMask)
}

func EmptyVariant() Variant {
	return Variant{}
}

func StringVariant(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VariantTypeString}
}

func MapVariant(cap int) Variant {
	m := make(map[string]Variant, cap)
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&m))
	return Variant{ptr: ptr, lenAndType: VariantTypeMap}
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
	dest.Len = v.lenAndType >> LenFieldBitShiftCount
	return s
}

func (v *Variant) Map() map[string]Variant {
	return *(*map[string]Variant)(unsafe.Pointer(&v.ptr))
}

func (v *Variant) Bytes() (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> LenFieldBitShiftCount
	dest.Cap = int(v.capOrVal)
	return b
}
