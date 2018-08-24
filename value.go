package goflagbuilder

import (
	"flag"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type flagKind interface {
	Set(v reflect.Value, s string) error
	Get(v reflect.Value) interface{}
	String(v reflect.Value) string
}

type value struct {
	value  reflect.Value
	kind   flagKind
	isBool bool
}

func (v value) Set(s string) error {
	return v.kind.Set(v.value, s)
}

func (v value) Get() interface{} { return v.kind.Get(v.value) }

func (v value) String() string { return v.kind.String(v.value) }

func (v value) IsBoolFlag() bool { return v.isBool }

func find(k reflect.Kind) flagKind {
	switch k {
	case reflect.Bool:
		return boolKind{}

	case reflect.Int:
		return intKind{}

	case reflect.Int64:
		return int64Kind{}

	case reflect.Uint:
		return uintKind{}

	case reflect.Uint64:
		return uint64Kind{}

	case reflect.Float64:
		return float64Kind{}

	case reflect.String:
		return stringKind{}
	}

	return nil
}

func isGetter(r reflect.Value) bool {
	return r.Type().Implements(reflect.TypeOf((*flag.Getter)(nil)).Elem())
}

func findKind(r reflect.Value) flagKind {
	kind := find(r.Kind())
	if kind != nil {
		return kind
	}

	if r.Kind() == reflect.Slice {
		kind := sliceKind{itemKind: find(r.Type().Elem().Kind())}
		if kind.itemKind != nil {
			return kind
		}
	}

	log.Println("Type:", r.Type())
	log.Println("Getter:", reflect.TypeOf((*flag.Getter)(nil)).Elem())
	log.Println("Implements:", reflect.PtrTo(r.Type()).Implements(reflect.TypeOf((*flag.Getter)(nil)).Elem()))

	if isGetter(r) || (r.CanAddr() && isGetter(r.Addr())) {
		return getterKind{}
	}

	return nil
}

// bool Kind

type boolKind struct{}

func (boolKind) Set(r reflect.Value, s string) error {
	ps, err := strconv.ParseBool(s)
	r.SetBool(ps)
	return err
}

func (boolKind) Get(r reflect.Value) interface{} { return r.Bool() }

func (boolKind) String(r reflect.Value) string {
	return strconv.FormatBool(r.Bool())
}

// int Kind

type intKind struct{}

func (intKind) Set(r reflect.Value, s string) error {
	ps, err := strconv.ParseInt(s, 0, strconv.IntSize)
	r.SetInt(ps)
	return err
}

func (intKind) Get(r reflect.Value) interface{} { return int(r.Int()) }

func (intKind) String(r reflect.Value) string {
	return strconv.FormatInt(r.Int(), 10)
}

// int64 Kind

type int64Kind struct{}

func (int64Kind) Set(r reflect.Value, s string) error {
	ps, err := strconv.ParseInt(s, 0, 64)
	r.SetInt(ps)
	return err
}

func (int64Kind) Get(r reflect.Value) interface{} { return r.Int() }

func (int64Kind) String(r reflect.Value) string {
	return strconv.FormatInt(r.Int(), 10)
}

// uint Kind

type uintKind struct{}

func (uintKind) Set(r reflect.Value, s string) error {
	ps, err := strconv.ParseUint(s, 0, strconv.IntSize)
	r.SetUint(ps)
	return err
}

func (uintKind) Get(r reflect.Value) interface{} { return uint(r.Uint()) }

func (uintKind) String(r reflect.Value) string {
	return strconv.FormatUint(r.Uint(), 10)
}

// uint64 Kind

type uint64Kind struct{}

func (uint64Kind) Set(r reflect.Value, s string) error {
	ps, err := strconv.ParseUint(s, 0, 64)
	r.SetUint(ps)
	return err
}

func (uint64Kind) Get(r reflect.Value) interface{} { return r.Uint() }

func (uint64Kind) String(r reflect.Value) string {
	return strconv.FormatUint(r.Uint(), 10)
}

// float64 Kind

type float64Kind struct{}

func (float64Kind) Set(r reflect.Value, s string) error {
	ps, err := strconv.ParseFloat(s, 64)
	r.SetFloat(ps)
	return err
}

func (float64Kind) Get(r reflect.Value) interface{} { return r.Float() }

func (float64Kind) String(r reflect.Value) string {
	return strconv.FormatFloat(r.Float(), 'g', -1, 64)
}

// string Kind

type stringKind struct{}

func (stringKind) Set(r reflect.Value, s string) error {
	r.SetString(s)
	return nil
}

func (stringKind) Get(r reflect.Value) interface{} { return r.String() }

func (stringKind) String(r reflect.Value) string { return r.String() }

// flag.Getter Value

type getterKind struct{}

func getter(r reflect.Value) flag.Getter {
	if isGetter(r) {
		return r.Interface().(flag.Getter)
	}

	if r.CanAddr() && isGetter(r.Addr()) {
		return r.Addr().Interface().(flag.Getter)
	}

	panic("not a getter")
}

func (getterKind) Set(r reflect.Value, s string) error {
	fv := getter(r)
	return fv.Set(s)
}

func (getterKind) Get(r reflect.Value) interface{} {
	fv := getter(r)
	return fv.Get()
}

func (getterKind) String(r reflect.Value) string {
	fv := getter(r)
	return fv.String()
}

// slice Kind

type sliceKind struct {
	itemKind flagKind
}

func (v sliceKind) Set(r reflect.Value, s string) error {
	itemValue := reflect.New(r.Type().Elem())
	if err := v.itemKind.Set(itemValue.Elem(), s); err != nil {
		return err
	}

	r.Set(reflect.Append(r, itemValue.Elem()))
	return nil
}

func (sliceKind) Get(r reflect.Value) interface{} { return r.Interface() }

func (sliceKind) String(r reflect.Value) string {
	var b strings.Builder
	b.WriteString("[")

	for i := 0; i < r.Len(); i++ {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(r.Index(i).String())
	}

	b.WriteString("]")
	return b.String()
}

// map Kind

type mapKind struct {
	keyval    reflect.Value
	valueKind flagKind
}

func (v mapKind) Set(r reflect.Value, s string) error {
	itemValue := reflect.New(r.Type().Elem())
	if err := v.valueKind.Set(itemValue.Elem(), s); err != nil {
		return err
	}

	r.SetMapIndex(v.keyval, itemValue.Elem())
	return nil
}

func (v mapKind) Get(r reflect.Value) interface{} {
	return v.valueKind.Get(r.MapIndex(v.keyval))
}

func (v mapKind) String(r reflect.Value) string {
	return v.valueKind.String(r.MapIndex(v.keyval))
}
