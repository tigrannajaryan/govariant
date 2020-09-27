package variant_test

import (
	"fmt"
	"strings"

	"github.com/tigrannajaryan/govariant/variant"
)

// variantToString converts a Variant to a human readable string.
func variantToString(v variant.Variant) string {
	switch v.Type() {
	case variant.TypeEmpty:
		return ""

	case variant.TypeInt:
		return fmt.Sprintf("%v", v.IntVal())

	case variant.TypeFloat64:
		return fmt.Sprintf("%v", v.Float64Val())

	case variant.TypeString:
		return fmt.Sprintf("%q", v.StringVal())

	case variant.TypeBytes:
		return fmt.Sprintf("0x%X", v.Bytes())

	case variant.TypeValueList:
		sb := strings.Builder{}
		sb.WriteString("[")
		for i, e := range v.ValueList() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(variantToString(e))
		}
		sb.WriteString("]")
		return sb.String()

	case variant.TypeKeyValueList:
		sb := strings.Builder{}
		sb.WriteString("{")
		for i, kv := range v.KeyValueList() {
			if i > 0 {
				sb.WriteString(", ")
			}
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

func Example() {
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

	// Output:
	// 123
	// 1.23
	// "Hello, World!"
	// 0xAFCD34
	// [10, "abc"]
	// {"intval": 10, "a string": "abc", "list": [10, "abc"]}
}
