// Package variant implements a Variant type that can store a value of one of the
// predefined types (see VType for the list of supported types).
//
// To use a Variant first create and store a value in it and then check the stored value
// type and read it. For example:
//
// v := NewInt(123)
// if v.Type() == VTypeInt {
//   x := v.IntVal() // x is now int value 123.
// }
//
// Variant uses an efficient data structure that is small and fast to operate on.
// On 64 bit systems the size of Variant is 24 bytes for scalar types (such as int or
// float64) plus any necessary additional data required by variable-sized types (String,
// List, etc).
//
// To maximize the performance Variant functions do not return errors. All functions define
// clear contracts that describe in which case the calls are valid. In such cases it is
// guaranteed that no errors occur. If the caller violates the contract and it results
// in erroneous situation one of the 2 things will happen (and such behavior is
// documented for each function):
//
// - Function will return an undefined value.
// - Function will panic.
//
// Typically panics are used to mimic the behavior of builtin Go types. For example
// accessing an element of a list using an index that is out of bounds will result in a
// panic similarly to how it will panic for Go slices.
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
const TypeFieldMask = (1 << TypeFieldBitCount) - 1

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
	return VType(v.lenAndType & TypeFieldMask)
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

// IntVal returns the stored int value.
// The returned value is undefined if the Variant type is not equal to VTypeInt.
func (v *Variant) IntVal() int {
	return int(v.capOrVal)
}

// Float64Val returns the stored float64 value.
// The returned value is undefined if the Variant type is not equal to VTypeFloat64.
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
// Elements in the returned slice are allowed to be modified after this call returns.
// Will panic if the Variant type is not VTypeValueList.
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

// ValueAt returns the value at specified index.
// Will panic if the Variant type is not VTypeValueList.
// Will panic if index is negative or exceeds current (length-1).
func (v *Variant) ValueAt(i int) Variant {
	if v.Type() != VTypeValueList {
		panic("Variant is not a VTypeValueList")
	}
	if v.ptr == nil {
		panic("index of empty VTypeValueList")
	}
	if i < 0 || i >= v.Len() {
		panic("out of index")
	}
	return *(*Variant)(unsafe.Pointer(uintptr(v.ptr) + uintptr(i)*unsafe.Sizeof(Variant{})))
}

// Len returns the length of contained slice-based type.
// Valid to call for VTypeString, VTypeBytes, VTypeValueList, VTypeKeyValueList types.
// For other types the returned value is undefined.
func (v *Variant) Len() int {
	return v.lenAndType >> TypeFieldBitCount
}

// Resize the length of contained slice-based type.
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
	v.lenAndType = (v.lenAndType & TypeFieldMask) | (len << TypeFieldBitCount)
}

// KeyValueList return the slice of stored KeyValue.
// Valid to call only if Type==VTypeKeyValueList otherwise will panic.
// Elements in the returned slice are allowed to be modified after this call returns.
// Such modification will affect the KeyValue stored in this Variant since returned
// slice is a reference type.
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

// KeyValueAt returns the KeyValue at specified index.
// Valid to call only if Type==VTypeKeyValueList otherwise will panic.
// The element is returned by pointer to allow the caller to modify the element
// by assigning to it if needed.
// Will panic if index is negative or exceeds current (length-1).
func (v *Variant) KeyValueAt(index int) *KeyValue {
	if v.Type() != VTypeKeyValueList {
		panic("Variant is not a VTypeKeyValueList")
	}
	if v.ptr == nil {
		panic("index of empty VTypeKeyValueList")
	}
	if index < 0 || index >= v.Len() {
		panic("out of index")
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
	panic("Invalid Variant type")
}
