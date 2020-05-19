package main

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestUVariant(t*testing.T) {
	fmt.Printf("Variant size=%v\n", unsafe.Sizeof(Variant{}))

	b1 := []byte{1,2,3}
	v := BytesVariant(b1)
	b2 := v.Bytes()
	assert.EqualValues(t, b1, b2)
	assert.EqualValues(t, VariantTypeBytes, v.Type())

	s1 := "abcdef"
	v = StringVariant(s1)
	s2 := v.String()
	assert.EqualValues(t, s1, s2)
	assert.EqualValues(t, VariantTypeString, v.Type())

	i1 := 1234
	v = IntVariant(i1)
	i2 := v.Int()
	assert.EqualValues(t, i1, i2)
	assert.EqualValues(t, VariantTypeInt, v.Type())

	f1 := 1234.567
	v = Float64Variant(f1)
	f2 := v.Float64()
	assert.EqualValues(t, f1, f2)
	assert.EqualValues(t, VariantTypeFloat64, v.Type())

	//assert.EqualValues(t, 8, unsafe.Sizeof(int(123)))
}

func createUVariantInt() Variant {
	for i :=0; i<1; i++ {
		return IntVariant(IntMagicVal)
	}
	return IntVariant(0)
}

func createUVariantFloat64() Variant {
	for i :=0; i<1; i++ {
		return Float64Variant(Float64MagicVal)
	}
	return Float64Variant(0)
}

func createUVariantString() Variant {
	for i :=0; i<1; i++ {
		return StringVariant(StrMagicVal)
	}
	return StringVariant("def")
}

func createUVariantBytes() Variant {
	for i :=0; i<1; i++ {
		return BytesVariant(BytesMagicVal)
	}
	return BytesVariant(BytesMagicVal)
}

func BenchmarkUVariantIntGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createUVariantInt()
		vi := v.Int()
		if vi!=vi {
			panic("invalid value")
		}
	}
}

func BenchmarkUVariantFloat64Get(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createUVariantFloat64()
		vf := v.Float64()
		if vf!= Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkUVariantStringTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createUVariantString()
		if v.Type()==VariantTypeString {
			if v.String()=="" {
				panic("empty string")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkUVariantBytesTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createUVariantBytes()
		if v.Type()== VariantTypeBytes {
			if v.Bytes()==nil {
				panic("nil bytes")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkUVariantIntTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createUVariantInt()
		if v.Type()==VariantTypeInt {
			vi := v.Int()
			if vi!=vi {
				panic("invalid value")
			}
		} else {
			panic("invalid type")
		}
	}
}

func createUVariantIntSlice(n int) []Variant {
	v := make([]Variant, n)
	for i :=0; i<n; i++ {
		v[i] = IntVariant(i)
	}
	return v
}

func createUVariantStringSlice(n int) []Variant {
	v := make([]Variant, n)
	for i :=0; i<n; i++ {
		v[i] = StringVariant("abc")
	}
	return v
}

func BenchmarkUVariantSliceIntGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		vv := createUVariantIntSlice(VariantSliceSize)
		for _,v := range vv {
			if v.Int()<0 {
				panic("zero int")
			}
		}
	}
}

func BenchmarkUVariantSliceIntTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		vv := createUVariantIntSlice(VariantSliceSize)
		for _,v := range vv {
			if v.Type()==VariantTypeInt {
				if v.Int()<0 {
					panic("zero int")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

func BenchmarkUVariantSliceStringTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		vv := createUVariantStringSlice(VariantSliceSize)
		for _,v := range vv {
			if v.Type()==VariantTypeString {
				if v.String()=="" {
					panic("empty string")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

