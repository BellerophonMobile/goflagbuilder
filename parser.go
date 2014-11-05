package goflagbuilder

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parser objects are constructed to manipulate specific objects via
// reading key/value pairs from a document source.
type Parser struct {
	flags map[string]flag.Value
}

func newparser() *Parser {
	x := &Parser{}
	x.flags = make(map[string]flag.Value)
	return x
}

func (x *Parser) add(flag string, set flag.Value) {
	x.flags[flag] = set
}

// Parse reads the given Reader line by line and parses key/value
// pairs to manipulate the associated object.
func (x *Parser) Parse(in io.Reader) error {

	var line int
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line++

		str := scanner.Text()

		index := strings.Index(str, "#")
		if index != -1 {
			if index == 0 || str[index-1] != '\\' {
				str = str[:index]
			}
		}

		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}

		index = strings.Index(str, "=")
		if index < 0 {
			return fmt.Errorf("Line %d has no key", line)
		}

		key := strings.TrimSpace(str[:index])
		value := strings.TrimSpace(str[index+1:])

		flag, ok := x.flags[key]
		if !ok {
			return fmt.Errorf("Unknown key '%s' on line %d", key, line)
		}
		if err := flag.Set(value); err != nil {
			return err
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil

}

// ParseFile reads the file indicated by filename line by line and
// parses key/value pairs to manipulate the associated object.  It is
// identical to calling Parse on a File opened from filename.  No
// error is returned if the file does not exist, on the theory that it
// is equivalent to a blank file which would set no flags.
func (x *Parser) ParseFile(filename string) error {

	in, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer in.Close()

	return x.Parse(in)

}
