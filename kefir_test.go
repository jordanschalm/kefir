package kefir_test

import (
	"strconv"
	"testing"

	"github.com/jordanschalm/kefir"
)

var setter = ConfigSetter{}

func init() {
	kefir.SetSource(&setter)
}

type ConfigSetter map[string]string

func (c *ConfigSetter) Get(key string) (string, bool) {
	v, ok := (*c)[key]
	return v, ok
}

func TestSetNilSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Setting nil source didn't panic")
		}
	}()
	kefir.SetSource(nil)
}

func TestSetNilFormatter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Setting nil formatter didn't panic")
		}
	}()
	kefir.SetFormatter(nil)
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
		t.Error("No error after populating non-struct pointer")
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
		t.Errorf("Failed to populate struct with source. Got: %s\tExpected: %s", c.A, "great")
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
		t.Errorf("Failed to populate struct with source. Got %s\tExpected: %s", c.A, "Great")
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
		t.Errorf("Failed to set field from source. Got: %s\tExpected: %s", c.Asdf, "great")
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
		t.Errorf("Failed to populate struct with source. Got: %s\tExpected: %s", c.A, "great")
	}

	// Reset formatter to prevent interactions with other tests
	// TODO refactor to allow creating multiple Kefir instances with their own state
	kefir.SetFormatter(new(kefir.Uppercaser))
}

func TestBoolField(t *testing.T) {
	c := struct{ A, B, C, D bool }{}

	setter["A"] = "true"
	setter["B"] = "false"
	setter["C"] = "T"
	setter["D"] = "F"

	kefir.Populate(&c)

	if !c.A || c.B || !c.C || c.D {
		t.Error("Failed to populate config with source")
	}
}

func TestIntField(t *testing.T) {
	c := struct {
		A int8
		B int16
		C int32
		D int64
		E int
	}{}
	setter["A"] = strconv.FormatInt(1>>7, 10)
	setter["B"] = strconv.FormatInt(1>>15, 10)
	setter["C"] = strconv.FormatInt(1>>31, 10)
	setter["D"] = strconv.FormatInt(1>>63, 10)
	setter["E"] = "42"

	kefir.Populate(&c)

	if c.A != 1>>7 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.A, 1>>7)
	}
	if c.B != 1>>15 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.B, 1>>15)
	}
	if c.C != 1>>31 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.C, 1>>31)
	}
	if c.D != 1>>63 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.D, 1>>63)
	}
	if c.E != 42 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.E, 42)
	}
}

func TestUIntField(t *testing.T) {
	c := struct {
		A uint8
		B uint16
		C uint32
		D uint64
		E uint
	}{}
	setter["A"] = strconv.FormatInt(1>>7, 10)
	setter["B"] = strconv.FormatInt(1>>15, 10)
	setter["C"] = strconv.FormatInt(1>>31, 10)
	setter["D"] = strconv.FormatInt(1>>63, 10)
	setter["E"] = "42"

	kefir.Populate(&c)

	if c.A != 1>>7 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.A, 1>>7)
	}
	if c.B != 1>>15 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.B, 1>>15)
	}
	if c.C != 1>>31 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.C, 1>>31)
	}
	if c.D != 1>>63 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.D, 1>>63)
	}
	if c.E != 42 {
		t.Errorf("Failed to populate field. Got: %d\tExpected: %d", c.E, 42)
	}
}

func TestFloatField(t *testing.T) {
	c := struct {
		A float32
		B float64
	}{}

	setter["A"] = "3.14159"
	setter["B"] = "3.141592653"

	kefir.Populate(&c)

	if c.A != 3.14159 {
		t.Errorf("Failed to populate field. Got: %f\tExpected: %f", c.A, 3.14159)
	}
	if c.B != 3.141592653 {
		t.Errorf("Failed to populate field. Got: %f\tExpected: %f", c.B, 3.141592653)
	}
}

func TestDefaultField(t *testing.T) {
	c := struct {
		A string `default:"defaultA"`
		B string `default:"defaultB"`
		C string
	}{}

	setter["A"] = "setA"
	delete(setter, "B")
	delete(setter, "C")

	kefir.Populate(&c)

	if c.A != "setA" {
		t.Errorf("Failed to populate field. Got: %s\tExpected: %s", c.A, "setA")
	}
	if c.B != "defaultB" {
		t.Errorf("Failed to populate field. Got: %s\tExpected: %s", c.B, "defaultB")
	}
	if c.C != "" {
		t.Errorf("Empty field populated. Got: %s\t", c.C)
	}
}
