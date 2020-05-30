package ptrstruct

import "github.com/tigrannajaryan/govariant/uvariant"

type Variant struct {
	typ   uvariant.VariantType
	bytes []byte
	str   string
	i     int
	f     float64
}

func (v *Variant) Type() uvariant.VariantType {
	return v.typ
}

func IntVariant(v int) *Variant {
	return &Variant{typ: uvariant.VariantTypeInt, i: v}
}

func StringVariant(v string) *Variant {
	return &Variant{typ: uvariant.VariantTypeString, str: v}
}

func BytesVariant(v []byte) *Variant {
	return &Variant{typ: uvariant.VariantTypeBytes, bytes: v}
}

func Float64Variant(v float64) *Variant {
	return &Variant{typ: uvariant.VariantTypeFloat64, f: v}
}

func (v *Variant) Int() int {
	return v.i
}

func (v *Variant) Float64() float64 {
	return v.f
}

func (v *Variant) String() (s string) {
	return v.str
}

func (v *Variant) Bytes() (b []byte) {
	return v.bytes
}
