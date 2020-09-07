// +build amd64

package cvariant

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestVariantFloat64FieldStorage(t *testing.T) {
	v := Variant{}

	// Ensure float64 can correctly fit in bits
	assert.EqualValues(t, unsafe.Sizeof(float64(0.0)), unsafe.Sizeof(v.bits))
}
