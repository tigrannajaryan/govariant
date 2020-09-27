// +build amd64

package variant

// This file contains Variant implementation specific to GOARCH=amd64

import (
	"reflect"
	"unsafe"
)

// Variant allows to store values of one of the VType data types.
//
// To create a Variant use one of New* functions in this package. Note that
// zero-initialized value of Variant is a valid VTypeEmpty value.
// To access the stored value call one of the *Val methods of Variant struct.
type Variant struct {
	// Pointer to the slice start for slice-based types.
	ptr unsafe.Pointer

	// Len and Type fields.
	// Type uses `typeFieldBitCount` least significant bits, Len uses the rest.
	// Len is used only for the slice-based types.
	lenAndType int

	// Capacity for slice-based types, or the value for other types. For Float64Val type
	// contains the 64 bits of the floating point value.
	capOrVal int
}

// NewInt creates a Variant of VTypeInt type.
func NewInt(v int) Variant {
	return Variant{
		lenAndType: int(VTypeInt),
		capOrVal:   v,
	}
}

// NewFloat64 creates a Variant of VTypeFloat64 type.
func NewFloat64(v float64) Variant {
	return Variant{
		lenAndType: int(VTypeFloat64),
		capOrVal:   *(*int)(unsafe.Pointer(&v)),
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
	if hdr.Len > maxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:        unsafe.Pointer(hdr.Data),
		lenAndType: (hdr.Len << typeFieldBitCount) | int(VTypeBytes),
		capOrVal:   hdr.Cap,
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
	if hdr.Len > maxSliceLen {
		panic("maximum len exceeded")
	}

	return Variant{
		ptr:        unsafe.Pointer(hdr.Data),
		lenAndType: (hdr.Len << typeFieldBitCount) | int(VTypeValueList),
		capOrVal:   hdr.Cap,
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
		ptr:        unsafe.Pointer(hdr.Data),
		lenAndType: (hdr.Len << typeFieldBitCount) | int(VTypeKeyValueList),
		capOrVal:   hdr.Cap,
	}
}
