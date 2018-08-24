/*

Package env provides a means to read environment variables into values backed
by a flag.FlagSet.

Environment variables are of the format:

	<FSNAME>_<KEYNAME>

Where FSNAME is the name of the given flag.FlagSet. KEYNAME is the name of
the flag. All spaces and periods are replaced with underscores in both
strings.

*/
package env

import (
	"flag"
	"os"
	"strings"
)

var replacer = strings.NewReplacer(" ", "_", ".", "_")

// Parse reads environment variables and parses into matching flags in the given
// flagset.  If flagSet is nil, the global flag.CommandLine FlagSet is used.
func Parse(flagSet *flag.FlagSet) error {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}

	name := format(flagSet.Name())
	if name != "" {
		name += "_"
	}

	var err error

	flagSet.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}

		keyName := format(f.Name)
		envName := name + keyName

		value, ok := os.LookupEnv(envName)
		if !ok {
			return
		}

		err = f.Value.Set(value)
	})

	return err
}

func format(x string) string {
	return strings.ToUpper(replacer.Replace(strings.TrimSpace(x)))
}
