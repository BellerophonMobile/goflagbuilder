package goflagbuilder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type flagvalue struct {
	name  string
	value reflect.Value
}

func (x *flagvalue) String() string {
	if x.value.Kind() == reflect.Slice {
		var b strings.Builder
		b.WriteString("[")

		for i := 0; i < x.value.Len(); i++ {
			if i != 0 {
				b.WriteString(", ")
			}
			b.WriteString(x.value.Index(i).String())
		}

		b.WriteString("]")
		return b.String()
	} else {
		return x.value.String()
	}
}

func (x *flagvalue) Set(str string) error {
	if !x.value.IsValid() {
		return fmt.Errorf("Flag variable value of type %s for %s is invalid", x.value.Type().Name(), x.name)
	}

	if !x.value.CanSet() {
		return fmt.Errorf("Flag variable value of type %s for %s cannot be set", x.value.Type(), x.name)
	}

	kind := x.value.Kind()
	if kind == reflect.Slice {
		kind = x.value.Type().Elem().Kind()
	}

	var argValue reflect.Value
	var err error

	switch kind {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(str)
		argValue = reflect.ValueOf(b)

	case reflect.Float64:
		var f float64
		f, err = strconv.ParseFloat(str, 64)
		argValue = reflect.ValueOf(f)

	case reflect.Int64:
		var i int64
		i, err = strconv.ParseInt(str, 0, 64)
		argValue = reflect.ValueOf(i)

	case reflect.Int:
		var i int64
		i, err = strconv.ParseInt(str, 0, 0)
		argValue = reflect.ValueOf(int(i))

	case reflect.String:
		argValue = reflect.ValueOf(str)

	case reflect.Uint64:
		var i uint64
		i, err = strconv.ParseUint(str, 0, 64)
		argValue = reflect.ValueOf(i)

	case reflect.Uint:
		var i uint64
		i, err = strconv.ParseUint(str, 0, 0)
		argValue = reflect.ValueOf(uint(i))

	default:
		return fmt.Errorf("Unsupported flag field variable type %v kind %v for prefix %s", x.value.Type(), x.value.Kind(), x.name)
	}

	if err != nil {
		return nil
	}

	if x.value.Kind() == reflect.Slice {
		argValue = reflect.Append(x.value, argValue)
	}

	x.value.Set(argValue)

	return nil
}

func (x *flagvalue) IsBoolFlag() bool {
	return x.value.Kind() == reflect.Bool
}
