package goflagbuilder

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

type foocomponent struct {
	Domain string
	Port   int
}

type nestedstruct struct {
	Index float64
}

type barcomponent struct {
	Label  string
	Nested *nestedstruct
}

func Example_Extended() {

	// Create some sample data
	masterconf := make(map[string]interface{})

	masterconf["Foo"] = &foocomponent{
		Domain: "example.com",
		Port:   9999,
	}

	masterconf["Bar"] = &barcomponent{
		Label: "Bar Component",
		Nested: &nestedstruct{
			Index: 79.3,
		},
	}

	// Construct the flags and parser
	parser, err := From(masterconf)
	if err != nil {
		log.Fatal("CONSTRUCTION ERROR: " + err.Error())
	}

	// Read from a config file
	err = parser.Parse(strings.NewReader(`# Comment
                                        Foo.Port = 1234
                                        Bar.Nested.Index=7.9 # SuccesS!`))
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}

	// Override settings from the command line
	flag.Parse()

	// Output our data
	fmt.Println(masterconf["Foo"].(*foocomponent).Port)
	fmt.Println(masterconf["Bar"].(*barcomponent).Nested.Index)

	// Output:
	// 1234
	// 7.9
}
