package goflagbuilder

import (
	"flag"
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

	// Output:
}
