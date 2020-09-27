// +build 386

package variant

// This file contains Variant implementation specific to GOARCH=386

import (
	"reflect"
	"unsafe"
)

type Variant struct {
	// Pointer to the slice start for slice-based types.
	ptr unsafe.Pointer

	// Len and Type fields.
	// Type uses `typeFieldBitCount` least significant bits, Len uses the rest.
	// Len is used only for the slice-based types.
	lenAndType int

	// Capacity for slice-based types, or the value for other types. For Float64Val type
	// contains the 64 bits of the floating point value.
	capOrVal int64
}

// NewInt creates a Variant of VTypeInt type.
func NewInt(v int) Variant {
	return Variant{
		lenAndType: int(VTypeInt),
		capOrVal:   int64(v),
	}
}

// NewFloat64 creates a Variant of VTypeFloat64 type.
func NewFloat64(v float64) (r Variant) {
	r.lenAndType = int(VTypeFloat64)
	*(*float64)(unsafe.Pointer(&r.capOrVal)) = v
	return r
}

// NewBytes creates a Variant of VTypeBytes type and initializes it with the specified
// slice of bytes.
//
// This function does not copy the slice. the Variant will point to
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
		capOrVal:   int64(hdr.Cap),
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
		capOrVal:   int64(hdr.Cap),
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
		capOrVal:   int64(hdr.Cap),
	}
}
