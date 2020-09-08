// +build 386

package variant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariantNewStringFromBytesTooLarge(t *testing.T) {
	assert.NotPanics(t, func() {
		b := make([]byte, MaxSliceLen)
		NewStringFromBytes(b)
	})
	assert.Panics(t, func() {
		b := make([]byte, MaxSliceLen+1)
		NewStringFromBytes(b)
	})
}

func TestVariantNewBytesTooLarge(t *testing.T) {
	assert.NotPanics(t, func() {
		b := make([]byte, MaxSliceLen)
		NewBytes(b)
	})
	assert.Panics(t, func() {
		b := make([]byte, MaxSliceLen+1)
		NewBytes(b)
	})
}

func TestVariantNewValueList(t *testing.T) {
	assert.Panics(t, func() {
		b := make([]Variant, MaxSliceLen+1)
		NewValueList(b)
	})
}

func TestVariantNewKeyValueList(t *testing.T) {
	assert.Panics(t, func() {
		b := make([]KeyValue, MaxSliceLen+1)
		NewKeyValueList(b)
	})
}

func TestVariantResizeTooLarge(t *testing.T) {
	assert.Panics(t, func() {
		b := make([]byte, 0, MaxSliceLen+1)
		v := NewBytes(b)
		v.Resize(MaxSliceLen + 1)
	})
}
