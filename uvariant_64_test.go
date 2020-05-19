// +build !386

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestFloat64Store(t*testing.T) {
	v := Variant{}

	// Ensure capOrVal is bit area where float64 can be stored.
	assert.EqualValues(t, 8, unsafe.Sizeof(v.capOrVal))
}