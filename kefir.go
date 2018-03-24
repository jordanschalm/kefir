package kefir

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

var (
	// preface is the string prepended to configuration keys
	preface string
	// source is the source to be used to populate configs
	source Source = new(OS)
)

// SetPreface sets the preface used when retrieving configuration values
func SetPreface(p string) {
	preface = p
}

// SetSource sets the source used when retrieving configuation values
func SetSource(s Source) {
	source = s
}

// Source defines a source for retrieving config values.
type Source interface {
	Get(string) string
}

// OS is a source that retrieves config from environment variables
type OS int

// Get simply wraps os.Getenv
func (o *OS) Get(key string) string {
	return os.Getenv(key)
}

// Populate populates
func Populate(conf interface{}) error {
	t := reflect.TypeOf(conf)
	if t.Kind() != reflect.Ptr {
		return newError("Populate must only be called with a pointer")
	}
	// Dereference the pointer and ensure it points to a struct
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return newError("Populate must only be called with a struct pointer")
	}

	v := reflect.ValueOf(conf)

	for i := 0; i < t.NumField(); i++ {
		fieldV := v.Elem().Field(i)
		fieldT := t.Field(i)

		key := strings.ToUpper(fieldT.Name)
		if len(preface) > 0 {
			key = preface + "_" + key
		}
		value := source.Get(key)

		if fieldV.CanSet() {
			fieldV.SetString(value)
		} else {
			panic("Failed to set config field: " + fieldT.Name)
		}
	}

	return nil
}

func newError(msg string) error {
	return errors.New("kefir: " + msg)
}
