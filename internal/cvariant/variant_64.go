// +build amd64

package cvariant

// This file contains Variant implementation specific to GOARCH=amd64

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	// Pointer to the slice start for slice-based types.
	ptr unsafe.Pointer

	bits uint
}

// NewInt creates a Variant of VTypeInt type.
func NewInt(v int) Variant {
	return Variant{
		ptr:  unsafe.Pointer(&intTypeMarker),
		bits: uint(v),
	}
}

// NewFloat64 creates a Variant of VTypeFloat64 type.
func NewFloat64(v float64) Variant {
	return Variant{
		ptr:  unsafe.Pointer(&floatTypeMarker),
		bits: *(*uint)(unsafe.Pointer(&v)),
	}
}

// NewBytes creates a Variant of VTypeBytes type and initializes it with the specified
// slice of bytes.
//
// This function does not copy the slice. The Variant will point to
// the same slice that is pointed to by the parameter v. Any changes made to the bytes
// in the slice v will be also reflected in the byte slice stored in this Variant.
func NewBytes(v []byte) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:  unsafe.Pointer(hdr.Data),
		bits: uint(hdr.Len<<lenFieldShiftCount) | uint(hdr.Cap<<capFieldShiftCount) | uint(VTypeBytes),
	}
}

// NewValueList creates a Variant of VTypeValueList type and initializes it with the
// specified slice of Variants.
//
// This function does not copy the slice. The Variant will point to the same slice that
// is pointed to by the parameter v. Any changes made to the elements in the slice v
// will be also reflected in the list stored in this Variant.
func NewValueList(v []Variant) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:  unsafe.Pointer(hdr.Data),
		bits: uint(hdr.Len<<lenFieldShiftCount) | uint(hdr.Cap<<capFieldShiftCount) | uint(VTypeValueList),
	}
}

// NewKeyValueList creates a Variant of VTypeKeyValueList type and initializes it with the
// specified slice of KeyValues.
//
// This function does not copy the slice. The Variant will point to the same slice that
// is pointed to by the parameter v. Any changes made to the elements in the slice v
// will be also reflected in the list stored in this Variant.
func NewKeyValueList(v []KeyValue) Variant {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))

	return Variant{
		ptr:  unsafe.Pointer(hdr.Data),
		bits: uint(hdr.Len<<lenFieldShiftCount) | uint(hdr.Cap<<capFieldShiftCount) | uint(VTypeKeyValueList),
	}
}
