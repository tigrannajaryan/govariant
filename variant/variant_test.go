package variant

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"

	"github.com/tigrannajaryan/govariant/internal/testutil"
)

func TestVariantFieldAliasing(t *testing.T) {
	v := Variant{}

	// Ensure fields correctly alias corresponding fields of StringHeader

	// Data/ptr field.
	assert.EqualValues(t, unsafe.Sizeof(reflect.StringHeader{}.Data), unsafe.Sizeof(v.ptr))

	// Len field.
	assert.EqualValues(t, unsafe.Sizeof(reflect.StringHeader{}.Len), unsafe.Sizeof(v.lenAndType))

	// Ensure fields correctly alias corresponding fields of SliceHeader

	// Data/ptr field.
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Data), unsafe.Sizeof(v.ptr))

	// Len field.
	assert.EqualValues(t, unsafe.Sizeof(reflect.SliceHeader{}.Len), unsafe.Sizeof(v.lenAndType))

	// Cap field.
	assert.True(t, unsafe.Sizeof(reflect.SliceHeader{}.Cap) <= unsafe.Sizeof(v.capOrVal))

	// Ensure float64 can correctly fit in capOrVal
	assert.EqualValues(t, unsafe.Sizeof(float64(0.0)), unsafe.Sizeof(v.capOrVal))
}

func TestVariant(t *testing.T) {
	fmt.Printf("Variant size=%v bytes\n", unsafe.Sizeof(Variant{}))

	v := NewEmpty()
	assert.EqualValues(t, VTypeEmpty, v.Type())
	assert.EqualValues(t, "", v.String())

	b1 := []byte{1, 2, 0xA}
	v = NewBytes(b1)
	b2 := v.Bytes()
	assert.EqualValues(t, b1, b2)
	assert.EqualValues(t, VTypeBytes, v.Type())
	assert.EqualValues(t, "0x01020A", v.String())

	s1 := "abcdef"
	v = NewString(s1)
	s2 := v.StringVal()
	assert.EqualValues(t, s1, s2)
	assert.EqualValues(t, VTypeString, v.Type())
	assert.EqualValues(t, `"abcdef"`, v.String())

	i1 := 1234
	v = NewInt(i1)
	i2 := v.IntVal()
	assert.EqualValues(t, i1, i2)
	assert.EqualValues(t, VTypeInt, v.Type())
	assert.EqualValues(t, "1234", v.String())

	f1 := 1234.567
	v = NewFloat64(f1)
	f2 := v.Float64Val()
	assert.EqualValues(t, f1, f2)
	assert.EqualValues(t, VTypeFloat64, v.Type())
	assert.EqualValues(t, "1234.567", v.String())
}

func TestVariantValueList(t *testing.T) {
	v := NewValueList(nil)
	assert.EqualValues(t, VTypeValueList, v.Type())
	assert.EqualValues(t, 0, v.Len())
	assert.EqualValues(t, []Variant(nil), v.ValueList())
	assert.EqualValues(t, "[]", v.String())
	assert.Panics(t, func() { v.ValueAt(0) }, "should panic on nil slice")

	v = NewValueList([]Variant{NewInt(123), NewString("abc")})
	assert.EqualValues(t, 2, v.Len())
	assert.EqualValues(t, []Variant{NewInt(123), NewString("abc")}, v.ValueList())
	assert.EqualValues(t, NewInt(123), v.ValueAt(0))
	assert.EqualValues(t, NewString("abc"), v.ValueAt(1))
	assert.EqualValues(t, `[123,"abc"]`, v.String())
	assert.Panics(t, func() { v.ValueAt(-1) }, "should panic on negative index")
	assert.Panics(t, func() { v.ValueAt(2) }, "should panic on out of bounds")
}

