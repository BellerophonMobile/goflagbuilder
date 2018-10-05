package goflagbuilder

import (
	"flag"
	"log"
)

type server struct {
	Domain string
	Port   int
}

func Example_simple() {
	myserver := &server{}

	// Construct the flags
	if err := From(myserver); err != nil {
		log.Fatal("error: " + err.Error())
	}

	// Read from the command line to establish the param
	flag.Parse()
	// Output:
}
