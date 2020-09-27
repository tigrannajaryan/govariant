// +build 386

package variant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariantNewStringFromBytesTooLarge(t *testing.T) {
	assert.NotPanics(t, func() {
		b := make([]byte, maxSliceLen)
		NewStringFromBytes(b)
	})
	assert.Panics(t, func() {
		b := make([]byte, maxSliceLen+1)
		NewStringFromBytes(b)
	})
}

func TestVariantNewBytesTooLarge(t *testing.T) {
	assert.NotPanics(t, func() {
		b := make([]byte, maxSliceLen)
		NewBytes(b)
	})
	assert.Panics(t, func() {
		b := make([]byte, maxSliceLen+1)
		NewBytes(b)
	})
}

func TestVariantNewValueList(t *testing.T) {
	assert.Panics(t, func() {
		b := make([]Variant, maxSliceLen+1)
		NewValueList(b)
	})
}

func TestVariantNewKeyValueList(t *testing.T) {
	assert.Panics(t, func() {
		b := make([]KeyValue, maxSliceLen+1)
		NewKeyValueList(b)
	})
}

func TestVariantResizeTooLarge(t *testing.T) {
	assert.Panics(t, func() {
		b := make([]byte, 0, maxSliceLen+1)
		v := NewBytes(b)
		v.Resize(maxSliceLen + 1)
	})
}
