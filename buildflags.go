package goflagbuilder

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
)

func populateMapFlags(flags FlagSet, prefix string, mapval reflect.Value) error {
	if prefix != "" {
		prefix += "."
	}

	for _, keyval := range mapval.MapKeys() {
		if keyval.Kind() != reflect.String {
			return fmt.Errorf("Map key must be string, got %s for prefix '%s'", keyval.Type(), prefix)
		}

		subprefix := prefix + keyval.String()
		elementval := mapval.MapIndex(keyval)

		mapKind := mapKind{
			keyval:    keyval,
			valueKind: findKind(elementval),
		}

		if mapKind.valueKind != nil {
			value := value{
				value:  mapval,
				kind:   mapKind,
				isBool: elementval.Kind() == reflect.Bool,
			}
			flags.Var(value, subprefix, "")

		} else if err := recurseBuildFlags(flags, subprefix, elementval); err != nil {
			return err
		}
	}

	return nil
}

func populateStructFlags(flags FlagSet, prefix string, structval reflect.Value) error {
	if prefix != "" {
		prefix += "."
	}

	structtype := structval.Type()
	for i := 0; i < structval.NumField(); i++ {
		field := structtype.Field(i)
		elementval := structval.Field(i)
		subprefix := prefix + field.Name
		help := field.Tag.Get("help")

		value := value{
			value:  elementval,
			kind:   findKind(elementval),
			isBool: elementval.Kind() == reflect.Bool,
		}

		if value.kind != nil {
			if !elementval.CanSet() {
				return fmt.Errorf("Value of type %s at %s cannot be set", field.Type.String(), subprefix)
			}
			flags.Var(value, subprefix, help)

		} else if err := recurseBuildFlags(flags, subprefix, elementval); err != nil {
			return err
		}
	}

	return nil
}

func recursePtrFlags(flags FlagSet, prefix string, ptrval reflect.Value) error {
	if ptrval.IsNil() {
		if ptrval.CanSet() {
			ptrval.Set(reflect.New(ptrval.Type().Elem()))
		} else {
			return fmt.Errorf("Cannot build flags from nil pointer for prefix '%s'", prefix)
		}
	}

	return recurseBuildFlags(flags, prefix, ptrval.Elem())
}

func recurseBuildFlags(flags FlagSet, prefix string, elementval reflect.Value) error {
	switch elementval.Kind() {
	case reflect.Map:
		return populateMapFlags(flags, prefix, elementval)

	case reflect.Struct:
		return populateStructFlags(flags, prefix, elementval)

	case reflect.Interface, reflect.Ptr:
		return recursePtrFlags(flags, prefix, elementval)

	default:
		return fmt.Errorf("Cannot build flags from type %v for prefix '%s'", elementval.Type(), prefix)
	}
}

// Into populates the given flag set with hierarchical fields from the
// given object.  It returns a Parser that may be used to read those
// same flags from a configuration file.
func Into(flags FlagSet, configuration interface{}) error {
	if configuration == nil {
		return errors.New("Cannot build flags from nil")
	}

	err := recurseBuildFlags(flags, "", reflect.ValueOf(configuration))
	if err != nil {
		return err
	}

	return nil
}

// From populates the top-level default flags with hierarchical fields
// from the given object.  It simply calls Into() with configuration
// on a facade of the top-level flag package functions, and returns
// the resultant Parser or error.
func From(configuration interface{}) error {
	return Into(flag.CommandLine, configuration)
}
