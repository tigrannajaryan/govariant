package main

type VariantType int
const (
	VariantTypeEmpty = iota
	VariantTypeInt
	VariantTypeFloat64
	VariantTypeString
	VariantTypeBytes
)

func EmptyVariant() Variant {
	return Variant{}
}
