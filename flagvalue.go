package goflagbuilder

import (
	"fmt"
	"reflect"
	"strconv"
)

type flagvalue struct {
	name  string
	value reflect.Value
}

func (x *flagvalue) String() string {
	return x.value.String()
}

func (x *flagvalue) Set(str string) error {

	if !x.value.IsValid() {
		return fmt.Errorf("Flag variable value of type %s for %s is invalid", x.value.Type().Name(), x.name)
	}

	if !x.value.CanSet() {
		return fmt.Errorf("Flag variable value of type %s for %s cannot be set", x.value.Type().Name(), x.name)
	}

	switch x.value.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		x.value.SetBool(b)

	case reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		x.value.SetFloat(f)

	case reflect.Int64:
		i, err := strconv.ParseInt(str, 0, 64)
		if err != nil {
			return err
		}
		x.value.SetInt(i)

	case reflect.Int:
		i, err := strconv.ParseInt(str, 0, 0)
		if err != nil {
			return err
		}
		x.value.SetInt(i)

	case reflect.String:
		x.value.SetString(str)

	case reflect.Uint64:
		i, err := strconv.ParseUint(str, 0, 64)
		if err != nil {
			return err
		}
		x.value.SetUint(i)

	case reflect.Uint:
		i, err := strconv.ParseUint(str, 0, 0)
		if err != nil {
			return err
		}
		x.value.SetUint(i)

	default:
		return fmt.Errorf("Unsupported flag field variable type %v kind %v for prefix %s", x.value.Type(), x.value.Kind(), x.name)
	}

	return nil
}

func (x *flagvalue) IsBool() bool {
	return x.value.Kind() == reflect.Bool
}
