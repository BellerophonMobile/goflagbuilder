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
			name:  "Parse Empty Doc",
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
			name:  "Parse Comment",
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
			name:  "Parse First Line",
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
			name: "Parse Two Lines",
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
			name: "Parse Two Lines Preceded by Comment",
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
			name: "Parse Two Lines Trailed by Comment",
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
			name: "Parse Two Lines Split by Comment",
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
			name: "Parse Two Lines Split by Comment and Blank Line",
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
			name: "Parse Two Lines With Comments and Blank Lines",
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
			name:  "Parse Comment on Line",
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
				t.Error("Unexpected error:", err)
			}

			for k, v := range item.end {
				f := flagSet.Lookup(k)
				getter := f.Value.(flag.Getter)
				if v != getter.Get() {
					t.Errorf("Values not equal\nExpected: %v\n  Actual: %v", v, getter.Get())
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
			name:  "Parse No Key",
			input: "Sushi",
			err:   "line 1 has no key",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
		},
		{
			name: "Parse No Key Line 2",
			input: `# Line 1
			        sushi`,
			err: "line 2 has no key",
			start: map[string]interface{}{
				"FieldA": "Banana",
				"FieldB": 7,
			},
		},
		{
			name:  "Parse Unknown Key",
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
				t.Error("Expected error:", item.err)
			} else if item.err != err.Error() {
				t.Errorf("Errors don't match\nExpected: %s\n  Actual: %s", item.err, err.Error())
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
		t.Error("Expected read error")
	} else {
		if err.Error() != "test" {
			t.Error("Unexpected error:", err.Error())
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
		t.Error("Expected set error")
	} else {
		if err.Error() != "test" {
			t.Error("Unexpected error:", err.Error())
		}
	}
}

func TestParseFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "goflagbuilder-test")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
		return
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("Foo=10"); err != nil {
		t.Fatal("Failed to write temp file:", err)
		return
	}
	if err := tmpFile.Sync(); err != nil {
		t.Fatal("Failed to sync temp file:", err)
		return
	}

	flagSet := flag.NewFlagSet("TestParseFile", flag.ContinueOnError)
	foo := flagSet.Int("Foo", 5, "")

	if err := ParseFile(tmpFile.Name(), flagSet); err != nil {
		t.Error("Failed to parse file:", err)
	}

	if *foo != 10 {
		t.Error("Failed to set value")
	}
}

func TestParseFileBad(t *testing.T) {
	flagSet := flag.NewFlagSet("TestParseFileBad", flag.ContinueOnError)

	// In order to test default FlagSet
	flag.CommandLine = flagSet

	err := ParseFile("/goflagbuilder-bad-test", nil)

	if err == nil {
		t.Error("Expected parse error")
	} else {
		if err.Error() != "open /goflagbuilder-bad-test: no such file or directory" {
			t.Error("Unexpected error:", err.Error())
		}
	}
}
