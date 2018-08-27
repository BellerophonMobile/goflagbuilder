package goflagbuilder

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mystruct struct {
	FieldA string `help:"Field A"`
	FieldB int
}

type myotherstruct struct {
	Grid     uint64
	Fraction float64

	Attrs map[string]string
}

type mystruct2 struct {
	Name  string
	Index int

	DoStuff bool

	Location myotherstruct
}

type mystruct3 struct {
	Name  string
	Index int

	Location *myotherstruct
}

type mystruct4 struct {
	Temp  int64
	Check uint
	Files []string
}

type mystruct5 struct {
	Getter testGetter
}

type testGetter struct {
	value string
}

func (t *testGetter) Set(s string) error {
	t.value = s
	return nil
}

func (t *testGetter) Get() interface{} { return t.value }
func (t *testGetter) String() string   { return t.value }

func TestInto_Invalid(t *testing.T) {
	suite := []struct {
		name string
		conf interface{}
		err  string
	}{
		{"Nil", nil, "Cannot build flags from nil"},
		{"String", "Banana", "Cannot build flags from type string for prefix ''"},
		{"Int", 7, "Cannot build flags from type int for prefix ''"},
		{"Float", 7.0, "Cannot build flags from type float64 for prefix ''"},
		{
			name: "Struct",
			conf: mystruct{"Banana", 7},
			err:  "Value of type string at FieldA cannot be set",
		},
		{
			name: "Map to Struct",
			conf: map[string]interface{}{"MyStruct": mystruct{"Banana", 7}},
			err:  "Value of type string at MyStruct.FieldA cannot be set",
		},
		{
			name: "Map without string keys",
			conf: map[int]interface{}{10: mystruct{"Banana", 7}},
			err:  "Map key must be string, got int for prefix ''",
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet(item.name, flag.ContinueOnError)
			err := Into(flagSet, item.conf)

			if err == nil {
				t.Error("Did not get expected error:", item.err)
			} else {
				actualError := err.Error()
				if item.err != actualError {
					t.Errorf("Expected error did not match actual.\nExpected: %s\n  Actual: %s", item.err, actualError)
				}
			}
		})
	}
}

type expectedVariable struct {
	value interface{}
	usage string
}

func TestInto(t *testing.T) {
	suite := []struct {
		name string
		conf interface{}
		args []string
		vars map[string]expectedVariable
	}{
		{
			name: "Empty map",
			conf: map[string]int{},
			args: []string{},
			vars: map[string]expectedVariable{},
		},
		{
			name: "Map to Int",
			conf: map[string]int{"Banana": 7},
			args: []string{"-Banana", "10"},
			vars: map[string]expectedVariable{
				"Banana": {value: 10},
			},
		},
		{
			name: "Map to Struct Ptr",
			conf: map[string]interface{}{"MyStruct": &mystruct{}},
			args: []string{"-MyStruct.FieldA", "asdf", "-MyStruct.FieldB", "12"},
			vars: map[string]expectedVariable{
				"MyStruct.FieldA": {value: "asdf", usage: "Field A"},
				"MyStruct.FieldB": {value: 12},
			},
		},
		{
			name: "Struct Ptr",
			conf: &mystruct{},
			args: []string{"-FieldA", "foo", "-FieldB", "21"},
			vars: map[string]expectedVariable{
				"FieldA": {value: "foo", usage: "Field A"},
				"FieldB": {value: 21},
			},
		},
		{
			name: "Nested Struct",
			conf: &mystruct2{},
			args: []string{"-Name", "foo", "-Index", "10", "-DoStuff", "-Location.Grid", "2048", "-Location.Fraction", "3.14"},
			vars: map[string]expectedVariable{
				"Name":              {value: "foo"},
				"Index":             {value: 10},
				"DoStuff":           {value: true},
				"Location.Grid":     {value: uint64(2048)},
				"Location.Fraction": {value: 3.14},
			},
		},
		{
			name: "Nested Struct Ptr",
			conf: &mystruct3{Location: &myotherstruct{}},
			args: []string{"-Name", "bar", "-Index", "20", "-Location.Grid", "1000", "-Location.Fraction", "2.71"},
			vars: map[string]expectedVariable{
				"Name":              {value: "bar"},
				"Index":             {value: 20},
				"Location.Grid":     {value: uint64(1000)},
				"Location.Fraction": {value: 2.71},
			},
		},
		{
			name: "Struct with Nested Map",
			conf: &myotherstruct{Attrs: map[string]string{"Foo": "Bar"}},
			args: []string{"-Grid", "12", "-Fraction", "1.23", "-Attrs.Foo", "AAA"},
			vars: map[string]expectedVariable{
				"Grid":      {value: uint64(12)},
				"Fraction":  {value: 1.23},
				"Attrs.Foo": {value: "AAA"},
			},
		},
		{
			name: "Struct with nil Pointer",
			conf: &mystruct3{},
			args: []string{"-Index", "10", "-Location.Fraction", "3.14"},
			vars: map[string]expectedVariable{
				"Name":              {value: ""},
				"Index":             {value: 10},
				"Location.Grid":     {value: uint64(0)},
				"Location.Fraction": {value: 3.14},
			},
		},
		{
			name: "Struct with Slice",
			conf: &mystruct4{},
			args: []string{"-Temp", "-10", "-Check", "5", "-Files", "foo.log", "-Files", "bar.txt"},
			vars: map[string]expectedVariable{
				"Temp":  {value: int64(-10)},
				"Check": {value: uint(5)},
				"Files": {value: []string{"foo.log", "bar.txt"}},
			},
		},
		{
			name: "Struct with Getter",
			conf: &mystruct5{},
			args: []string{"-Getter", "Foo"},
			vars: map[string]expectedVariable{
				"Getter": {value: "Foo"},
			},
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet(item.name, flag.ContinueOnError)

			if err := Into(flagSet, item.conf); err != nil {
				t.Error("Unexpected error:", err)
				return
			}

			if err := flagSet.Parse(item.args); err != nil {
				t.Error("Error parsing args:", err)
				return
			}

			flagSet.VisitAll(func(f *flag.Flag) {
				if _, ok := item.vars[f.Name]; ok {
					return
				}
				t.Error("Unexpected variable:", f.Name)
			})

			for name, expected := range item.vars {
				f := flagSet.Lookup(name)
				if f == nil {
					t.Error("Expected variable was not found:", name)
					return
				}

				if expected.usage != f.Usage {
					t.Errorf("Usage doesn't match\nExpected: %s\n  Actual: %s", expected.usage, f.Usage)
				}

				getter, ok := f.Value.(flag.Getter)
				if !ok {
					t.Fatal("Value not getter?")
					return
				}

				assert.Equal(t, expected.value, getter.Get(), "Values for %s not equal", name)
			}
		})
	}
}
