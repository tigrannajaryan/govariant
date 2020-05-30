package ivariant

type IVariant interface {
	Int() int
	String() string
	Float64() float64
	Bytes() []byte
}

type IVariantInt int

func (i IVariantInt) Int() int {
	return int(i)
}

func (i IVariantInt) Float64() float64 {
	panic("invalid type")
}

func (i IVariantInt) String() string {
	panic("invalid type")
}

func (i IVariantInt) Bytes() []byte {
	panic("invalid type")
}

type IVariantFloat64 float64

func (i IVariantFloat64) Float64() float64 {
	return float64(i)
}

func (i IVariantFloat64) Int() int {
	panic("invalid type")
}

func (i IVariantFloat64) String() string {
	panic("invalid type")
}

func (i IVariantFloat64) Bytes() []byte {
	panic("invalid type")
}

type IVariantString string

func (i IVariantString) Int() int {
	panic("invalid type")
}

func (i IVariantString) Float64() float64 {
	panic("invalid type")
}

func (i IVariantString) String() string {
	return string(i)
}

func (i IVariantString) Bytes() []byte {
	panic("invalid type")
}

type IVariantBytes []byte

func (i IVariantBytes) Int() int {
	panic("invalid type")
}

func (i IVariantBytes) Float64() float64 {
	panic("invalid type")
}

func (i IVariantBytes) Bytes() []byte {
	return []byte(i)
}

func (i IVariantBytes) String() string {
	panic("invalid type")
}
