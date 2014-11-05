package goflagbuilder

import (
	"fmt"
	"reflect"
	"strconv"
)

type mapvalue struct {
	name       string
	mapval     reflect.Value
	keyval     reflect.Value
	elementval reflect.Value
}

func (x *mapvalue) String() string {
	return x.elementval.String()
}

func (x *mapvalue) Set(str string) error {

	switch x.elementval.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(b))

	case reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(f))

	case reflect.Int64:
		i, err := strconv.ParseInt(str, 0, 64)
		if err != nil {
			return err
		}
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(i))

	case reflect.Int:
		i, err := strconv.ParseInt(str, 0, 0)
		if err != nil {
			return err
		}
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(i))

	case reflect.String:
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(str))

	case reflect.Uint64:
		i, err := strconv.ParseUint(str, 0, 64)
		if err != nil {
			return err
		}
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(i))

	case reflect.Uint:
		i, err := strconv.ParseUint(str, 0, 0)
		if err != nil {
			return err
		}
		x.mapval.SetMapIndex(x.keyval, reflect.ValueOf(i))

	default:
		return fmt.Errorf("Unsupported type %v kind %v for prefix %s", x.elementval.Type(), x.elementval.Kind(), x.name)
	}

	return nil
}

func (x *mapvalue) IsBoolFlag() bool {
	return x.elementval.Kind() == reflect.Bool
}
