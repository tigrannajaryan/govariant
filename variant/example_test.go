package variant_test

import (
	"fmt"

	"github.com/tigrannajaryan/govariant/variant"
)

func ExampleNewInt() {
	v := variant.NewInt(123)
	if v.Type() == variant.VTypeInt {
		fmt.Print(v.IntVal())
	}

	// Output: 123
}

func ExampleVariant_String() {
	v := variant.NewBytes([]byte{1, 2, 0xA})
	fmt.Println(v.String())

	v = variant.NewString("abcdef")
	fmt.Println(v.String())

	v = variant.NewInt(1234)
	fmt.Println(v.String())

	v = variant.NewFloat64(0.12)
	fmt.Println(v.String())

	v = variant.NewValueList(
		[]variant.Variant{
			variant.NewInt(10),
			variant.NewString("abc"),
		},
	)
	fmt.Println(v.String())

	// Output:
	// 0x01020A
	// "abcdef"
	// 1234
	// 0.12
	// [10,"abc"]
}

func ExampleVariant_ValueList() {
	v := variant.NewValueList(
		[]variant.Variant{
			variant.NewInt(10),
			variant.NewString("abc"),
		},
	)

	for i, e := range v.ValueList() {
		fmt.Printf("v[%d]=%s\n", i, e.String())
	}

	// Output:
	// v[0]=10
	// v[1]="abc"
}

func ExampleVariant_KeyValueList() {
	v := variant.NewKeyValueList(
		[]variant.KeyValue{
			{Key: "intval", Value: variant.NewInt(10)},
			{Key: "a string", Value: variant.NewString("abc")},
		},
	)

	for _, kv := range v.KeyValueList() {
		fmt.Printf("%q=%s\n", kv.Key, kv.Value.String())
	}

	// Output:
	// "intval"=10
	// "a string"="abc"
}

func ExampleNewStringFromBytes() {
	bytes := []byte{'a', 'b', 'c'}
	v := variant.NewStringFromBytes(bytes)
	fmt.Println(v.StringVal())

	bytes[2] = 'd'
	fmt.Println(v.StringVal())

	// Output:
	// abc
	// abd
}
