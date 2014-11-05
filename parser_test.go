package goflagbuilder

import (
	"testing"
)

func Test_ParserValid(t *testing.T) {

	var test *testflags

	test = newtest(t, "Parse Empty Doc", &mystruct{"Banana", 7})
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse("")

	test = newtest(t, "Parse Comment", &mystruct{"Banana", 7})
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse("# Mushi sushi")

	a := &mystruct{"Banana", 7}
	test = newtest(t, "Parse First Line", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse("FieldA=sushi")
	if a.FieldA != "sushi" {
		t.Error("Did not set field")
	}
	if a.FieldB != 7 {
		t.Error("Set incorrect field")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Two Lines", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse(`FieldA=sushi
              FieldB=9`)
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 9 {
		t.Error("Did not set field FieldB")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Two Lines Preceded by Comment", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse(`# This is a comment
              FieldA=sushi
              FieldB=9`)
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 9 {
		t.Error("Did not set field FieldB")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Two Lines Trailed by Comment", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse(`FieldA=sushi
              FieldB=9
              # This is a comment`)
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 9 {
		t.Error("Did not set field FieldB")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Two Lines Split by Comment", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse(`FieldA=sushi
              # This is a comment
              FieldB=9`)
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 9 {
		t.Error("Did not set field FieldB")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Two Lines Split by Comment and Blank Line", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse(`FieldA=sushi

              # This is a comment
              FieldB=9`)
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 9 {
		t.Error("Did not set field FieldB")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Two Lines With Comments and Blank Lines", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse(`
              FieldA=sushi

              # This is a comment
              # Another comment
              FieldB=9
              `)
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 9 {
		t.Error("Did not set field FieldB")
	}

	a = &mystruct{"Banana", 7}
	test = newtest(t, "Parse Comment on Line", a)
	test.variable("FieldA")
	test.variable("FieldB")
	test.parse("FieldA=sushi # Bananana")
	if a.FieldA != "sushi" {
		t.Error("Did not set field FieldA")
	}
	if a.FieldB != 7 {
		t.Error("Set incorrect field FieldB")
	}

}

func Test_ParserInValid(t *testing.T) {

	var test *testflags

	test = newtest(t, "Parse No Key", &mystruct{"Banana", 7})
	test.variable("FieldA")
	test.variable("FieldB")
	test.parseerror("Line 1 has no key")
	test.parse("Sushi")

	test = newtest(t, "Parse No Key Line 2", &mystruct{"Banana", 7})
	test.variable("FieldA")
	test.variable("FieldB")
	test.parseerror("Line 2 has no key")
	test.parse(`# Line 1
             sushi`)

}
