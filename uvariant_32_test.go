// +build 386

package main

import (
	"math/rand"
	"testing"
	"unsafe"
)

func TestFloat64Store(t*testing.T) {
	//v := Variant{}

	// Ensure capOrVal and last32bit form a continuous 64 bit area where float64 can be stored.
	//assert.EqualValues(t, 4, unsafe.Sizeof(v.capOrVal))
	//assert.EqualValues(t, 4, unsafe.Sizeof(v.last32bit))
	//assert.EqualValues(t, unsafe.Offsetof(v.capOrVal)+unsafe.Sizeof(v.capOrVal), unsafe.Offsetof(v.last32bit))
}

func createUVariantF64Float64() VariantF64 {
	for i :=0; i<1; i++ {
		return Float64VariantF64(Float64MagicVal)
	}
	return VariantF64{}
}

//func BenchmarkUVariantF64Float64Get(b *testing.B) {
//	for i:=0; i<b.N; i++ {
//		v := createUVariantF64Float64()
//		vf := v.Float64()
//		if vf!=vf {
//			panic("invalid value")
//		}
//	}
//}


type VS struct {
	ptr unsafe.Pointer
	lenOrType int
	//capOrVal int
	//last32bit int // used for second half of float64.
	capOrVal uint64
}

func fff() Variant {
	return Float64Variant(1.5)
	for i:=0; i<10; i++ {
		return Float64Variant(1.5)
	}
	return Variant{}
}

var f3val float64

func f3() float64 {
	return f3val
}

func BenchmarkFloat64Bits(b *testing.B) {
	f3val = rand.Float64()
	for i:=0; i<b.N; i++ {
		v:=fff()
		f1 := v.Float64()
		//f := f3()
		if f1!=f1 {
			panic("bad bits")
		}
	}
}