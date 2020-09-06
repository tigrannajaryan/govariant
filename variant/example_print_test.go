package variant_test

import (
	"fmt"
	"strings"

	"github.com/tigrannajaryan/govariant/variant"
)

func variantToString(v variant.Variant) string {
	switch v.Type() {
	case variant.VTypeEmpty:
		return ""

	case variant.VTypeInt:
		return fmt.Sprintf("%v", v.IntVal())

	case variant.VTypeFloat64:
		return fmt.Sprintf("%v", v.Float64Val())

	case variant.VTypeString:
		return fmt.Sprintf("%q", v.StringVal())

	case variant.VTypeBytes:
		return fmt.Sprintf("0x%X", v.Bytes())

	case variant.VTypeValueList:
		sb := strings.Builder{}
		sb.WriteString("[")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(variantToString(v.ValueAt(i)))
		}
		sb.WriteString("]")
		return sb.String()

	case variant.VTypeKeyValueList:
		sb := strings.Builder{}
		sb.WriteString("{")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				sb.WriteString(", ")
			}
			kv := v.KeyValueAt(i)
			sb.WriteString(fmt.Sprintf("%q: ", kv.Key))
			sb.WriteString(variantToString(kv.Value))
		}
		sb.WriteString("}")
		return sb.String()
	}
	panic("Unknown variant type")
}

func printVariant(v variant.Variant) {
	fmt.Printf("%s\n", variantToString(v))
}

func ExamplePrintVariant() {
	v := variant.NewInt(123)
	printVariant(v)

	v = variant.NewFloat64(1.23)
	printVariant(v)

	v = variant.NewString("Hello, World!")
	printVariant(v)

	v = variant.NewBytes([]byte{0xAF, 0xCD, 0x34})
	printVariant(v)

	v = variant.NewValueList(
		[]variant.Variant{
			variant.NewInt(10),
			variant.NewString("abc"),
		},
	)
	printVariant(v)

	v = variant.NewKeyValueList(
		[]variant.KeyValue{
			{Key: "intval", Value: variant.NewInt(10)},
			{Key: "a string", Value: variant.NewString("abc")},
			{Key: "list", Value: v},
		},
	)
	printVariant(v)

	// Output: 123
	// 1.23
	// "Hello, World!"
	// 0xAFCD34
	// [10, "abc"]
	// {"intval": 10, "a string": "abc", "list": [10, "abc"]}
}
