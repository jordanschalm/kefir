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
		t.Error("Failed to populate struct with mirror source")
	}
}
