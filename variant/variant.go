package variant

import (
	"reflect"
	"unsafe"
)

type VType int

const (
	VTypeEmpty = iota
	VTypeInt
	VTypeFloat64
	VTypeString
	VTypeMap
	VTypeBytes
	VTypeValueList
	VTypeKeyValueList
)

const TypeFieldMask = 0x07
const LenFieldBitShiftCount = 3
const MaxSliceLen = int((^uint(0))>>1) >> LenFieldBitShiftCount

type KeyValue struct {
	Key   string
	Value Variant
}

type KeyValueList []KeyValue

func (v *Variant) Type() VType {
	return VType(v.lenAndType & TypeFieldMask)
}

func NewEmpty() Variant {
	return Variant{}
}

func NewString(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << LenFieldBitShiftCount) | VTypeString}
}

func NewMap(cap int) Variant {
	m := make(map[string]Variant, cap)
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&m))
	return Variant{ptr: ptr, lenAndType: VTypeMap}
}

func (v *Variant) Int() int {
	return int(v.capOrVal)
}

func (v *Variant) Float64() float64 {
	return *(*float64)(unsafe.Pointer(&v.capOrVal))
}

func (v *Variant) String() (s string) {
	if v.Type() != VTypeString {
		panic("Variant is not a string")
	}
	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> LenFieldBitShiftCount
	return s
}

func (v *Variant) Map() map[string]Variant {
	return *(*map[string]Variant)(unsafe.Pointer(&v.ptr))
}

func (v *Variant) Bytes() (b []byte) {
	if v.Type() != VTypeBytes {
		panic("Variant is not a bytes")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> LenFieldBitShiftCount
	dest.Cap = int(v.capOrVal)
	return b
}

func (v *Variant) ValueList() (s []Variant) {
	if v.Type() != VTypeValueList {
		panic("Variant is not a slice")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> LenFieldBitShiftCount
	dest.Cap = int(v.capOrVal)
	return s
}

func (v *Variant) ValueAt(i int) Variant {
	if v.Type() != VTypeValueList {
		panic("Variant is not a slice")
	}
	if v.ptr == nil {
		panic("index of nil slice")
	}
	if i < 0 || i >= v.Len() {
		panic("out of index")
	}
	return *(*Variant)(unsafe.Pointer(uintptr(v.ptr) + uintptr(i)*unsafe.Sizeof(Variant{})))
}

func (v *Variant) Len() int {
	return v.lenAndType >> LenFieldBitShiftCount
}

func (v *Variant) Resize(len int) {
	if len < 0 {
		panic("negative len is not allowed")
	}
	if len > int(v.capOrVal) {
		panic("cannot resize beyond capacity")
	}
	if len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	v.lenAndType = (v.lenAndType & TypeFieldMask) | (len << LenFieldBitShiftCount)
}

func (v *Variant) KeyValueList() (s []KeyValue) {
	if v.Type() != VTypeKeyValueList {
		panic("Variant is not a KeyValueList")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> LenFieldBitShiftCount
	dest.Cap = int(v.capOrVal)
	return s
}

func (v *Variant) KeyValueAt(i int) *KeyValue {
	if v.Type() != VTypeKeyValueList {
		panic("Variant is not a KeyValueList")
	}
	if v.ptr == nil {
		panic("index of nil KeyValueList")
	}
	if i < 0 || i >= v.Len() {
		panic("out of index")
	}
	return (*KeyValue)(unsafe.Pointer(uintptr(v.ptr) + uintptr(i)*unsafe.Sizeof(KeyValue{})))
}
