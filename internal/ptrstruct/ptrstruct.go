package ptrstruct

import "github.com/tigrannajaryan/govariant/variant"

type Variant struct {
	typ   variant.Type
	bytes []byte
	str   string
	i     int
	f     float64
}

func (v *Variant) Type() variant.Type {
	return v.typ
}

func IntVariant(v int) *Variant {
	return &Variant{typ: variant.TypeInt, i: v}
}

func StringVariant(v string) *Variant {
	return &Variant{typ: variant.TypeString, str: v}
}

func BytesVariant(v []byte) *Variant {
	return &Variant{typ: variant.TypeBytes, bytes: v}
}

func Float64Variant(v float64) *Variant {
	return &Variant{typ: variant.TypeFloat64, f: v}
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
