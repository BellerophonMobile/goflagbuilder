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
	expectations map[string]*expectation

	errors []string
}

func newtestflags() *testflags {
	fmt.Println("New test")
	x := &testflags{}
	x.expectations = make(map[string]*expectation)
	x.errors = make([]string, 0)
	return x
}

func (x *testflags) expect(name string) {
	x.expectations[name] = &expectation{}
}

func (x *testflags) check(t *testing.T) {

	var res string

	for k, v := range x.expectations {
		if !v.found {
			res += "Expectation " + k + " was not found\n"
		}
	}

	for _, s := range x.errors {
		res += s
	}

	if res != "" {
		t.Error("Failed expectations:\n" + res)
	}

}

func (x *testflags) Var(value flag.Value, name string, usage string) {
	fmt.Printf("  Add flag %s %v --- \"%v\"\n", name, value, usage)

	exp, ok := x.expectations[name]
	if !ok {
		x.errors = append(x.errors, "Unexpected variable '"+name+"'\n")
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

	var flags *testflags

	flags = newtestflags()
	err := Into(flags, nil)
	if err == nil {
		t.Error("No error thrown on nil")
	} else if err.Error() != "Cannot build flags from nil" {
		t.Error("Incorrect error thrown on nil: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	err = Into(flags, "Banana")
	if err == nil {
		t.Error("No error thrown on string")
	} else if err.Error() != "Cannot build flags from type string for prefix ''" {
		t.Error("Incorrect error thrown on string: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	err = Into(flags, 7)
	if err == nil {
		t.Error("No error thrown on int")
	} else if err.Error() != "Cannot build flags from type int for prefix ''" {
		t.Error("Incorrect error thrown on int: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	err = Into(flags, 7.0)
	if err == nil {
		t.Error("No error thrown on float")
	} else if err.Error() != "Cannot build flags from type float64 for prefix ''" {
		t.Error("Incorrect error thrown on float: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	err = Into(flags, mystruct{"Banana", 7})
	if err == nil {
		t.Error("No error thrown on struct")
	} else if err.Error() != "Value of type string at FieldA cannot be set" {
		t.Error("Incorrect error thrown on struct: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	err = Into(flags, map[string]interface{}{"MyStruct": mystruct{}})
	if err == nil {
		t.Error("No error thrown on map to value struct")
	} else if err.Error() != "Value of type string at MyStruct.FieldA cannot be set" {
		t.Error("Incorrectly thrown error on map to value struct: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	flags.expect("Name")
	flags.expect("Index")
	err = Into(flags, &mystruct3{})
	if err == nil {
		t.Error("No error thrown on nil nested struct")
	} else if err.Error() != "Cannot build flags from nil pointer for prefix 'Location'" {
		t.Error("Incorrectly thrown error on nil nested struct: " + err.Error())
	}
	flags.check(t)

}

func Test_From_Map(t *testing.T) {

	var flags *testflags

	flags = newtestflags()
	err := Into(flags, make(map[string]int))
	if err != nil {
		t.Error("Incorrectly thrown error on empty map: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	flags.expect("Banana")
	err = Into(flags, map[string]int{"Banana": 7})
	if err != nil {
		t.Error("Incorrectly thrown error on string->int map: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	flags.expect("MyStruct.FieldA")
	flags.expect("MyStruct.FieldB")
	err = Into(flags, map[string]interface{}{"MyStruct": &mystruct{}})
	if err != nil {
		t.Error("Incorrectly thrown error on string->int map: " + err.Error())
	}
	flags.check(t)

}

func Test_From_Struct(t *testing.T) {

	var flags *testflags

	flags = newtestflags()
	flags.expect("FieldA")
	flags.expect("FieldB")
	err := Into(flags, &mystruct{})
	if err != nil {
		t.Error("Incorrectly thrown error on struct: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	flags.expect("Name")
	flags.expect("Index")
	flags.expect("Location.Grid")
	flags.expect("Location.Fraction")
	err = Into(flags, &mystruct2{})
	if err != nil {
		t.Error("Incorrectly thrown error on struct: " + err.Error())
	}
	flags.check(t)

	flags = newtestflags()
	flags.expect("Name")
	flags.expect("Index")
	flags.expect("Location.Grid")
	flags.expect("Location.Fraction")
	err = Into(flags, &mystruct3{Location: &myotherstruct{}})
	if err != nil {
		t.Error("Incorrectly thrown error on struct: " + err.Error())
	}
	flags.check(t)

}
