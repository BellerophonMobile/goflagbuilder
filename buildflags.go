package goflagbuilder

import (
	"errors"
	"fmt"
	"reflect"
)

var Strict = false

func populateMapFlags(flags FlagSet, parser *Parser, prefix string, mapval reflect.Value) error {

	if prefix != "" {
		prefix += "."
	}

	for _, keyval := range mapval.MapKeys() {

		subprefix := prefix + keyval.String()
		elementval := mapval.MapIndex(keyval)

		switch elementval.Kind() {
		case reflect.Bool:
			fallthrough

		case reflect.Float64:
			fallthrough

		case reflect.Int64:
			fallthrough

		case reflect.Int:
			fallthrough

		case reflect.String:
			fallthrough

		case reflect.Uint64:
			fallthrough

		case reflect.Uint:
			// Potentially get a usage string from a tag?

			set := &mapvalue{subprefix, mapval, keyval, elementval}
			flags.Var(set, subprefix, "")
			parser.add(subprefix, set)

		default:
			err := recurseBuildFlags(flags, parser, subprefix, elementval)
			if err != nil {
				return err
			}

		}

	}

	return nil

}

func populateStructFlags(flags FlagSet, parser *Parser, prefix string, structval reflect.Value) error {

	if prefix != "" {
		prefix += "."
	}

	structtype := structval.Type()
	for i := 0; i < structval.NumField(); i++ {

		field := structtype.Field(i)
		elementval := structval.Field(i)
		subprefix := prefix + field.Name

		switch elementval.Kind() {
		case reflect.Bool:
			fallthrough

		case reflect.Float64:
			fallthrough

		case reflect.Int64:
			fallthrough

		case reflect.Int:
			fallthrough

		case reflect.String:
			fallthrough

		case reflect.Uint64:
			fallthrough

		case reflect.Uint:
			if !elementval.CanSet() {
				if Strict {
					return fmt.Errorf("Value of type %s at %s cannot be set", elementval.Type().Name(), subprefix)
				} else {
					continue
				}
			}

			set := &flagvalue{subprefix, elementval}
			flags.Var(set, subprefix, field.Tag.Get("help"))
			parser.add(subprefix, set)

		default:
			err := recurseBuildFlags(flags, parser, subprefix, elementval)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func recursePtrFlags(flags FlagSet, parser *Parser, prefix string, ptrval reflect.Value) error {

	if ptrval.IsNil() {
		if Strict {
			return fmt.Errorf("Cannot build flags from nil pointer for prefix '%s'", prefix)
		} else {
			return nil
		}
	}

	return recurseBuildFlags(flags, parser, prefix, ptrval.Elem())

}

func recurseBuildFlags(flags FlagSet, parser *Parser, prefix string, elementval reflect.Value) error {

	switch elementval.Kind() {
	case reflect.Map:
		return populateMapFlags(flags, parser, prefix, elementval)

	case reflect.Struct:
		return populateStructFlags(flags, parser, prefix, elementval)

	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		return recursePtrFlags(flags, parser, prefix, elementval)

	default:
		if Strict {
			return fmt.Errorf("Cannot build flags from type %v for prefix '%s'", elementval.Type(), prefix)
		} else {
			return nil
		}
	}

	return nil

}

// Into populates the given flag set with hierarchical fields from the
// given object.  It returns a Parser that may be used to read those
// same flags from a configuration file.
func Into(flags FlagSet, configuration interface{}) (*Parser, error) {

	if configuration == nil {
		return nil, errors.New("Cannot build flags from nil")
	}

	var parser = newparser()

	err := recurseBuildFlags(flags, parser, "", reflect.ValueOf(configuration))
	if err != nil {
		return nil, err
	}

	return parser, nil

}

// From populates the top-level default flags with hierarchical fields
// from the given object.  It simply calls Into() with configuration
// on a facade of the top-level flag package functions, and returns
// the resultant Parser or error.
func From(configuration interface{}) (*Parser, error) {
	return Into(toplevelflags, configuration)
}
