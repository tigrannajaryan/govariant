// +build amd64

package variant

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestUVariantFieldAliasing(t *testing.T) {
	v := Variant{}

	// Ensure fields correctly alias corresponding fields of StringHeader

	// Data/ptr field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.StringHeader{}.Data), unsafe.Offsetof(v.ptr))
	assert.EqualValues(t, unsafe.Sizeof(reflect.StringHeader{}.Data), unsafe.Sizeof(v.ptr))

	// Len field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.StringHeader{}.Len), unsafe.Offsetof(v.lenAndType))
	assert.EqualValues(t, unsafe.Sizeof(reflect.StringHeader{}.Len), unsafe.Sizeof(v.lenAndType))

	// Ensure fields correctly alias corresponding fields of SliceHeader

	// Data/ptr field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.SliceHeader{}.Data), unsafe.Offsetof(v.ptr))
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Data), unsafe.Sizeof(v.ptr))

	// Len field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.SliceHeader{}.Len), unsafe.Offsetof(v.lenAndType))
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Len), unsafe.Sizeof(v.lenAndType))

	// Cap field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.SliceHeader{}.Cap), unsafe.Offsetof(v.capOrVal))
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Cap), unsafe.Sizeof(v.capOrVal))

	// Ensure float64 can correctly fit in capOrVal
	assert.EqualValues(t, unsafe.Sizeof(float64(0.0)), unsafe.Sizeof(v.capOrVal))

	// Ensure map fits in the ptr.
	m := make(map[string]Variant, 1)
	assert.EqualValues(t, unsafe.Sizeof(m), unsafe.Sizeof(v.ptr))
}
