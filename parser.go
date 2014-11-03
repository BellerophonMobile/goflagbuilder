package gobuildflags

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

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
			return errors.New(fmt.Sprintf("Line %d has no key", line))
		}

		key := strings.TrimSpace(str[:index])
		value := strings.TrimSpace(str[index+1:])

		fmt.Printf("%s -> %s\n", key, value)

		flag, ok := x.flags[key]
		if !ok {
			return errors.New(fmt.Sprintf("Unknown key '%s' on line %d", key, line))
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

func (x *Parser) ParseFile(filename string) error {

	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	return x.Parse(in)

}
