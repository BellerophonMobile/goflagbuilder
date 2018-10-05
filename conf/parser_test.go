package conf

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"strings"
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
		input string
		start map[string]interface{}
		end   map[string]interface{}
	}{
		{
			name:  "parse empty doc",
			input: "",
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
			name:  "parse comment",
			input: "# Mushi sushi",
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
			name:  "parse first line",
			input: "FieldA=sushi",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 7,
			},
		},
		{
			name: "parse two lines",
			input: `FieldA=sushi
			        FieldB=9`,
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 9,
			},
		},
		{
			name: "parse two lines preceded by comment",
			input: `# This is a comment
			        FieldA=sushi
			        FieldB=9`,
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 9,
			},
		},
		{
			name: "parse two lines trailed by comment",
			input: `FieldA=sushi
			        FieldB=9
			        # This is a comment`,
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 9,
			},
		},
		{
			name: "parse two lines split by comment",
			input: `FieldA=sushi
			        # This is a comment
			        FieldB=9`,
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 9,
			},
		},
		{
			name: "parse two lines split by comment and blank line",
			input: `FieldA=sushi

			        # This is a comment
			        FieldB=9`,
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 9,
			},
		},
		{
			name: "parse two lines with comments and blank lines",
			input: `FieldA=sushi

			        # This is a comment
			        # This is another comment
			        FieldB=9`,
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 9,
			},
		},
		{
			name:  "parse comment on line",
			input: "FieldA=sushi # Banananana",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
			end: map[string]interface{}{
				"FieldA": "sushi",
				"FieldB": 7,
			},
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			flagSet := makeFlagSet(item.name, item.start)
			reader := strings.NewReader(item.input)

			// In order to test default FlagSet
			flag.CommandLine = flagSet

			if err := Parse(reader, nil); err != nil {
				t.Error("unexpected error:", err)
			}

			for k, v := range item.end {
				f := flagSet.Lookup(k)
				getter := f.Value.(flag.Getter)
				if v != getter.Get() {
					t.Errorf("values not equal\nexpected: %v\n  actual: %v", v, getter.Get())
				}
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	suite := []struct {
		name  string
		input string
		start map[string]interface{}
		err   string
	}{
		{
			name:  "parse no key",
			input: "Sushi",
			err:   "line 1 has no key",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
		},
		{
			name: "parse no key line 2",
			input: `# Line 1
			        sushi`,
			err: "line 2 has no key",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
		},
		{
			name:  "parse unknown key",
			input: "Foo=10",
			err:   "unknown key 'Foo' on line 1",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			flagSet := makeFlagSet(item.name, item.start)
			reader := strings.NewReader(item.input)

			err := Parse(reader, flagSet)
			if err == nil {
				t.Error("expected error:", item.err)
			} else if item.err != err.Error() {
				t.Errorf("errors don't match\nexpected: %s\n  actual: %s", item.err, err.Error())
			}
		})
	}
}

type badreader struct{}

func (badreader) Read(p []byte) (int, error) { return 0, errors.New("test") }

func TestParseBadReader(t *testing.T) {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	err := Parse(badreader{}, flagSet)

	if err == nil {
		t.Error("expected read error")
	} else {
		if err.Error() != "test" {
			t.Error("unexpected error:", err.Error())
		}
	}
}

type badvar struct{}

func (badvar) Set(string) error { return errors.New("test") }
func (badvar) String() string   { return "" }

func TestParseBadValue(t *testing.T) {
	flagSet := flag.NewFlagSet("TestBadValue", flag.ContinueOnError)
	flagSet.Var(badvar{}, "Foo", "")

	reader := strings.NewReader("Foo=10")

	err := Parse(reader, flagSet)

	if err == nil {
		t.Error("expected set error")
	} else {
		if err.Error() != "test" {
			t.Error("unexpected error:", err.Error())
		}
	}
}

func TestParseFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "goflagbuilder-test")
	if err != nil {
		t.Fatal("failed to create temp file:", err)
		return
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("Foo=10"); err != nil {
		t.Fatal("failed to write temp file:", err)
		return
	}
	if err := tmpFile.Sync(); err != nil {
		t.Fatal("failed to sync temp file:", err)
		return
	}

	flagSet := flag.NewFlagSet("TestParseFile", flag.ContinueOnError)
	foo := flagSet.Int("Foo", 5, "")

	if err := ParseFile(tmpFile.Name(), flagSet); err != nil {
		t.Error("failed to parse file:", err)
	}

	if *foo != 10 {
		t.Error("failed to set value")
	}
}

func TestParseFileBad(t *testing.T) {
	flagSet := flag.NewFlagSet("TestParseFileBad", flag.ContinueOnError)

	// In order to test default FlagSet
	flag.CommandLine = flagSet

	err := ParseFile("/goflagbuilder-bad-test", nil)

	if err == nil {
		t.Error("expected parse error")
	} else {
		if err.Error() != "open /goflagbuilder-bad-test: no such file or directory" {
			t.Error("unexpected error:", err.Error())
		}
	}
}
