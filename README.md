gobuildflags
============

Build command line flags to set fields in a configuration structure
using Go's built-in flag package.

[![Build Status](https://travis-ci.org/BellerophonMobile/gobuildflags.svg)](https://travis-ci.org/BellerophonMobile/gobuildflags) [![GoDoc](https://godoc.org/github.com/BellerophonMobile/gobuildflags?status.svg)](https://godoc.org/github.com/BellerophonMobile/gobuildflags)

## Example

A very simple example:

```go
package main

import (
	"flag"
  "gobuildflags"
	"log"
)

type server struct {
	Domain string
	Port   int
}

func Example_Simple() {

	myserver := &server{}

	// Construct the flags
	_, err := From(myserver)
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// Read from the command line to establish the param
	flag.Parse()

}
```

A more elaborate example including nested structures and using the
parser is available
[here](https://github.com/BellerophonMobile/gobuildflags/blob/master/doc_extended_test.go).
There are also a series of tests in the package outlining exactly what
input structures are valid.


## Major Release Changelog

 * **2014/11/03: Release 1.0!** Though not mature at all, we consider
   gobuildflags to be usable.


## License

gobuildflags is provided under the open source
[MIT license](http://opensource.org/licenses/MIT):

> The MIT License (MIT)
>
> Copyright (c) 2014 [Bellerophon Mobile](http://bellerophonmobile.com/)
> 
>
> Permission is hereby granted, free of charge, to any person
> obtaining a copy of this software and associated documentation files
> (the "Software"), to deal in the Software without restriction,
> including without limitation the rights to use, copy, modify, merge,
> publish, distribute, sublicense, and/or sell copies of the Software,
> and to permit persons to whom the Software is furnished to do so,
> subject to the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
> BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
> ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
> CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
