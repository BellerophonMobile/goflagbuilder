/**

Package gobuildflags constructs command line flags and a file parser
to manipulate a given structure.  It uses reflection to traverse a
potentially hierarchical object of structs and maps and install
handlers in Go's standard flag package.  Constructed flags have the
form Foo.Bar... where Foo and Bar are keys of the map or exposed
fields of the containing object.  The associated parser simply scans
line by line and applies any such key/value pairs that it finds.

Primitive types understood by gobuildflags include bool, float64,
int64, int, string, uint64, and uint.

Primitive fields in the given object and sub-objects must be settable.
Primitive types included in a map may be set, as well as fields in
structs passed by pointer.

*/
package gobuildflags
