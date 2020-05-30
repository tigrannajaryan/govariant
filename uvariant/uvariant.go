package uvariant

type VariantType int

const (
	VariantTypeEmpty = iota
	VariantTypeInt
	VariantTypeFloat64
	VariantTypeString
	VariantTypeMap
	VariantTypeBytes
)

func EmptyVariant() Variant {
	return Variant{}
}
