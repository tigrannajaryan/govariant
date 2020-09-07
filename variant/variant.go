package variant

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// Type of a value stored in Variant.
type VType int

// Possible value types that can be stored in Variant.
const (
	// Empty or no value. The default state of zero-initialized Variant.
	VTypeEmpty VType = iota

	// An int number.
	VTypeInt

	// A float64 number.
	VTypeFloat64

	// A string.
	VTypeString

	// A []byte slice.
	VTypeBytes

	// A list of Variant.
	VTypeValueList

	// A list of KeyValue.
	VTypeKeyValueList
)

// Number of bits to use for Type field. This should be wide enough to fit all VType values.
const TypeFieldBitCount = 3

// Bit mask for Type part of lenAndType field.
const typeFieldMask = (1 << TypeFieldBitCount) - 1

// Maximum length of a slice-type that can be stored in Variant. The length of Go slices
// can be at most maxint, however Variant is not able to store lengths of maxint. Len field
// in Variant uses TypeFieldBitCount bits less than int, i.e. the maximum length of a slice
// stored in Variant is maxint / (2^TypeFieldBitCount), which we calculate below.
const MaxSliceLen = int((^uint(0))>>1) >> TypeFieldBitCount

// KeyValue is an element that is used for VTypeKeyValueList storage.
type KeyValue struct {
	Key   string
	Value Variant
}

// Type returns the type of the currently stored value.
func (v *Variant) Type() VType {
	return VType(v.lenAndType & typeFieldMask)
}

// NewEmpty creates a Variant of VTypeEmpty type.
func NewEmpty() Variant {
	return Variant{}
}

// NewString creates a Variant of VTypeString type.
func NewString(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:        unsafe.Pointer(hdr.Data),
		lenAndType: (hdr.Len << TypeFieldBitCount) | int(VTypeString),
	}
}

// NewStringFromBytes creates a Variant of VTypeString type from a slice of bytes
// that represent the string.
//
// WARNING: the string stored inside this Variant will be aliased in the memory and will
// share its storage with the byte slice provided. This means any changes to the bytes
// in the slice will also modify the string in this Variant.
//
// This function should be only used when it is guaranteed that the bytes
// in the slice will not be modified or when the immutability of the string
// stored inside this Variant is not required. In such cases NewStringFromBytes(v)
// provides significant performance advantage over NewString(string(v)) call,
// which will create a copy of byte slice 'v'.
func NewStringFromBytes(v []byte) (r Variant) {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:        unsafe.Pointer(hdr.Data),
		lenAndType: (hdr.Len << TypeFieldBitCount) | int(VTypeString),
	}
}

// IntVal returns the stored int value.
// The returned value is undefined if the Variant type is not VTypeInt.
func (v *Variant) IntVal() int {
	return int(v.capOrVal)
}

// Float64Val returns the stored float64 value.
// The returned value is undefined if the Variant type is not VTypeFloat64.
func (v *Variant) Float64Val() float64 {
	return *(*float64)(unsafe.Pointer(&v.capOrVal))
}

// StringVal returns the stored string value.
// Will panic if the Variant type is not VTypeString.
func (v *Variant) StringVal() (s string) {
	if v.Type() != VTypeString {
		panic("Variant is not a VTypeString")
	}
	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> TypeFieldBitCount
	return s
}

