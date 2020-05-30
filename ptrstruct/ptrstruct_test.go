package ptrstruct

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/tigrannajaryan/govariant/testutil"

	"github.com/tigrannajaryan/govariant/uvariant"

	"github.com/stretchr/testify/assert"
)

func TestUVariant(t *testing.T) {
	fmt.Printf("Variant size=%v bytes\n", unsafe.Sizeof(Variant{}))

	b1 := []byte{1, 2, 3}
	v := BytesVariant(b1)
	b2 := v.Bytes()
	assert.EqualValues(t, b1, b2)
	assert.EqualValues(t, uvariant.VariantTypeBytes, v.Type())

	s1 := "abcdef"
	v = StringVariant(s1)
	s2 := v.String()
	assert.EqualValues(t, s1, s2)
	assert.EqualValues(t, uvariant.VariantTypeString, v.Type())

	i1 := 1234
	v = IntVariant(i1)
	i2 := v.Int()
	assert.EqualValues(t, i1, i2)
	assert.EqualValues(t, uvariant.VariantTypeInt, v.Type())

	f1 := 1234.567
	v = Float64Variant(f1)
	f2 := v.Float64()
	assert.EqualValues(t, f1, f2)
	assert.EqualValues(t, uvariant.VariantTypeFloat64, v.Type())

	//assert.EqualValues(t, 8, unsafe.Sizeof(int(123)))
}

func createVariantInt() *Variant {
	for i := 0; i < 1; i++ {
		return IntVariant(testutil.IntMagicVal)
	}
	return IntVariant(0)
}

func createVariantFloat64() *Variant {
	for i := 0; i < 1; i++ {
		return Float64Variant(testutil.Float64MagicVal)
	}
	return Float64Variant(0)
}

func createVariantString() *Variant {
	for i := 0; i < 1; i++ {
		return StringVariant(testutil.StrMagicVal)
	}
	return StringVariant("def")
}

func createVariantBytes() *Variant {
	for i := 0; i < 1; i++ {
		return BytesVariant(testutil.BytesMagicVal)
	}
	return BytesVariant(testutil.BytesMagicVal)
}

func BenchmarkVariantIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		vi := v.Int()
		if vi != vi {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantFloat64Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantFloat64()
		vf := v.Float64()
		if vf != testutil.Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantString()
		if v.Type() == uvariant.VariantTypeString {
			if v.String() == "" {
				panic("empty string")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkVariantBytesTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantBytes()
		if v.Type() == uvariant.VariantTypeBytes {
			if v.Bytes() == nil {
				panic("nil bytes")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkVariantIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		if v.Type() == uvariant.VariantTypeInt {
			vi := v.Int()
			if vi != vi {
				panic("invalid value")
			}
		} else {
			panic("invalid type")
		}
	}
}

func createVariantIntSlice(n int) []*Variant {
	v := make([]*Variant, n)
	for i := 0; i < n; i++ {
		v[i] = IntVariant(i)
	}
	return v
}

func createVariantStringSlice(n int) []*Variant {
	v := make([]*Variant, n)
	for i := 0; i < n; i++ {
		v[i] = StringVariant("abc")
	}
	return v
}

func BenchmarkVariantSliceIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.Int() < 0 {
				panic("zero int")
			}
		}
	}
}

func BenchmarkVariantIntSliceTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.Type() == uvariant.VariantTypeInt {
				if v.Int() < 0 {
					panic("zero int")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

func BenchmarkVariantStringSliceTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantStringSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			if v.Type() == uvariant.VariantTypeString {
				if v.String() == "" {
					panic("empty string")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}