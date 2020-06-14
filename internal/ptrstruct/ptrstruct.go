package ptrstruct

import "github.com/tigrannajaryan/govariant/variant"

type Variant struct {
	typ   variant.VType
	bytes []byte
	str   string
	i     int
	f     float64
}

func (v *Variant) Type() variant.VType {
	return v.typ
}

func IntVariant(v int) *Variant {
	return &Variant{typ: variant.VTypeInt, i: v}
}

func StringVariant(v string) *Variant {
	return &Variant{typ: variant.VTypeString, str: v}
}

func BytesVariant(v []byte) *Variant {
	return &Variant{typ: variant.VTypeBytes, bytes: v}
}

func Float64Variant(v float64) *Variant {
	return &Variant{typ: variant.VTypeFloat64, f: v}
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
