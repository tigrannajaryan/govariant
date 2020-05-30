// +build 386

package uvariant

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/tigrannajaryan/govariant/testutil"

	"github.com/stretchr/testify/assert"
)

func TestUVariantFieldAliasing(t *testing.T) {
	v := Variant{}

	// Ensure fields correctly alias corresponding fields of StringHeader

	// Data/ptr field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.StringHeader{}.Data), unsafe.Offsetof(v.ptr))
	assert.EqualValues(t, unsafe.Sizeof(reflect.StringHeader{}.Data), unsafe.Sizeof(v.ptr))

	// Len field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.StringHeader{}.Len), unsafe.Offsetof(v.lenOrType))
	assert.EqualValues(t, unsafe.Sizeof(reflect.StringHeader{}.Len), unsafe.Sizeof(v.lenOrType))

	// Ensure fields correctly alias corresponding fields of SliceHeader

	// Data/ptr field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.SliceHeader{}.Data), unsafe.Offsetof(v.ptr))
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Data), unsafe.Sizeof(v.ptr))

	// Len field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.SliceHeader{}.Len), unsafe.Offsetof(v.lenOrType))
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Len), unsafe.Sizeof(v.lenOrType))

	// Cap field.
	assert.EqualValues(t, unsafe.Offsetof(reflect.SliceHeader{}.Cap), unsafe.Offsetof(v.capOrVal))
	assert.True(t, unsafe.Sizeof(reflect.SliceHeader{}.Cap) <= unsafe.Sizeof(v.capOrVal))
}

func createUVariantF64Float64() VariantF64 {
	for i := 0; i < 1; i++ {
		return Float64VariantF64(testutil.Float64MagicVal)
	}
	return VariantF64{}
}

//func BenchmarkUnionVariantF64Float64Get(b *testing.B) {
//	for i:=0; i<b.N; i++ {
//		v := createUVariantF64Float64()
//		vf := v.Float64()
//		if vf!=vf {
//			panic("invalid value")
//		}
//	}
//}
