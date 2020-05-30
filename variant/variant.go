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
	VTypeBytes
	VTypeValueList
	VTypeKeyValueList
)

// Number of buts to use for Type field. This should be wide to fit all VType values.
const TypeFieldBitCount = 3

// Bit mask for Type part of lenAndType field.
const TypeFieldMask = (1 << TypeFieldBitCount) - 1

// Maximum length of a slice-type that can be stored in Variant. The length of Go slices
// can be at most maxint, however Variant is not able to store lengths of maxint. Len field
// in Variant uses TypeFieldBitCount bits less than int.
const MaxSliceLen = int((^uint(0))>>1) >> TypeFieldBitCount

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
	return Variant{ptr: unsafe.Pointer(hdr.Data), lenAndType: (hdr.Len << TypeFieldBitCount) | VTypeString}
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
	dest.Len = v.lenAndType >> TypeFieldBitCount
	return s
}

func (v *Variant) Bytes() (b []byte) {
	if v.Type() != VTypeBytes {
		panic("Variant is not a bytes")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> TypeFieldBitCount
	dest.Cap = int(v.capOrVal)
	return b
}

// Return the slice of values.
// Valid to call only if Type==VTypeValueList, will panic otherwise.
// Elements in the returned slice are allowed to be modified after this call returns.
func (v *Variant) ValueList() (s []Variant) {
	if v.Type() != VTypeValueList {
		panic("Variant is not a slice")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> TypeFieldBitCount
	dest.Cap = int(v.capOrVal)
	return s
}

// Return the value at specified index.
// Valid to call only if Type==VTypeValueList, will panic otherwise.
// Will panic if index is negative or exceeds current length or if the slice is nil.
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

// Return the length of contained slice-based type.
// Valid to call for String, Bytes, ValueList, KeyValueList types. For other types
// the returned value is undefined.
func (v *Variant) Len() int {
	return v.lenAndType >> TypeFieldBitCount
}

// Resize the length of contained slice-based type.
// Valid to call for String, Bytes, ValueList, KeyValueList types.
// Will panic if len is negative or exceeds the current capacity of the slice or if
// len exceeds MaxSliceLen.
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
	v.lenAndType = (v.lenAndType & TypeFieldMask) | (len << TypeFieldBitCount)
}

// Return the slice of KeyValues.
// Valid to call only if Type==VTypeKeyValueList, will panic otherwise.
// Elements in the returned slice are allowed to be modified after this call returns.
func (v *Variant) KeyValueList() (s []KeyValue) {
	if v.Type() != VTypeKeyValueList {
		panic("Variant is not a KeyValueList")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> TypeFieldBitCount
	dest.Cap = int(v.capOrVal)
	return s
}

// Return the KeyValye at specified index.
// Valid to call only if Type==VTypeKeyValueList, will panic otherwise.
// The element is returned by pointer to allowed the caller to modify the element
// if needed.
// Will panic if index is negative or exceeds current length or if the slice is nil.
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
