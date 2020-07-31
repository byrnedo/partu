package partu

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	errTagMissing        = errors.New("struct field must have `sql` tag if included in Columns() output")
	errFieldNotSupported = errors.New("at least one of the field types in Columns() is not supported")
)


func (p Builder) ColName(table Table, field interface{}) string {

	v := reflect.ValueOf(table)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	ft, err := p.findFieldTag(v, reflect.ValueOf(field))
	if err != nil {
		panic(err)
	}
	return ft
}

func (p Builder) NamedFields(table Table) (ret namedFields) {
	t := reflect.TypeOf(table)

	v := reflect.ValueOf(table)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i, col := range table.Columns() {
		ft, err := p.findFieldTag(v, reflect.ValueOf(col))
		if err != nil {
			panic(fmt.Sprintf("Columns index %d: %s", i, err.Error()))
		}

		ret = append(ret, namedField{Name: ft, Field: col})
	}
	return
}

func (b Builder) findFieldTag(structValue reflect.Value, fieldValue reflect.Value) (string, error) {
	sf := b.findStructField(structValue, fieldValue)
	if sf == nil {
		return "", errFieldNotSupported
	}

	colName := sf.Tag.Get(b.Tag())
	if colName == "" {
		return "", errTagMissing
	}
	colName = strings.Split(colName, ",")[0]
	return colName, nil
}

// Author: Copied from https://github.com/go-ozzo/ozzo-validation/
//
// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the field info will be returned. Otherwise, nil will be returned.
func (b Builder) findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	t := fieldValue.Elem().Type()

	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		sfV := structValue.Field(i)
		if ptr == sfV.UnsafeAddr() {
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if sf.Type.Kind() == reflect.Struct {
				if sf.Type == t {
					return &sf
				}
			} else {
				return &sf
			}
		}
		if sf.Anonymous {
			// delve into anonymous struct to look for the field
			fi := structValue.Field(i)
			if sf.Type.Kind() == reflect.Ptr {
				fi = fi.Elem()
			}
			if fi.Kind() == reflect.Struct {
				if f := b.findStructField(fi, fieldValue); f != nil {
					return f
				}
			}
		}
	}
	return nil
}
