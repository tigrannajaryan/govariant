package main

type SVariant struct {
	typ VariantType
	bytes []byte
	str string
	i int
	f float64
}

func (v* SVariant) Type() VariantType {
	return v.typ
}

func IntSVariant(v int) SVariant {
	return SVariant{typ:VariantTypeInt, i: v}
}

func StringSVariant(v string) SVariant {
	return SVariant{typ:VariantTypeString, str: v}
}

func BytesSVariant(v []byte) SVariant {
	return SVariant{typ:VariantTypeBytes, bytes: v}
}

func Float64SVariant(v float64) SVariant {
	return SVariant{typ:VariantTypeFloat64, f: v}
}

func (v* SVariant) Int() int {
	return v.i
}

func (v* SVariant) Float64() float64 {
	return v.f
}

func (v* SVariant) String() (s string) {
	return v.str
}

func (v* SVariant) Bytes() (b []byte) {
	return v.bytes
}
