package gobuildflags

import (
	"flag"
	"fmt"
	"testing"
)

type expectation struct {
	found bool
}

type testflags struct {
	label string
	id    int

	data interface{}

	expectedVars  map[string]*expectation
	expectedError string

	faults []string
}

var testcount int = 0

func newtest(label string, data interface{}) *testflags {

	testcount++

	x := &testflags{}

	x.label = label
	x.id = testcount

	x.data = data

	x.expectedVars = make(map[string]*expectation)
	x.faults = make([]string, 0)

	fmt.Printf("Test %d: %s\n", x.id, x.label)

	return x
}

func (x *testflags) variable(name string) {
	x.expectedVars[name] = &expectation{}
}

func (x *testflags) error(error string) {
	x.expectedError = error
}

func (x *testflags) check(err error, t *testing.T) {

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
		t.Error(fmt.Sprintf("Failed test %d: %s\n  Data: %v\n  Faults:\n%s", x.id, x.label, x.data, res))
	}

}

func (x *testflags) run(t *testing.T) {

	_, err := Into(x, x.data)
	x.check(err, t)

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

type mystruct struct {
	FieldA string
	FieldB int
}

type myotherstruct struct {
	Grid     uint64
	Fraction float64
}

type mystruct2 struct {
	Name  string
	Index int

	Location myotherstruct
}

type mystruct3 struct {
	Name  string
	Index int

	Location *myotherstruct
}

func Test_From_Invalid(t *testing.T) {

	var test *testflags

	test = newtest("Nil", nil)
	test.error("Cannot build flags from nil")
	test.run(t)

	test = newtest("String", "Banana")
	test.error("Cannot build flags from type string for prefix ''")
	test.run(t)

	test = newtest("Int", 7)
	test.error("Cannot build flags from type int for prefix ''")
	test.run(t)

	test = newtest("Float", 7.0)
	test.error("Cannot build flags from type float64 for prefix ''")
	test.run(t)

	test = newtest("Struct", mystruct{"Banana", 7})
	test.error("Value of type string at FieldA cannot be set")
	test.run(t)

	test = newtest("Map to Struct",
		map[string]interface{}{"MyStruct": mystruct{}})
	test.error("Value of type string at MyStruct.FieldA cannot be set")
	test.run(t)

	test = newtest("Struct with Nested Nil", &mystruct3{})
	test.variable("Name")
	test.variable("Index")
	test.error("Cannot build flags from nil pointer for prefix 'Location'")
	test.run(t)

}

func Test_From_Map(t *testing.T) {

	var test *testflags

	test = newtest("Empty map", make(map[string]int))
	test.run(t)

	test = newtest("Map to Int", map[string]int{"Banana": 7})
	test.variable("Banana")
	test.run(t)

	test = newtest("Map to Struct Ptr",
		map[string]interface{}{"MyStruct": &mystruct{}})
	test.variable("MyStruct.FieldA")
	test.variable("MyStruct.FieldB")
	test.run(t)

}

func Test_From_Struct(t *testing.T) {

	var test *testflags

	test = newtest("Struct Ptr", &mystruct{})
	test.variable("FieldA")
	test.variable("FieldB")
	test.run(t)

	test = newtest("Nested Struct", &mystruct2{})
	test.variable("Name")
	test.variable("Index")
	test.variable("Location.Grid")
	test.variable("Location.Fraction")
	test.run(t)

	test = newtest("Nested Struct Ptr", &mystruct3{Location: &myotherstruct{}})
	test.variable("Name")
	test.variable("Index")
	test.variable("Location.Grid")
	test.variable("Location.Fraction")
	test.run(t)

}
