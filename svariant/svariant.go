package svariant

import "github.com/tigrannajaryan/govariant/uvariant"

type SVariant struct {
	typ   uvariant.VariantType
	bytes []byte
	str   string
	i     int
	f     float64
}

func (v *SVariant) Type() uvariant.VariantType {
	return v.typ
}

func IntSVariant(v int) SVariant {
	return SVariant{typ: uvariant.VariantTypeInt, i: v}
}

func StringSVariant(v string) SVariant {
	return SVariant{typ: uvariant.VariantTypeString, str: v}
}

func BytesSVariant(v []byte) SVariant {
	return SVariant{typ: uvariant.VariantTypeBytes, bytes: v}
}

func Float64SVariant(v float64) SVariant {
	return SVariant{typ: uvariant.VariantTypeFloat64, f: v}
}

func (v *SVariant) Int() int {
	return v.i
}

func (v *SVariant) Float64() float64 {
	return v.f
}

func (v *SVariant) String() (s string) {
	return v.str
}

func (v *SVariant) Bytes() (b []byte) {
	return v.bytes
}
