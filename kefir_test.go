package kefir_test

import (
	"testing"

	"github.com/jordanschalm/kefir"
)

var setter = ConfigSetter{}

func init() {
	kefir.SetSource(&setter)
}

type ConfigSetter map[string]string

func (c *ConfigSetter) Get(key string) string {
	return (*c)[key]
}

func TestNonPointer(t *testing.T) {
	c := struct{}{}
	err := kefir.Populate(c)
	if err == nil {
		t.Error("No error after populating struct literal")
	}
}

func TestNonStruct(t *testing.T) {
	c := 1
	err := kefir.Populate(&c)
	if err == nil {
		t.Error("No error after populating pointer to int")
	}
}

func TestSimpleStruct(t *testing.T) {
	c := struct{ A string }{}
	setter["A"] = "great"

	err := kefir.Populate(&c)
	if err != nil {
		t.Error("Unexpected error: " + err.Error())
	}

	if c.A != "great" {
		t.Error("Failed to populate struct with source")
	}
}

func TestPrivateField(t *testing.T) {
	c := struct {
		A string
		a string
	}{}
	setter["A"] = "Great"
	setter["a"] = "great"

	err := kefir.Populate(&c)
	if err != nil {
		t.Error("Unexpected error: " + err.Error())
	}

	if c.A != "Great" {
		t.Error("Failed to populate struct with source")
	}

	if c.a != "" {
		t.Error("Private struct field set")
	}
}

func TestUppercaser(t *testing.T) {
	formatter := &kefir.Uppercaser{Prefix: "PREFIX"}
	kefir.SetFormatter(formatter)
	c := struct{ Asdf string }{}
	setter["PREFIX_ASDF"] = "great"

	kefir.Populate(&c)
	if c.Asdf != "great" {
		t.Error("Failed to set field from source")
	}
}

type CustomFormatter int

func (f *CustomFormatter) Format(field string) string {
	return "BEFORE" + field + "AFTER"
}

func TestCustomFormatter(t *testing.T) {
	formatter := new(CustomFormatter)
	kefir.SetFormatter(formatter)
	c := struct{ A string }{}
	setter[formatter.Format("A")] = "great"

	kefir.Populate(&c)
	if c.A != "great" {
		t.Error("Failed to populate struct with source")
	}
}

func TestIntField(t *testing.T) {}

func TestFloatField(t *testing.T) {}

func TestBoolField(t *testing.T) {}
