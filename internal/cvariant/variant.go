// +build amd64

package cvariant

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// Type of a value stored in Variant.
type Type int

// Possible value types that can be stored in Variant.
const (
	// Empty or no value. The default state of zero-initialized Variant.
	TypeEmpty Type = iota

	// An int number.
	TypeInt

	// A float64 number.
	TypeFloat64

	// A string.
	TypeString

	// A []byte slice.
	TypeBytes

	// A list of Variant.
	TypeValueList

	// A list of KeyValue.
	TypeKeyValueList
)

// Number of bits to use for Type field. This should be wide enough to fit all Type values.
const TypeFieldBitCount = 3

const totalBitCount = 8 * unsafe.Sizeof(int64(0))
const lenAndCapBitCount = totalBitCount - TypeFieldBitCount

const capFieldBitCount = lenAndCapBitCount / 2
const capFieldShiftCount = TypeFieldBitCount
const capFieldMask = (1 << capFieldBitCount) - 1

const lenFieldBitCount = lenAndCapBitCount - capFieldBitCount
const lenFieldShiftCount = TypeFieldBitCount + capFieldBitCount
const lenFieldMask = (1 << lenFieldBitCount) - 1

// Bit mask for Type part of lenAndType field.
const typeFieldMask = (1 << TypeFieldBitCount) - 1

//const lenFieldShiftCount = TypeFieldBitCount

// Maximum length of a slice-type that can be stored in Variant. The length of Go slices
// can be at most maxint, however Variant is not able to store lengths of maxint. Len field
// in Variant uses lenFieldShiftCount bits less than int, i.e. the maximum length of a slice
// stored in Variant is maxint / (2^lenFieldShiftCount), which we calculate below.
const MaxSliceLen = int((^uint(0))>>1) >> lenFieldShiftCount

// A slice of Type values which is used as a marker of the type to which the Variant's
// ptr field points to for non pointer types.
var intTypeMarker = TypeInt
var floatTypeMarker = TypeFloat64

// KeyValue is an element that is used for TypeKeyValueList storage.
type KeyValue struct {
	Key   string
	Value Variant
}

// Type returns the type of the currently stored value.
func (v *Variant) Type() Type {
	if v.ptr != nil {
		switch v.ptr {
		case unsafe.Pointer(&intTypeMarker):
			return TypeInt
		case unsafe.Pointer(&floatTypeMarker):
			return TypeFloat64
		}
	}

	return Type(v.bits & typeFieldMask)
}

// NewEmpty creates a Variant of TypeEmpty type.
func NewEmpty() Variant {
	return Variant{}
}

// NewString creates a Variant of TypeString type.
func NewString(v string) Variant {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:  unsafe.Pointer(hdr.Data),
		bits: uint(hdr.Len<<lenFieldShiftCount) | uint(TypeString),
	}
}

// NewStringFromBytes creates a Variant of TypeString type from a slice of bytes
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
		ptr:  unsafe.Pointer(hdr.Data),
		bits: uint(hdr.Len<<lenFieldShiftCount) | uint(TypeString),
	}
}

// IntVal returns the stored int value.
// The returned value is undefined if the Variant type is not TypeInt.
func (v *Variant) IntVal() int {
	return int(v.bits)
}

// Float64Val returns the stored float64 value.
// The returned value is undefined if the Variant type is not TypeFloat64.
func (v *Variant) Float64Val() float64 {
	return *(*float64)(unsafe.Pointer(&v.bits))
}

// StringVal returns the stored string value.
// Will panic if the Variant type is not TypeString.
func (v *Variant) StringVal() (s string) {
	switch v.ptr {
	case unsafe.Pointer(&intTypeMarker):
		fallthrough
	case unsafe.Pointer(&floatTypeMarker):
		panic("Variant is not a TypeString")
	}

	if Type(v.bits&typeFieldMask) != TypeString {
		panic("Variant is not a TypeString")
	}

	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = int(v.bits >> lenFieldShiftCount)
	return s
}

// Bytes returns the stored byte slice.
// Will panic if the Variant type is not TypeBytes.
func (v *Variant) Bytes() (b []byte) {
	switch v.ptr {
	case unsafe.Pointer(&intTypeMarker):
		fallthrough
	case unsafe.Pointer(&floatTypeMarker):
		panic("Variant is not a TypeBytes")
	}

	if Type(v.bits&typeFieldMask) != TypeBytes {
		panic("Variant is not a TypeBytes")
	}

	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = int(v.bits >> lenFieldShiftCount)
	dest.Cap = int((v.bits >> capFieldShiftCount) & capFieldMask)
	return b
}

// ValueList returns the slice of stored Variant values.
// Elements in the returned slice are allowed to be modified after this call returns.
// Will panic if the Variant type is not TypeValueList.
func (v *Variant) ValueList() (s []Variant) {
	switch v.ptr {
	case unsafe.Pointer(&intTypeMarker):
		fallthrough
	case unsafe.Pointer(&floatTypeMarker):
		panic("Variant is not a TypeValueList")
	}

	if Type(v.bits&typeFieldMask) != TypeValueList {
		panic("Variant is not a TypeValueList")
	}

	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = int(v.bits >> lenFieldShiftCount)
	dest.Cap = int((v.bits >> capFieldShiftCount) & capFieldMask)
	return s
}

