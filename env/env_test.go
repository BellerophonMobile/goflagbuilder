package env

import (
	"errors"
	"flag"
	"os"
	"testing"
)

func makeFlagSet(name string, flags map[string]interface{}) *flag.FlagSet {
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)

	for k, v := range flags {
		switch v := v.(type) {
		case int:
			flagSet.Int(k, v, "")

		case string:
			flagSet.String(k, v, "")

		default:
			panic("unexpected type in test")
		}
	}

	return flagSet
}

func TestParseValid(t *testing.T) {
	suite := []struct {
		name  string
		env   map[string]string
		start map[string]interface{}
		end   map[string]interface{}
	}{
		{
			name: "parse empty env",
			env:  map[string]string{},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
		},
		{
			name: "parse single value",
			env: map[string]string{
				"PARSE_SINGLE_VALUE_FIELDB": "10",
			},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 10,
			},
		},
		{
			name: "parse multiple values",
			env: map[string]string{
				"PARSE_MULTIPLE_VALUES_FIELDA": "Sushi",
				"PARSE_MULTIPLE_VALUES_FIELDB": "20",
			},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Sushi",
				"FieldB": 20,
			},
		},
		{
			name: "parse extraneous values",
			env: map[string]string{
				"PARSE_EXTRANEOUS_VALUES_FIELDA": "Sushi",
				"PARSE_EXTRANEOUS_VALUES_FIELDB": "20",
				"HELLO":                          "3.14",
				"FOO_BAR":                        "AAA",
			},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Sushi",
				"FieldB": 20,
			},
		},
		{
			name: "./parse/file/names/test",
			env: map[string]string{
				"TEST_FIELDA": "Sushi",
				"TEST_FIELDB": "20",
				"HELLO":       "3.14",
				"FOO_BAR":     "AAA",
			},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Sushi",
				"FieldB": 20,
			},
		},
		{
			name: "/another/file/names/test",
			env: map[string]string{
				"TEST_FIELDA": "Sushi",
				"TEST_FIELDB": "20",
				"HELLO":       "3.14",
				"FOO_BAR":     "AAA",
			},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Sushi",
				"FieldB": 20,
			},
		},
		{
			name: "/file/names/test/",
			env: map[string]string{
				"TEST_FIELDA": "Sushi",
				"TEST_FIELDB": "20",
				"HELLO":       "3.14",
				"FOO_BAR":     "AAA",
			},
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "Sushi",
				"FieldB": 20,
			},
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range item.env {
				os.Setenv(k, v)
			}

			flagSet := makeFlagSet(item.name, item.start)
			// In order to test global FlagSet
			flag.CommandLine = flagSet

			if err := Parse(nil); err != nil {
				t.Error("unexpected error:", err)
			}

			for k, v := range item.end {
				f := flagSet.Lookup(k)
				getter := f.Value.(flag.Getter)
				if v != getter.Get() {
					t.Errorf("values not equal\nExpected: %v\n  Actual: %v", v, getter.Get())
				}
			}
		})
	}
}

type badvar struct{}

func (badvar) Set(string) error { return errors.New("test") }
func (badvar) String() string   { return "" }

func TestParseBadValue(t *testing.T) {
	flagSet := flag.NewFlagSet("TestParseBadValue", flag.ContinueOnError)
	flagSet.Var(badvar{}, "Bar", "")
	flagSet.Var(badvar{}, "Baz", "")
	flagSet.Var(badvar{}, "Foo", "")

	os.Clearenv()
	os.Setenv("TESTPARSEBADVALUE_BAZ", "BAR")

	err := Parse(flagSet)

	if err == nil {
		t.Error("expected set error")
	} else {
		if err.Error() != "test" {
			t.Error("unexpected error:", err.Error())
		}
	}
}
