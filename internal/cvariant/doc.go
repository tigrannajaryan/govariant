/*
Package cvariant implements experimental, more compact Variant data type.

Compared to Variant type found in the "variant" package "cvariant"
implementation trades a small amount of speed for 8 byte reduction of
space per Variant on 64 bits systems. The implementation currently targets amd64 GOARCH
only (it can be extended to other architectures).

On 64 bit systems the size of Variant is 16 bytes for primitive types such as int or
float64.
*/
package cvariant
