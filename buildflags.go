package gobuildflags

import (
	"errors"
	"fmt"
	"reflect"
)

func populateElement(flags FlagSet, elementval reflect.Value, subprefix string) error {

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

		if !elementval.CanSet() {
			return errors.New(fmt.Sprintf("Value of type %s at %s cannot be set", elementval.Type().Name(), subprefix))
		}

		flags.Var(&flagvalue{subprefix, elementval}, subprefix, "")

	default:
		err := recurseBuildFlags(flags, subprefix, elementval)
		if err != nil {
			return err
		}

	}

	return nil

}

func populateMapFlags(flags FlagSet, prefix string, mapval reflect.Value) error {

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

			flags.Var(&mapvalue{subprefix, mapval, keyval, elementval}, subprefix, "")

		default:
			err := recurseBuildFlags(flags, subprefix, elementval)
			if err != nil {
				return err
			}

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

		err := populateElement(flags, structval.Field(i), prefix+structtype.Field(i).Name)
		if err != nil {
			return err
		}

	}

	return nil
}

func recursePtrFlags(flags FlagSet, prefix string, ptrval reflect.Value) error {

	if ptrval.IsNil() {
		return errors.New(fmt.Sprintf("Cannot build flags from nil pointer for prefix '%s'", prefix))
	}

	return recurseBuildFlags(flags, prefix, ptrval.Elem())

}

func recurseBuildFlags(flags FlagSet, prefix string, elementval reflect.Value) error {

	switch elementval.Kind() {
	case reflect.Map:
		return populateMapFlags(flags, prefix, elementval)

	case reflect.Struct:
		return populateStructFlags(flags, prefix, elementval)

	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		return recursePtrFlags(flags, prefix, elementval)

	default:
		return errors.New(fmt.Sprintf("Cannot build flags from type %v for prefix '%s'", elementval.Type(), prefix))
	}

	return nil

}

// Into populates the given flag set with hierarchical fields from the
// given object.
func Into(flags FlagSet, configuration interface{}) error {

	if configuration == nil {
		return errors.New("Cannot build flags from nil")
	}

	return recurseBuildFlags(flags, "", reflect.ValueOf(configuration))

}

// From populates the top-level default flags with hierarchical fields
// from the given object.  It simply calls Into() with configuration
// on a facade of the top-level flag package functions.
func From(configuration interface{}) error {
	return Into(toplevelflags, configuration)
}
