package goflagbuilder

import (
	"flag"
	"fmt"
	"strings"
	"testing"
)

type expectation struct {
	found bool
}

type testflags struct {
	t *testing.T

	label string
	id    int

	data interface{}

	expectedVars  map[string]*expectation
	expectedError string

	faults []string

	parser              *Parser
	expectedParserError string
}

var testcount int

func newtest(t *testing.T, label string, data interface{}) *testflags {

	testcount++

	x := &testflags{}

	x.t = t

	x.label = label
	x.id = testcount

	x.data = data

	x.expectedVars = make(map[string]*expectation)
	x.faults = make([]string, 0)

	fmt.Printf("Test %d: %s\n", x.id, x.label)

	return x
}

func (x *testflags) Var(value flag.Value, name string, usage string) {
	fmt.Printf("  Add flag %s %v --- \"%v\"\n", name, value, usage)

	exp, ok := x.expectedVars[name]
	if !ok {
		x.faults = append(x.faults, "Unexpected variable '"+name+"'")
	} else {
		exp.found = true
	}

}

func (x *testflags) variable(name string) {
	x.expectedVars[name] = &expectation{}
}

func (x *testflags) error(error string) {
	x.expectedError = error
}

func (x *testflags) check(err error) {

	var res string

	if err == nil {
		if x.expectedError != "" {
			res += "    Did not get expected error:\n  " + x.expectedError + "\n"
		}
	} else {
		str := err.Error()

		if x.expectedError == "" {
			res += "    Did not expect error:\n  " + str + "\n"
		} else if x.expectedError != str {
			res += "    Expected error did not match received:\n  " + x.expectedError + "\n  " + str + "\n"
		}

	}

	for k, v := range x.expectedVars {
		if !v.found {
			res += "    Expected variable " + k + " was not found\n"
		}
	}

	for _, s := range x.faults {
		res += "    " + s + "\n"
	}

	if res != "" {
		x.t.Error(fmt.Sprintf("Failed test %d: %s\n  Data: %v\n  Faults:\n%s", x.id, x.label, x.data, res))
	}

}

func (x *testflags) run() {
	x.execute(true)
}

func (x *testflags) execute(strict bool) {

	Strict = strict
	var err error
	x.parser, err = Into(x, x.data)
	x.check(err)

}

func (x *testflags) parseerror(err string) {
	x.expectedParserError = err
}

func (x *testflags) parse(text string) {

	x.run()

	reader := strings.NewReader(text)
	err := x.parser.Parse(reader)

	if err == nil {
		if x.expectedParserError != "" {
			x.t.Error("Did not get expected error:\n  " + x.expectedParserError)
		}
	} else {
		if x.expectedParserError == "" {
			x.t.Error("Did not expect error:\n  " + err.Error())
		} else if err.Error() != x.expectedParserError {
			x.t.Error("Expected error did not match received:\n  " + x.expectedParserError + "\n  " + err.Error())
		}
	}

}
