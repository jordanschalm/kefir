package kefir

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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
		panic("Attempted to set formatter to nil")
	}
	formatter = f
}

// SetSource sets the source used when retrieving configuation values.
func SetSource(s Source) {
	// Guarantee that calls on the source will never panic
	if s == nil {
		panic("Attempted to set source to nil")
	}
	source = s
}

// Source defines a source for retrieving config values.
type Source interface {
	Get(string) (string, bool)
}

// Formatter defines how config keys should be formatted. Format defines
// a function that converts field names in the configuation struct to
// corresponding key names in the source.
type Formatter interface {
	Format(string) string
}

// OS is a source that retrieves config from environment variables.
type OS struct{}

// Get simply wraps os.LookupEnv.
func (o *OS) Get(key string) (string, bool) {
	return os.LookupEnv(key)
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
	v = v.Elem()

	for i := 0; i < t.NumField(); i++ {
		fieldV := v.Field(i)
		fieldT := t.Field(i)

		// Ignore the field if we can't set it
		if !fieldV.CanSet() {
			continue
		}

		key := formatter.Format(fieldT.Name)
		value, ok := source.Get(key)
		if !ok {
			// Couldn't find value. Use default if one is set.
			value, _ = fieldT.Tag.Lookup("default")
		}

		var err error

		switch fieldV.Kind() {
		case reflect.String:
			fieldV.SetString(value)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				break
			}
			fieldV.SetBool(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var v int64
			// Check for case of time.Duration
			tp := fieldV.Type()
			if tp.Name() == "Duration" && tp.PkgPath() == "time" {
				dv, err := time.ParseDuration(value)
				if err != nil {
					break
				}
				v = int64(dv)
			} else {
				v, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					break
				}
			}
			fieldV.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				break
			}
			fieldV.SetUint(v)
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				break
			}
			fieldV.SetFloat(v)
		}
	}

	return nil
}

func newError(msg string) error {
	return errors.New("kefir: " + msg)
}
