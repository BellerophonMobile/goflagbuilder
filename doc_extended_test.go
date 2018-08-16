package goflagbuilder

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/BellerophonMobile/goflagbuilder/conf"
	"github.com/BellerophonMobile/goflagbuilder/env"
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

func Example_extended() {

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

	// Construct the flags
	if err := From(masterconf); err != nil {
		log.Fatal("CONSTRUCTION ERROR:", err)
	}

	// Read from a config file
	reader := strings.NewReader(`
		# Comment
		Foo.Port = 1234
		Bar.Nested.Index=7.9 # SuccesS!
	`)
	if err := conf.Parse(reader, nil); err != nil {
		log.Fatal("Error:", err)
	}

	// Override settings from the environment
	if err := env.Parse(nil); err != nil {
		log.Fatal("Error:", err)
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
