package main

import "testing"

const VariantSliceSize = 10
const IntMagicVal = 12345678
const Float64MagicVal = 12345678.9
const StrMagicVal = "abc"
var BytesMagicVal = []byte(StrMagicVal)

func createIVariantInt() IVariant {
	for i :=0; i<1; i++ {
		ivi := IVariantInt(IntMagicVal)
		return &ivi
	}
	return nil
}

func createIVariantFloat64() IVariant {
	for i :=0; i<1; i++ {
		ivi := IVariantFloat64(Float64MagicVal)
		return &ivi
	}
	return nil
}

func createIVariantString() IVariant {
	for i :=0; i<1; i++ {
		ivs := IVariantString(StrMagicVal)
		return &ivs
	}
	return nil
}

func createIVariantBytes() IVariant {
	for i :=0; i<1; i++ {
		ivs := IVariantBytes(BytesMagicVal)
		return &ivs
	}
	return nil
}

func createIVariantIntSlice(n int) []IVariant {
	v := make([]IVariant, n)
	for i :=0; i<n; i++ {
		ivi := IVariantInt(i)
		v[i] = &ivi
	}
	return v
}

func createIVariantStringSlice(n int) []IVariant {
	v := make([]IVariant, n)
	for i :=0; i<n; i++ {
		ivs := IVariantString(StrMagicVal)
		v[i] = &ivs
	}
	return v
}

func BenchmarkIVariantIntGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createIVariantInt()
		vi := v.Int()
		if vi!=IntMagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkIVariantFloat64Get(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createIVariantFloat64()
		vi := v.Float64()
		if vi!=Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkIVariantIntTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createIVariantInt()
		switch val:=v.(type) {
		case *IVariantInt:
			vi := val.Int()
			if vi!=IntMagicVal {
				panic("invalid value")
			}
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkIVariantStringTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createIVariantString()
		switch val:=v.(type) {
		case *IVariantString: val.String()
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkIVariantBytesTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		v := createIVariantBytes()
		switch val:=v.(type) {
		case *IVariantBytes: val.Bytes()
		default:
			panic("invalid type")
		}
	}
}

func BenchmarkIVariantSliceIntGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		vv := createIVariantIntSlice(VariantSliceSize)
		for _,v := range vv {
			v.Int()
		}
	}
}

func BenchmarkIVariantSliceIntTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		vv := createIVariantIntSlice(VariantSliceSize)
		for _,v := range vv {
			switch val:=v.(type) {
			case *IVariantInt: val.Int()
			default:
				panic("invalid type")
			}
		}
	}
}

func BenchmarkIVariantSliceStringTypeAndGet(b *testing.B) {
	for i:=0; i<b.N; i++ {
		vv := createIVariantStringSlice(VariantSliceSize)
		for _,v := range vv {
			switch val:=v.(type) {
			case *IVariantString: val.String()
			default:
				panic("invalid type")
			}
		}
	}
}
