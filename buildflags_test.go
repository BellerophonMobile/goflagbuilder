package goflagbuilder

import (
	"testing"
)

type mystruct struct {
	FieldA string
	FieldB int
}

type myotherstruct struct {
	Grid     uint64
	Fraction float64
}

type mystruct2 struct {
	Name  string
	Index int

	Location myotherstruct
}

type mystruct3 struct {
	Name  string
	Index int

	Location *myotherstruct
}

func Test_From_Invalid(t *testing.T) {
	
	var test *testflags

	test = newtest(t, "Nil", nil)
	test.error("Cannot build flags from nil")
	test.run()

	test = newtest(t, "String", "Banana")
	test.error("Cannot build flags from type string for prefix ''")
	test.run()

	test = newtest(t, "String (not strict)", "Banana")
	test.execute(false)

	test = newtest(t, "Int", 7)
	test.error("Cannot build flags from type int for prefix ''")
	test.run()

	test = newtest(t, "Float", 7.0)
	test.error("Cannot build flags from type float64 for prefix ''")
	test.run()

	test = newtest(t, "Struct", mystruct{"Banana", 7})
	test.error("Value of type string at FieldA cannot be set")
	test.run()

	test = newtest(t, "Struct (not strict)", mystruct{"Banana", 7})
	test.execute(false)

	test = newtest(t, "Map to Struct",
		map[string]interface{}{"MyStruct": mystruct{}})
	test.error("Value of type string at MyStruct.FieldA cannot be set")
	test.run()

	test = newtest(t, "Struct with Nested Nil", &mystruct3{})
	test.variable("Name")
	test.variable("Index")
	test.error("Cannot build flags from nil pointer for prefix 'Location'")
	test.run()

}

func Test_From_Map(t *testing.T) {

	var test *testflags

	test = newtest(t, "Empty map", make(map[string]int))
	test.run()

	test = newtest(t, "Map to Int", map[string]int{"Banana": 7})
	test.variable("Banana")
	test.run()

	test = newtest(t, "Map to Struct Ptr",
		map[string]interface{}{"MyStruct": &mystruct{}})
	test.variable("MyStruct.FieldA")
	test.variable("MyStruct.FieldB")
	test.run()

}

func Test_From_Struct(t *testing.T) {

	var test *testflags

	test = newtest(t, "Struct Ptr", &mystruct{})
	test.variable("FieldA")
	test.variable("FieldB")
	test.run()

	test = newtest(t, "Nested Struct", &mystruct2{})
	test.variable("Name")
	test.variable("Index")
	test.variable("Location.Grid")
	test.variable("Location.Fraction")
	test.run()

	test = newtest(t, "Nested Struct Ptr", &mystruct3{Location: &myotherstruct{}})
	test.variable("Name")
	test.variable("Index")
	test.variable("Location.Grid")
	test.variable("Location.Fraction")
	test.run()

}
