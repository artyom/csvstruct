// Package csvstruct allows decoding of string slice obtained from a
// csv.Reader.Read call into a struct type.
//
// It supports decoding to string, integer, and float struct fields.
package csvstruct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// NewDecoder takes CSV header and dst which must be a struct type or a pointer
// to a struct type with struct "csv" tags mapped to header names, and returns
// a Decoder function for this type and field ordering. It does not modify dst.
//
// Only exported fields are processed.
func NewDecoder(header []string, dst interface{}) (Decoder, error) {
	st := reflect.ValueOf(dst)
	// if st.Kind() != reflect.Ptr {
	// 	return nil, errors.New("dst must be a pointer to a struct type")
	// }
	st = reflect.Indirect(st)
	if !st.IsValid() || st.Type().Kind() != reflect.Struct {
		return nil, errors.New("dst must be either struct type or a pointer to a struct type")
	}
	var setters []setter
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get("csv")
		if tag == "" {
			continue
		}
		val := st.Field(i)
		// if !val.CanAddr() {
		// 	return nil, fmt.Errorf("field with csv tag %q has unsupported type", tag)
		// }
		csvIdx := indexOf(header, tag)
		if csvIdx == -1 {
			continue
		}
		var fn func(reflect.Value, string) error // setter.fn
		switch val.Interface().(type) {
		case int:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseInt(s, 0, 0)
				if err != nil {
					return err
				}
				field.SetInt(x)
				return nil
			}
		case int8:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseInt(s, 0, 8)
				if err != nil {
					return err
				}
				field.SetInt(x)
				return nil
			}
		case int16:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseInt(s, 0, 16)
				if err != nil {
					return err
				}
				field.SetInt(x)
				return nil
			}
		case int32:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseInt(s, 0, 32)
				if err != nil {
					return err
				}
				field.SetInt(x)
				return nil
			}
		case int64:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseInt(s, 0, 64)
				if err != nil {
					return err
				}
				field.SetInt(x)
				return nil
			}
		case uint:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseUint(s, 0, 0)
				if err != nil {
					return err
				}
				field.SetUint(x)
				return nil
			}
		case uint8:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseUint(s, 0, 8)
				if err != nil {
					return err
				}
				field.SetUint(x)
				return nil
			}
		case uint16:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseUint(s, 0, 16)
				if err != nil {
					return err
				}
				field.SetUint(x)
				return nil
			}
		case uint32:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseUint(s, 0, 32)
				if err != nil {
					return err
				}
				field.SetUint(x)
				return nil
			}
		case uint64:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseUint(s, 0, 64)
				if err != nil {
					return err
				}
				field.SetUint(x)
				return nil
			}
		case float32, float64:
			fn = func(field reflect.Value, s string) error {
				x, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}
				field.SetFloat(x)
				return nil
			}
		case bool:
			fn = func(field reflect.Value, s string) error {
				val, err := strconv.ParseBool(s)
				if err != nil {
					return err
				}
				field.SetBool(val)
				return nil
			}
		case string:
			fn = func(field reflect.Value, s string) error {
				field.SetString(s)
				return nil
			}
		default:
			return nil, fmt.Errorf("field with csv tag %q has unsupported type", tag)
		}
		setters = append(setters, setter{
			csvIdx:   csvIdx,
			fieldIdx: i,
			fn:       fn,
		})
	}
	if len(setters) == 0 {
		return nil, errors.New("no matches found between header and csv-tagged struct fields")
	}
	return func(row []string, dst interface{}) error {
		st := reflect.Indirect(reflect.ValueOf(dst))
		for _, s := range setters {
			if err := s.fn(st.Field(s.fieldIdx), row[s.csvIdx]); err != nil {
				return err
			}
		}
		return nil
	}, nil
}

type setter struct {
	csvIdx   int
	fieldIdx int
	fn       func(field reflect.Value, s string) error
}

// Decoder is a function that decodes CSV row to dst, which must be a pointer
// to a struct type. Decoder must be called on the same type it was created
// from by the NewDecoder call.
type Decoder func(row []string, dst interface{}) error

func indexOf(s []string, x string) int {
	for i := range s {
		if x == s[i] {
			return i
		}
	}
	return -1
}