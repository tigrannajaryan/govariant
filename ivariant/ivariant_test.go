package ivariant

import (
	"testing"

	"github.com/tigrannajaryan/govariant/testutil"
)

func createIVariantInt() IVariant {
	for i := 0; i < 1; i++ {
		ivi := IVariantInt(testutil.IntMagicVal)
		return &ivi
	}
	return nil
}

func createIVariantFloat64() IVariant {
	for i := 0; i < 1; i++ {
		ivi := IVariantFloat64(testutil.Float64MagicVal)
		return &ivi
	}
	return nil
}

func createIVariantString() IVariant {
	for i := 0; i < 1; i++ {
		ivs := IVariantString(testutil.StrMagicVal)
		return &ivs
	}
	return nil
}

func createIVariantBytes() IVariant {
	for i := 0; i < 1; i++ {
		ivs := IVariantBytes(testutil.BytesMagicVal)
		return &ivs
	}
	return nil
}

func createIVariantIntSlice(n int) []IVariant {
	v := make([]IVariant, n)
	for i := 0; i < n; i++ {
		ivi := IVariantInt(i)
		v[i] = &ivi
	}
	return v
}

func createIVariantStringSlice(n int) []IVariant {
	v := make([]IVariant, n)
	for i := 0; i < n; i++ {
		ivs := IVariantString(testutil.StrMagicVal)
		v[i] = &ivs
	}
	return v
}

func BenchmarkInterfaceVariantIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createIVariantInt()
		vi := v.Int()
		if vi != testutil.IntMagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkInterfaceVariantFloat64Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createIVariantFloat64()
		vi := v.Float64()
		if vi != testutil.Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkInterfaceVariantIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createIVariantInt()
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

func BenchmarkInterfaceVariantStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createIVariantString()
		switch val := v.(type) {
		case *IVariantString:
			val.String()
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkInterfaceVariantBytesTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createIVariantBytes()
		switch val := v.(type) {
		case *IVariantBytes:
			val.Bytes()
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkInterfaceVariantSliceIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createIVariantIntSlice(testutil.VariantSliceSize)
		for _, v := range vv {
			v.Int()
		}
	}
}

func BenchmarkInterfaceVariantSliceIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createIVariantIntSlice(testutil.VariantSliceSize)
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

func BenchmarkInterfaceVariantSliceStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createIVariantStringSlice(testutil.VariantSliceSize)
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
