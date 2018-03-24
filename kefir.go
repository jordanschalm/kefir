package kefir

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

var (
	// The transformer used when populating configs. Uses unprefixed uppercaser
	// by default.
	formatter Formatter = new(Uppercaser)
	// The source to be used to populate configs. Uses OS environment variables
	// by default.
	source Source = new(OS)
)

// SetFormatter sets the formatter used when populating configs.
func SetFormatter(f Formatter) {
	// Guarantee that calls on the formatter will never panic
	if f == nil {
		return
	}
	formatter = f
}

// SetSource sets the source used when retrieving configuation values.
func SetSource(s Source) {
	// Guarantee that calls on the source will never panic
	if s == nil {
		return
	}
	source = s
}

// Source defines a source for retrieving config values.
type Source interface {
	Get(string) string
}

// Formatter defines how config keys should be formatted. Format defines
// a function that converts field names in the configuation struct to
// corresponding key names in the source.
type Formatter interface {
	Format(string) string
}

// OS is a source that retrieves config from environment variables.
type OS struct{}

// Get simply wraps os.Getenv.
func (o *OS) Get(key string) string {
	return os.Getenv(key)
}

// Uppercaser formats config keys by optionally prepending a prefix and
// uppercasing the key.
type Uppercaser struct {
	Prefix string
}

// Format optionally prepends a prefix and uppercases the key.
func (u *Uppercaser) Format(field string) string {
	if len(u.Prefix) > 0 {
		return strings.ToUpper(u.Prefix + "_" + field)
	}
	return strings.ToUpper(field)
}

// Populate populates a configuration struct by retrieving values from the set
// source that correspond to fields in the input struct.
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

		if !fieldV.CanSet() {
			continue
		}

		key := formatter.Format(fieldT.Name)
		value := source.Get(key)

		fieldV.SetString(value)
	}

	return nil
}

func newError(msg string) error {
	return errors.New("kefir: " + msg)
}
