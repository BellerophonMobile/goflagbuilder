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
	/*
		BoolVar(p *bool, name string, value bool, usage string)
		Float64Var(p *float64, name string, value float64, usage string)
		Int64Var(p *int64, name string, value int64, usage string)
		IntVar(p *int, name string, value int, usage string)
		StringVar(p *string, name string, value string, usage string)
		Uint64Var(p *uint64, name string, value uint64, usage string)
		UintVar(p *uint, name string, value uint, usage string)
	*/
}

type toplevel struct{}

var toplevelflags = &toplevel{}

func (f *toplevel) Var(value flag.Value, name string, usage string) {
	flag.Var(value, name, usage)
}

/*
func (x *toplevel) BoolVar(p *bool, name string, value bool, usage string) {
	flag.BoolVar(p, name, value, usage)
}

func (x *toplevel) Float64Var(p *float64, name string, value float64, usage string) {
	flag.Float64Var(p, name, value, usage)
}

func (x *toplevel) Int64Var(p *int64, name string, value int64, usage string) {
	flag.Int64Var(p, name, value, usage)
}

func (x *toplevel) IntVar(p *int, name string, value int, usage string) {
	flag.IntVar(p, name, value, usage)
}

func (x *toplevel) StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)
}

func (x *toplevel) Uint64Var(p *uint64, name string, value uint64, usage string) {
	flag.Uint64Var(p, name, value, usage)
}

func (x *toplevel) UintVar(p *uint, name string, value uint, usage string) {
	flag.UintVar(p, name, value, usage)
}
*/
