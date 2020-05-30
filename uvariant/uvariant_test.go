package uvariant

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"unsafe"

	"github.com/tigrannajaryan/govariant/testutil"

	"github.com/stretchr/testify/assert"
)

func TestUVariant(t *testing.T) {
	fmt.Printf("Variant size=%v bytes\n", unsafe.Sizeof(Variant{}))

	v := NewEmpty()
	assert.EqualValues(t, VariantTypeEmpty, v.Type())

	b1 := []byte{1, 2, 3}
	v = NewBytes(b1)
	b2 := v.Bytes()
	assert.EqualValues(t, b1, b2)
	assert.EqualValues(t, VariantTypeBytes, v.Type())

	s1 := "abcdef"
	v = NewString(s1)
	s2 := v.String()
	assert.EqualValues(t, s1, s2)
	assert.EqualValues(t, VariantTypeString, v.Type())

	i1 := 1234
	v = NewInt(i1)
	i2 := v.Int()
	assert.EqualValues(t, i1, i2)
	assert.EqualValues(t, VariantTypeInt, v.Type())

	f1 := 1234.567
	v = NewFloat64(f1)
	f2 := v.Float64()
	assert.EqualValues(t, f1, f2)
	assert.EqualValues(t, VariantTypeFloat64, v.Type())

	v = NewMap(0)
	assert.EqualValues(t, VariantTypeMap, v.Type())
	assert.EqualValues(t, map[string]Variant{}, v.Map())

	v.Map()["k"] = NewInt(123)
	assert.EqualValues(t, map[string]Variant{"k": NewInt(123)}, v.Map())

	v = NewSlice(nil)
	assert.EqualValues(t, VariantTypeSlice, v.Type())
	assert.EqualValues(t, []Variant(nil), v.Slice())

	v = NewSlice([]Variant{NewInt(123)})
	assert.EqualValues(t, []Variant{NewInt(123)}, v.Slice())
}

func TestUVariantGC(t *testing.T) {

	var bb []*Variant

	var v1 Variant
	s1 := strconv.Itoa(1234)
	v1 = NewString(s1)
	v1 = v1

	for i := 0; i < 10000; i++ {
		s1 := strconv.Itoa(i)
		vi := new(Variant)
		*vi = NewString(s1)
		b := vi
		bb = append(bb, b)
	}

	var v Variant
	s1 = strconv.Itoa(1234)
	v = NewString(s1)

	s2 := v.String()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	bb = nil

	runtime.GC()

	runtime.ReadMemStats(&ms)

	s2 = v.String()
	s2 = s2
}

func createUVariantInt() Variant {
	for i := 0; i < 1; i++ {
		return NewInt(testutil.IntMagicVal)
	}
	return NewInt(testutil.IntMagicVal)
}

func createUVariantFloat64() Variant {
	for i := 0; i < 1; i++ {
		return NewFloat64(testutil.Float64MagicVal)
	}
	return NewFloat64(0)
}

func createUVariantString() Variant {
	for i := 0; i < 1; i++ {
		return NewString(testutil.StrMagicVal)
	}
	return NewString("def")
}

func createUVariantBytes() Variant {
	for i := 0; i < 1; i++ {
		return NewBytes(testutil.BytesMagicVal)
	}
	return NewBytes(testutil.BytesMagicVal)
}

func BenchmarkUnionVariantIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createUVariantInt()
		vi := v.Int()
		if vi != testutil.IntMagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkUnionVariantFloat64Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createUVariantFloat64()
		vf := v.Float64()
		if vf != testutil.Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkUnionVariantStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createUVariantString()
		if v.Type() == VariantTypeString {
			if v.String() == "" {
				panic("empty string")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkUnionVariantBytesTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createUVariantBytes()
		if v.Type() == VariantTypeBytes {
			if v.Bytes() == nil {
				panic("nil bytes")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkUnionVariantIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createUVariantInt()
		if v.Type() == VariantTypeInt {
			vi := v.Int()
			if vi != testutil.IntMagicVal {
				panic("invalid value")
			}
		} else {
			panic("invalid type")
		}
	}
}

func createUVariantIntSlice(n int) []Variant {
	v := make([]Variant, n)
	for i := 0; i < n; i++ {
		v[i] = NewInt(i)
	}
	return v
}

func createUVariantStringSlice(n int) []Variant {
	v := make([]Variant, n)
	for i := 0; i < n; i++ {
		v[i] = NewString("abc")
	}
	return v
}

func BenchmarkUnionVariantSliceIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createUVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.Int() < 0 {
				panic("zero int")
			}
		}
	}
}

func BenchmarkUnionVariantIntSliceTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createUVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.Type() == VariantTypeInt {
				if v.Int() < 0 {
					panic("zero int")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

func BenchmarkUnionVariantStringSliceTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createUVariantStringSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.Type() == VariantTypeString {
				if v.String() == "" {
					panic("empty string")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

func BenchmarkUnionVariantStringSliceGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createUVariantStringSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.String() == "" {
				panic("empty string")
			}
		}
	}
}

func BenchmarkUnionVariantSliceOfVariantGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := NewSlice(createUVariantStringSlice(testutil.VariantSliceSize))
		for _, v := range vv.Slice() {
			if v.String() == "" {
				panic("empty string")
			}
		}
	}
}

func createStringSlice(n int) []string {
	v := make([]string, n)
	for i := 0; i < n; i++ {
		v[i] = "abc"
	}
	return v
}

func BenchmarkStringSliceGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createStringSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v == "" {
				panic("empty string")
			}
		}
	}
}