// Bytes returns the stored byte slice.
// Will panic if the Variant type is not VTypeBytes.
func (v *Variant) Bytes() (b []byte) {
	if v.Type() != VTypeBytes {
		panic("Variant is not a VTypeBytes")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> TypeFieldBitCount
	dest.Cap = int(v.capOrVal)
	return b
}

// ValueList returns the slice of stored Variant values.
//
// Elements in the returned slice are allowed to be modified after this call returns.
// Will panic if the Variant type is not VTypeValueList.
//
// It is recommended to use this function for iteration over the list, e.g.
// 		for i, e := range v.ValueList() {
//			// Do something with item e
//		}
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

// ValueAt returns the value at the specified index.
//
// Valid to call only if Variant type is VTypeValueList otherwise will panic.
// Will panic if index is negative or is greater or equal the current length.
//
// ValueAt() and Len() can be used to iterate over the list using a for loop,
// however instead it is recommended to call ValueList() and use for-range
// loop over the returned value (the later approach is faster and safer). See
// ValueList() for an example.
func (v *Variant) ValueAt(i int) Variant {
	if v.Type() != VTypeValueList {
		panic("Variant is not a VTypeValueList")
	}
	if v.ptr == nil {
		panic("index of empty VTypeValueList")
	}
	if i < 0 || i >= v.Len() {
		panic("index out of bounds")
	}
	return *(*Variant)(unsafe.Pointer(uintptr(v.ptr) + uintptr(i)*unsafe.Sizeof(Variant{})))
}

// Len returns the length of contained slice-based type.
//
// Valid to call for VTypeString, VTypeBytes, VTypeValueList, VTypeKeyValueList types.
// For other types the returned value is undefined.
func (v *Variant) Len() int {
	return v.lenAndType >> TypeFieldBitCount
}

// Resize the length of contained slice-based type.
//
// Valid to call for VTypeString, VTypeBytes, VTypeValueList, VTypeKeyValueList types.
// Will panic for other types.
// Will panic if len is negative or exceeds the current capacity of the slice or if
// len exceeds MaxSliceLen.
func (v *Variant) Resize(len int) {
	switch v.Type() {
	case VTypeEmpty, VTypeInt, VTypeFloat64:
		panic(fmt.Sprintf("Cannot resize Variant type %d", v.Type()))
	}

	if len < 0 {
		panic("negative len is not allowed")
	}
	if len > int(v.capOrVal) {
		panic("cannot resize beyond capacity")
	}
	if len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	v.lenAndType = (v.lenAndType & typeFieldMask) | (len << TypeFieldBitCount)
}

// KeyValueList return the slice of stored KeyValue.
//
// Valid to call only if Type==VTypeKeyValueList otherwise will panic.
// Elements in the returned slice are allowed to be modified after this call returns.
// Such modification will affect the KeyValue stored in this Variant since returned
// slice is a reference type.
//
// It is recommended to use this function for iteration over the list, e.g.
// 		for i, kv := range v.KeyValueList() {
//			// Do something with item kv.Key and kv.Value
//		}
func (v *Variant) KeyValueList() (s []KeyValue) {
	if v.Type() != VTypeKeyValueList {
		panic("Variant is not a VTypeKeyValueList")
	}
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = v.lenAndType >> TypeFieldBitCount
	dest.Cap = int(v.capOrVal)
	return s
}

// KeyValueAt returns the KeyValue at the specified index.
//
// Valid to call only if Variant type is VTypeKeyValueList otherwise will panic.
// The element is returned by pointer to allow the caller to modify the element
// by assigning to it if needed.
// Will panic if index is negative or is greater or equal the current length.
//
// KeyValueAt() and Len() can be used to iterate over the list using a for loop,
// however instead it is recommended to call KeyValueList() and use for-range
// loop over the returned value (the later approach is faster and safer). See
// KeyValueList() for an example.
func (v *Variant) KeyValueAt(index int) *KeyValue {
	if v.Type() != VTypeKeyValueList {
		panic("Variant is not a VTypeKeyValueList")
	}
	if v.ptr == nil {
		panic("index of empty VTypeKeyValueList")
	}
	if index < 0 || index >= v.Len() {
		panic("index out of bounds")
	}
	return (*KeyValue)(unsafe.Pointer(uintptr(v.ptr) + uintptr(index)*unsafe.Sizeof(KeyValue{})))
}

// String returns a human readable string representation of the stored value.
func (v Variant) String() string {
	switch v.Type() {
	case VTypeEmpty:
		return ""
	case VTypeInt:
		return strconv.Itoa(v.IntVal())
	case VTypeFloat64:
		return strconv.FormatFloat(v.Float64Val(), 'g', -1, 64)
	case VTypeString:
		return fmt.Sprintf("%q", v.StringVal())
	case VTypeBytes:
		return fmt.Sprintf("0x%X", v.Bytes())
	case VTypeValueList:
		var strs []string
		for _, e := range v.ValueList() {
			strs = append(strs, e.String())
		}
		return "[" + strings.Join(strs, ",") + "]"
	case VTypeKeyValueList:
		var strs []string
		for _, e := range v.KeyValueList() {
			strs = append(strs, fmt.Sprintf("%q:%s", e.Key, e.Value.String()))
		}
		return "{" + strings.Join(strs, ",") + "}"
	}
	panic("invalid Variant type")
}