func TestVariantKeyValueList(t *testing.T) {
	var nilKvl []KeyValue
	v := NewKeyValueList(nilKvl)
	assert.EqualValues(t, VTypeKeyValueList, v.Type())
	assert.EqualValues(t, 0, v.Len())
	assert.EqualValues(t, nilKvl, v.KeyValueList())
	assert.EqualValues(t, "{}", v.String())
	assert.Panics(t, func() { v.KeyValueAt(0) }, "should panic on nil slice")

	v = NewKeyValueList(make([]KeyValue, 0, 2))
	assert.EqualValues(t, 0, v.Len())
	assert.EqualValues(t, "{}", v.String())
	assert.Panics(t, func() { v.KeyValueAt(1) })

	v.Resize(2)
	assert.EqualValues(t, 2, v.Len())
	assert.EqualValues(t, []KeyValue{{}, {}}, v.KeyValueList())
	assert.EqualValues(t, `{"":,"":}`, v.String())
	assert.Panics(t, func() { v.KeyValueAt(-1) }, "should panic on negative index")
	assert.Panics(t, func() { v.KeyValueAt(3) }, "should panic on out of bounds")

	kv := v.KeyValueAt(0)
	assert.NotNil(t, kv)
	kv.Key = "key1"
	kv.Value = NewString("value1")
	assert.EqualValues(t, `{"key1":"value1","":}`, v.String())

	kv = v.KeyValueAt(1)
	assert.NotNil(t, kv)
	kv.Key = `key2"`
	kv.Value = NewFloat64(1.23)
	assert.EqualValues(t, `{"key1":"value1","key2\"":1.23}`, v.String())

	kv = v.KeyValueAt(0)
	assert.NotNil(t, kv)
	assert.EqualValues(t, "key1", kv.Key)
	assert.EqualValues(t, "value1", kv.Value.StringVal())

	kv = v.KeyValueAt(1)
	assert.NotNil(t, kv)
	assert.EqualValues(t, `key2"`, kv.Key)
	assert.EqualValues(t, 1.23, kv.Value.Float64Val())

	list := v.KeyValueList()
	assert.EqualValues(t, "key1", list[0].Key)
	assert.EqualValues(t, "value1", list[0].Value.StringVal())
	assert.EqualValues(t, `key2"`, list[1].Key)
	assert.EqualValues(t, 1.23, list[1].Value.Float64Val())

	// Update an element.
	list[0] = KeyValue{"newKey", NewInt(123)}

	// Verify that it updated correctly.
	list = v.KeyValueList()
	assert.EqualValues(t, "newKey", list[0].Key)
	assert.EqualValues(t, 123, list[0].Value.IntVal())
	assert.EqualValues(t, `{"newKey":123,"key2\"":1.23}`, v.String())
}

func TestVariantGC(t *testing.T) {

	var bb []*Variant

	var v1 Variant
	s1 := strconv.Itoa(1234)
	v1 = NewString(s1)
	v1 = v1

	for i := 0; i < 10000; i++ {
		s1 := strconv.Itoa(i)
		vi := new(Variant)
		*vi = NewString(s1)
		b := vi
		bb = append(bb, b)
	}

	var v Variant
	s1 = strconv.Itoa(1234)
	v = NewString(s1)

	s2 := v.StringVal()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	bb = nil

	runtime.GC()

	runtime.ReadMemStats(&ms)

	s2 = v.StringVal()
	s2 = s2
}

func TestVariantPanics(t *testing.T) {
	vals := []Variant{
		NewEmpty(),
		NewInt(123),
		NewFloat64(1.23),
	}
	for _, v := range vals {
		t.Run(v.String(), func(t *testing.T) {
			// Check that get value of the type that is not what the function expects panics.
			assert.Panics(t, func() { v.StringVal() })
			assert.Panics(t, func() { v.Bytes() })
			assert.Panics(t, func() { v.ValueList() })
			assert.Panics(t, func() { v.ValueAt(0) })
			assert.Panics(t, func() { v.Resize(1) })
			assert.Panics(t, func() { v.KeyValueList() })
			assert.Panics(t, func() { v.KeyValueAt(0) })
		})
	}
}

func TestResize(t *testing.T) {
	v := NewBytes([]byte("abc"))
	assert.EqualValues(t, 3, v.Len())
	assert.EqualValues(t, []byte("abc"), v.Bytes())

	v.Resize(2)
	assert.EqualValues(t, 2, v.Len())
	assert.EqualValues(t, []byte("ab"), v.Bytes())

	v.Resize(3)
	assert.EqualValues(t, 3, v.Len())
	assert.EqualValues(t, []byte("abc"), v.Bytes())

	assert.Panics(t, func() { v.Resize(-1) })
	assert.Panics(t, func() { v.Resize(4) })
}

func createVariantInt() Variant {
	for i := 0; i < 1; i++ {
		return NewInt(testutil.IntMagicVal)
	}
	return NewInt(testutil.IntMagicVal)
}

func createVariantFloat64() Variant {
	for i := 0; i < 1; i++ {
		return NewFloat64(testutil.Float64MagicVal)
	}
	return NewFloat64(0)
}

func createVariantString() Variant {
	for i := 0; i < 1; i++ {
		return NewString(testutil.StrMagicVal)
	}
	return NewString("def")
}

func createVariantBytes() Variant {
	for i := 0; i < 1; i++ {
		return NewBytes(testutil.BytesMagicVal)
	}
	return NewBytes(testutil.BytesMagicVal)
}

func BenchmarkVariantIntGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		vi := v.IntVal()
		if vi != testutil.IntMagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantFloat64Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantFloat64()
		vf := v.Float64Val()
		if vf != testutil.Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		if v.Type() == VTypeInt {
			vi := v.IntVal()
			if vi != testutil.IntMagicVal {
				panic("invalid value")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkVariantStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantString()
		if v.Type() == VTypeString {
			if v.StringVal() == "" {
				panic("empty string")
			}
		} else {
			panic("invalid type")
		}
	}
}

func BenchmarkVariantBytesTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantBytes()
		if v.Type() == VTypeBytes {
			if v.Bytes() == nil {
				panic("nil bytes")
			}
		} else {
			panic("invalid type")
		}
	}
}

func createVariantIntSlice(n int) []Variant {
	v := make([]Variant, n)
	for i := 0; i < n; i++ {
		v[i] = NewInt(i)
	}
	return v
}

func createVariantStringSlice(n int) []Variant {
	v := make([]Variant, n)
	for i := 0; i < n; i++ {
		v[i] = NewString("abc")
	}
	return v
}

func BenchmarkVariantIntSliceGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantIntSlice(testutil.VariantListSize)
		for _, v := range vv {
			if v.IntVal() < 0 {
				panic("zero int")
			}
		}
	}
}

func BenchmarkVariantIntSliceTypeAndGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantIntSlice(testutil.VariantListSize)
		for _, v := range vv {
			if v.Type() == VTypeInt {
				if v.IntVal() < 0 {
					panic("zero int")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

func BenchmarkVariantStringSliceTypeAndGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantStringSlice(testutil.VariantListSize)
		for _, v := range vv {
			if v.Type() == VTypeString {
				if v.StringVal() == "" {
					panic("empty string")
				}
			} else {
				panic("invalid type")
			}
		}
	}
}

func BenchmarkVariantStringSliceGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createVariantStringSlice(testutil.VariantListSize)
		for _, v := range vv {
			if v.StringVal() == "" {
				panic("empty string")
			}
		}
	}
}

func BenchmarkVariantValueListGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := NewValueList(createVariantStringSlice(testutil.VariantListSize))
		for _, v := range vv.ValueList() {
			if v.StringVal() == "" {
				panic("empty string")
			}
		}
	}
}

func BenchmarkVariantValueListForRangeAll(b *testing.B) {
	vv := NewValueList(createVariantStringSlice(testutil.VariantListSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range vv.ValueList() {
			if v.Len() == 0 {
				panic("empty string")
			}
		}
	}
}

func BenchmarkVariantValueListAtLenIter(b *testing.B) {
	vv := NewValueList(createVariantStringSlice(testutil.VariantListSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < vv.Len(); j++ {
			v := vv.ValueAt(j)
			if v.Len() == 0 {
				panic("empty string")
			}
		}
	}
}

func BenchmarkVariantValueListAt(b *testing.B) {
	vv := NewValueList(createVariantStringSlice(testutil.VariantListSize))
	b.ResetTimer()
	l := vv.Len()
	j := 0
	for i := 0; i < b.N; i++ {
		v := vv.ValueAt(j)
		if v.Len() == 0 {
			panic("empty string")
		}
		j++
		if j >= l {
			j = 0
		}
	}
}

func BenchmarkVariantValueListLen(b *testing.B) {
	vv := NewValueList(createVariantStringSlice(testutil.VariantListSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if vv.Len() < -i {
			panic("empty string")
		}
	}
}

func BenchmarkStringFromBytes(b *testing.B) {
	bytes := []byte{'a', 'b', 'c'}
	for i := 0; i < b.N; i++ {
		v := NewString(string(bytes))
		str := v.StringVal()
		if str != "abc" {
			panic("invalid string")
		}
	}
}

func BenchmarkStringOptimizedFromBytes(b *testing.B) {
	bytes := []byte{'a', 'b', 'c'}
	for i := 0; i < b.N; i++ {
		v := NewStringFromBytes(bytes)
		str := v.StringVal()
		if str != "abc" {
			panic("invalid string")
		}
	}
}

func createStringSlice(n int) []string {
	v := make([]string, n)
	for i := 0; i < n; i++ {
		v[i] = "abc"
	}
	return v
}

func BenchmarkNativeStringSliceGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := createStringSlice(testutil.VariantListSize)
		for _, v := range vv {
			if len(v) == 0 {
				panic("empty string")
			}
		}
	}
}

func BenchmarkNativeStringSliceForRangeAll(b *testing.B) {
	vv := createStringSlice(testutil.VariantListSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range vv {
			if len(v) == 0 {
				panic("empty string")
			}
		}
	}
}

func BenchmarkNativeStringSliceAtLenIter(b *testing.B) {
	vv := createStringSlice(testutil.VariantListSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(vv); j++ {
			if len(vv[j]) == 0 {
				panic("empty string")
			}
		}
	}
}

func BenchmarkNativeStringSliceAt(b *testing.B) {
	vv := createStringSlice(testutil.VariantListSize)
	b.ResetTimer()
	l := len(vv)
	j := 0
	for i := 0; i < b.N; i++ {
		v := vv[j]
		if len(v) == 0 {
			panic("empty string")
		}
		j++
		if j >= l {
			j = 0
		}
	}
}

func BenchmarkNativeStringSliceLen(b *testing.B) {
	vv := createStringSlice(testutil.VariantListSize)
	for i := 0; i < b.N; i++ {
		if len(vv) < -i {
			panic("empty string")
		}
	}
}
