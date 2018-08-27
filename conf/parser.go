/*

Package conf provides a simple configuration file format that reads values
backed by a flag.FlagSet.

The configuration file format is single lines of "Key=Value", and comments
marked by "#". Everything after the equals sign is assigned to the key.
Whitespace is trimmed. Comments can be escaped by \#.

For example:

	# A full-line comment
	Foo = 10
	Bar.Baz = hello

*/
package conf

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parse reads the given Reader line by line and parses key/value pairs, setting
// values to matching flags in the given flagset.  If flagSet is nil, then the
// global flag.CommandLine FlagSet is used.
func Parse(in io.Reader, flagSet *flag.FlagSet) error {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}

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
			return fmt.Errorf("line %d has no key", line)
		}

		key := strings.TrimSpace(str[:index])
		value := strings.TrimSpace(str[index+1:])

		flag := flagSet.Lookup(key)
		if flag == nil {
			return fmt.Errorf("unknown key '%s' on line %d", key, line)
		}
		if err := flag.Value.Set(value); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// ParseFile reads the file indicated by filename line by line and parses
// key/value pairs, setting values to matching flags in the given flagset. It is
// identical to calling Parse on a File opened from filename. If flagSet is nil,
// the global flag.CommandLine FlagSet is used.
func ParseFile(filename string, flagSet *flag.FlagSet) error {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}

	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	return Parse(in, flagSet)
}
