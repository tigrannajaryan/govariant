package interfacev

import (
	"testing"

	"github.com/tigrannajaryan/govariant/testutil"
)

func createVariantInt() IVariant {
	for i := 0; i < 1; i++ {
		ivi := IVariantInt(testutil.IntMagicVal)
		return &ivi
	}
	return nil
}

func createVariantFloat64() IVariant {
	for i := 0; i < 1; i++ {
		ivi := IVariantFloat64(testutil.Float64MagicVal)
		return &ivi
	}
	return nil
}

func createVariantString() IVariant {
	for i := 0; i < 1; i++ {
		ivs := IVariantString(testutil.StrMagicVal)
		return &ivs
	}
	return nil
}

func createVariantBytes() IVariant {
	for i := 0; i < 1; i++ {
		ivs := IVariantBytes(testutil.BytesMagicVal)
		return &ivs
	}
	return nil
}

func createVariantIntSlice(n int) []IVariant {
	v := make([]IVariant, n)
	for i := 0; i < n; i++ {
		ivi := IVariantInt(i)
		v[i] = &ivi
	}
	return v
}

func createVariantStringSlice(n int) []IVariant {
	v := make([]IVariant, n)
	for i := 0; i < n; i++ {
		ivs := IVariantString(testutil.StrMagicVal)
		v[i] = &ivs
	}
	return v
}

func BenchmarkVariantIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		vi := v.Int()
		if vi != testutil.IntMagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantFloat64Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantFloat64()
		vi := v.Float64()
		if vi != testutil.Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		switch val := v.(type) {
		case *IVariantInt:
			vi := val.Int()
			if vi != testutil.IntMagicVal {
				panic("invalid value")
			}
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkVariantStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantString()
		switch val := v.(type) {
		case *IVariantString:
			val.String()
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkVariantBytesTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantBytes()
		switch val := v.(type) {
		case *IVariantBytes:
			val.Bytes()
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkVariantSliceIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			v.Int()
		}
	}
}

func BenchmarkVariantIntSliceTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			switch val := v.(type) {
			case *IVariantInt:
				val.Int()
			default:
				panic("invalid type")
			}
		}
	}
}

func BenchmarkVariantStringSliceTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantStringSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			switch val := v.(type) {
			case *IVariantString:
				val.String()
			default:
				panic("invalid type")
			}
		}
	}
}
