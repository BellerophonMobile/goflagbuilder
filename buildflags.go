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
			return fmt.Errorf("map key must be string, got %s for prefix '%s'", keyval.Type(), prefix)
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

		// This is true if the field is unexported.  Borrowed from JSON encoder.
		if field.PkgPath != "" {
			continue
		}

		// This is borrowed from stdlib's JSON encoder
		// https://golang.org/src/encoding/json/encode.go#L1102
		isUnexported := field.PkgPath != ""
		if field.Anonymous {
			t := field.Type
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if isUnexported && t.Kind() != reflect.Struct {
				// Ignore embedded fields of unexported non-struct types.
				continue
			}
			// Do not ignore embedded fields of unexported struct types
			// since they may have exported fields.
		} else if isUnexported {
			// Ignore unexported non-embedded fields.
			continue
		}

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
				return fmt.Errorf("value of type %s at %s cannot be set", field.Type.String(), subprefix)
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
			return fmt.Errorf("cannot build flags from nil pointer for prefix '%s'", prefix)
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
		return fmt.Errorf("cannot build flags from type %v for prefix '%s'", elementval.Type(), prefix)
	}
}

// Into populates the given flag set with hierarchical fields from the
// given object.  It returns a Parser that may be used to read those
// same flags from a configuration file.
func Into(flags FlagSet, configuration interface{}) error {
	if configuration == nil {
		return errors.New("cannot build flags from nil")
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
