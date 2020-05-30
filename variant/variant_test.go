package variant

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"unsafe"

	"github.com/tigrannajaryan/govariant/testutil"

	"github.com/stretchr/testify/assert"
)

func TestVariant(t *testing.T) {
	fmt.Printf("Variant size=%v bytes\n", unsafe.Sizeof(Variant{}))

	v := NewEmpty()
	assert.EqualValues(t, VTypeEmpty, v.Type())

	b1 := []byte{1, 2, 3}
	v = NewBytes(b1)
	b2 := v.Bytes()
	assert.EqualValues(t, b1, b2)
	assert.EqualValues(t, VTypeBytes, v.Type())

	s1 := "abcdef"
	v = NewString(s1)
	s2 := v.String()
	assert.EqualValues(t, s1, s2)
	assert.EqualValues(t, VTypeString, v.Type())

	i1 := 1234
	v = NewInt(i1)
	i2 := v.Int()
	assert.EqualValues(t, i1, i2)
	assert.EqualValues(t, VTypeInt, v.Type())

	f1 := 1234.567
	v = NewFloat64(f1)
	f2 := v.Float64()
	assert.EqualValues(t, f1, f2)
	assert.EqualValues(t, VTypeFloat64, v.Type())
}

func TestVariantMap(t *testing.T) {
	v := NewMap(0)
	assert.EqualValues(t, VTypeMap, v.Type())
	assert.EqualValues(t, map[string]Variant{}, v.Map())

	v.Map()["k"] = NewInt(123)
	assert.EqualValues(t, map[string]Variant{"k": NewInt(123)}, v.Map())
}

func TestVariantValueList(t *testing.T) {
	v := NewValueList(nil)
	assert.EqualValues(t, VTypeValueList, v.Type())
	assert.EqualValues(t, 0, v.Len())
	assert.EqualValues(t, []Variant(nil), v.ValueList())

	v = NewValueList([]Variant{NewInt(123), NewString("abc")})
	assert.EqualValues(t, 2, v.Len())
	assert.EqualValues(t, []Variant{NewInt(123), NewString("abc")}, v.ValueList())
	assert.EqualValues(t, NewInt(123), v.ValueAt(0))
	assert.EqualValues(t, NewString("abc"), v.ValueAt(1))
}

func TestVariantKeyValueList(t *testing.T) {
	v := NewKeyValueList(0)
	assert.EqualValues(t, VTypeKeyValueList, v.Type())
	assert.EqualValues(t, 0, v.Len())
	assert.EqualValues(t, []KeyValue{}, v.KeyValueList())

	v = NewKeyValueList(2)
	assert.EqualValues(t, 0, v.Len())

	v.Resize(2)
	assert.EqualValues(t, 2, v.Len())
	assert.EqualValues(t, []KeyValue{{}, {}}, v.KeyValueList())

	kv := v.KeyValueAt(0)
	assert.NotNil(t, kv)
	kv.Key = "key1"
	kv.Value = NewString("value1")

	kv = v.KeyValueAt(1)
	assert.NotNil(t, kv)
	kv.Key = "key2"
	kv.Value = NewFloat64(1.23)

	kv = v.KeyValueAt(0)
	assert.NotNil(t, kv)
	assert.EqualValues(t, "key1", kv.Key)
	assert.EqualValues(t, "value1", kv.Value.String())

	kv = v.KeyValueAt(1)
	assert.NotNil(t, kv)
	assert.EqualValues(t, "key2", kv.Key)
	assert.EqualValues(t, 1.23, kv.Value.Float64())

	list := v.KeyValueList()
	assert.EqualValues(t, "key1", list[0].Key)
	assert.EqualValues(t, "value1", list[0].Value.String())
	assert.EqualValues(t, "key2", list[1].Key)
	assert.EqualValues(t, 1.23, list[1].Value.Float64())

	// Update an element.
	list[0] = KeyValue{"newkey", NewInt(123)}

	// Verify that it updated correctly.
	list = v.KeyValueList()
	assert.EqualValues(t, "newkey", list[0].Key)
	assert.EqualValues(t, 123, list[0].Value.Int())
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

	s2 := v.String()

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	bb = nil

	runtime.GC()

	runtime.ReadMemStats(&ms)

	s2 = v.String()
	s2 = s2
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
		vi := v.Int()
		if vi != testutil.IntMagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantFloat64Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantFloat64()
		vf := v.Float64()
		if vf != testutil.Float64MagicVal {
			panic("invalid value")
		}
	}
}

func BenchmarkVariantStringTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantString()
		if v.Type() == VTypeString {
			if v.String() == "" {
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

func BenchmarkVariantIntTypeAndGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := createVariantInt()
		if v.Type() == VTypeInt {
			vi := v.Int()
			if vi != testutil.IntMagicVal {
				panic("invalid value")
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
			if v.Int() < 0 {
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
				if v.Int() < 0 {
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
				if v.String() == "" {
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
			if v.String() == "" {
				panic("empty string")
			}
		}
	}
}

func BenchmarkVariantValueListGetAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		vv := NewValueList(createVariantStringSlice(testutil.VariantListSize))
		for _, v := range vv.ValueList() {
			if v.String() == "" {
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
