package goflagbuilder

import (
	"flag"
)

// FlagSet is the interface for handling flags identified by
// GoFlagBuilder.  FlagSet objects from Go's standard flag package
// meet this specification and are the intended primary target, in
// addition to an internal facade in front of the flag package's top
// level function, and the GoFlagBuilder Parser.
type FlagSet interface {
	Var(value flag.Value, name string, usage string)
}
