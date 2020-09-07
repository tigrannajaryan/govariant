/*
Package variant implements Variant data type.

Variant (also known as tagged union: https://en.wikipedia.org/wiki/Tagged_union)
allows to store values of one of the following types:

 - int,
 - float64,
 - string,
 - []byte slice,
 - ordered list of Variant,
 - ordered key/value list of Variant, where key is a string.
 - empty or no value.

Variant implementation is optimized for performance: for minimal CPU and
memory usage. The implementation currently targets amd64 or 386 GOARCH
only (it can be extended to other architectures).

Usage

To use a Variant first create and store a value in it and then check the stored value
type and read it. For example:

 v := variant.NewInt(123)
 if v.Type() == variant.VTypeInt {
   x := v.IntVal() // x is now int value 123.
 }

Variant uses an efficient data structure that is small and fast to operate on.
On 64 bit systems the size of Variant is 24 bytes for primitive types such as int or
float64. For variable-sized types (String, List, etc) Variant stores a pointer to the
data and length counter (similarly to how Go's built-in string and slice types do).

To maximize the performance Variant functions do not return errors. All functions define
clear contracts that describe in which case the calls are valid. In such cases it is
guaranteed that no errors occur. If the caller violates the contract and it results
in erroneous situation the functions will either return an undefined value or will panic.
The documentation for each function will describe which of those two things will happen.

If a list is stored in the Variant it uses panics to mimic the behavior of builtin Go
slice type. For example accessing an element of a VTypeValueList using an index that is
out of bounds will result in a panic.
*/
package variant