// ValueAt returns the value at the specified index.
//
// Valid to call only if Variant type is TypeValueList otherwise will panic.
// Will panic if index is negative or is greater or equal the current length.
func (v *Variant) ValueAt(index int) Variant {
	switch v.ptr {
	case nil:
		fallthrough
	case unsafe.Pointer(&intTypeMarker):
		fallthrough
	case unsafe.Pointer(&floatTypeMarker):
		panic("Variant is not a TypeValueList or is empty")
	}

	if Type(v.bits&typeFieldMask) != TypeValueList {
		panic("Variant is not a TypeValueList")
	}

	if index < 0 || index >= v.Len() {
		panic("index out of bounds")
	}

	return *(*Variant)(unsafe.Pointer(uintptr(v.ptr) + uintptr(index)*unsafe.Sizeof(Variant{})))
}

// Len returns the length of contained slice-based type.
// Valid to call for TypeString, TypeBytes, TypeValueList, TypeKeyValueList types.
// For other types the returned value is undefined.
func (v *Variant) Len() int {
	return int(v.bits >> lenFieldShiftCount)
}

// Resize the length of contained slice-based type.
// Valid to call for TypeString, TypeBytes, TypeValueList, TypeKeyValueList types.
// Will panic for other types.
// Will panic if len is negative or exceeds the current capacity of the slice or if
// len exceeds MaxSliceLen.
func (v *Variant) Resize(len int) {
	switch v.Type() {
	case TypeEmpty, TypeInt, TypeFloat64:
		panic(fmt.Sprintf("Cannot resize Variant type %d", v.Type()))
	}

	if len < 0 {
		panic("negative len is not allowed")
	}
	if len > int((v.bits>>capFieldShiftCount)&capFieldMask) {
		panic("cannot resize beyond capacity")
	}
	if len > MaxSliceLen {
		panic("maximum len exceeded")
	}
	v.bits = (v.bits & (typeFieldMask | (capFieldMask << capFieldShiftCount))) | (uint(len) << lenFieldShiftCount)
}

// KeyValueList return the slice of stored KeyValue.
// Valid to call only if Type==TypeKeyValueList otherwise will panic.
// Elements in the returned slice are allowed to be modified after this call returns.
// Such modification will affect the KeyValue stored in this Variant since returned
// slice is a reference type.
func (v *Variant) KeyValueList() (s []KeyValue) {
	switch v.ptr {
	case unsafe.Pointer(&intTypeMarker):
		fallthrough
	case unsafe.Pointer(&floatTypeMarker):
		panic("Variant is not a TypeKeyValueList")
	}

	if Type(v.bits&typeFieldMask) != TypeKeyValueList {
		panic("Variant is not a TypeKeyValueList")
	}

	dest := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = int(v.bits >> lenFieldShiftCount)
	dest.Cap = int((v.bits >> capFieldShiftCount) & capFieldMask)
	return s
}

// KeyValueAt returns the KeyValue at the specified index.
//
// Valid to call only if Variant type is TypeKeyValueList otherwise will panic.
// The element is returned by pointer to allow the caller to modify the element
// by assigning to it if needed.
// Will panic if index is negative or is greater or equal the current length.
func (v *Variant) KeyValueAt(index int) *KeyValue {
	switch v.ptr {
	case nil:
		fallthrough
	case unsafe.Pointer(&intTypeMarker):
		fallthrough
	case unsafe.Pointer(&floatTypeMarker):
		panic("Variant is not a TypeKeyValueList or is empty")
	}

	if Type(v.bits&typeFieldMask) != TypeKeyValueList {
		panic("Variant is not a TypeKeyValueList")
	}

	if index < 0 || index >= v.Len() {
		panic("index out of bounds")
	}

	return (*KeyValue)(unsafe.Pointer(uintptr(v.ptr) + uintptr(index)*unsafe.Sizeof(KeyValue{})))
}

// String returns a human readable string representation of the stored value.
func (v Variant) String() string {
	switch v.Type() {
	case TypeEmpty:
		return ""
	case TypeInt:
		return strconv.Itoa(v.IntVal())
	case TypeFloat64:
		return strconv.FormatFloat(v.Float64Val(), 'g', -1, 64)
	case TypeString:
		return fmt.Sprintf("%q", v.StringVal())
	case TypeBytes:
		return fmt.Sprintf("0x%X", v.Bytes())
	case TypeValueList:
		var strs []string
		for _, e := range v.ValueList() {
			strs = append(strs, e.String())
		}
		return "[" + strings.Join(strs, ",") + "]"
	case TypeKeyValueList:
		var strs []string
		for _, e := range v.KeyValueList() {
			strs = append(strs, fmt.Sprintf("%q:%s", e.Key, e.Value.String()))
		}
		return "{" + strings.Join(strs, ",") + "}"
	}
	panic("invalid Variant type")
}
