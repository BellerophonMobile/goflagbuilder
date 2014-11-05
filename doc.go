/*

Package goflagbuilder constructs command line flags and a file parser
to manipulate a given structure.  It uses reflection to traverse a
potentially hierarchical object of structs and maps and install
handlers in Go's standard flag package.  Constructed parsers scan
given documents line by line and apply any key/value pairs found.

Constructed flags have the form Foo=..., Bar=...,Obj.Field=... and so
on.  Maps and exposed fields of structs with primitive types are
consumed, so in this case Foo and Bar might be map keys or public
struct fields to a primitive.  Nested maps and structs are followed,
producing dot-notation hierarchical keys such as Obj.Field.

Primitive types understood by goflagbuilder include bool, float64,
int64, int, string, uint64, and uint.

Primitive fields in the given object and sub-objects must be settable.
In general this means structs should be passed in as pointers.  Maps
may also be set directly.

*/
package goflagbuilder
